package app

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"code-review-agent/backend/internal/auth"
	"code-review-agent/backend/internal/config"
	"code-review-agent/backend/internal/database"
	"code-review-agent/backend/internal/models"
	"code-review-agent/backend/internal/respond"
	reviewservice "code-review-agent/backend/internal/review"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	cfg     config.Config
	db      *pgxpool.Pool
	reviews *reviewservice.Service
}

func New(cfg config.Config, db *pgxpool.Pool) *App {
	if err := database.EnsureSchema(context.Background(), db); err != nil {
		panic(err)
	}
	return &App{cfg: cfg, db: db, reviews: &reviewservice.Service{DB: db, WorkDir: cfg.WorkDir}}
}

func (a *App) Routes() http.Handler {
	router := chi.NewRouter()
	router.Use(middleware.RequestID, middleware.RealIP, middleware.Recoverer)
	router.Use(a.cors)
	router.Get("/health", a.health)
	router.Post("/api/auth/register", a.register)
	router.Post("/api/auth/login", a.login)
	router.Post("/api/webhooks/github", a.githubWebhook)
	router.Group(func(protected chi.Router) {
		protected.Use(auth.Middleware(a.cfg.JWTSecret))
		protected.Get("/api/me", a.me)
		protected.Get("/api/projects", a.listProjects)
		protected.Post("/api/projects", a.createProject)
		protected.Post("/api/projects/generate", a.generateProject)
		protected.Get("/api/projects/generated", a.listGeneratedProjects)
		protected.Get("/api/projects/generated/{generatedProjectID}", a.getGeneratedProject)
		protected.Get("/api/projects/generated/{generatedProjectID}/file", a.getGeneratedProjectFile)
		protected.Put("/api/projects/generated/{generatedProjectID}/file", a.saveGeneratedProjectFile)
		protected.Post("/api/projects/generated/{generatedProjectID}/sync-git", a.syncGeneratedProjectGit)
		protected.Get("/api/projects/{projectID}/reviews", a.listReviews)
		protected.Post("/api/projects/{projectID}/reviews/upload", a.uploadReview)
		protected.Post("/api/projects/{projectID}/reviews/upload-file", a.uploadSingleReview)
		protected.Post("/api/projects/{projectID}/reviews/paste", a.pasteReview)
		protected.Post("/api/projects/{projectID}/reviews/git", a.gitReview)
		protected.Get("/api/reviews/{reviewID}", a.getReview)
		protected.Get("/api/reviews/{reviewID}/findings", a.listFindings)
	})
	return router
}

func (a *App) cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if a.isAllowedOrigin(origin) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
		}
		w.Header().Set("Vary", "Origin")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, Accept, X-Requested-With, X-Hub-Signature-256, X-GitHub-Event, X-GitHub-Delivery")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (a *App) isAllowedOrigin(origin string) bool {
	if origin == "" {
		return false
	}
	for _, allowed := range a.cfg.CORSOrigins {
		if allowed == "*" || allowed == origin {
			return true
		}
	}
	for _, prefix := range []string{"http://localhost:", "http://127.0.0.1:", "http://[::1]:", "https://localhost:", "https://127.0.0.1:", "https://[::1]:"} {
		if strings.HasPrefix(origin, prefix) {
			return true
		}
	}
	return false
}

func (a *App) health(w http.ResponseWriter, r *http.Request) {
	if err := a.db.Ping(r.Context()); err != nil {
		respond.Error(w, 503, "database unavailable")
		return
	}
	respond.JSON(w, 200, map[string]string{"status": "ok"})
}

type credentials struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (a *App) register(w http.ResponseWriter, r *http.Request) {
	var input credentials
	if respond.DecodeJSON(r, &input) != nil || strings.TrimSpace(input.Name) == "" || !strings.Contains(input.Email, "@") || len(input.Password) < 8 {
		respond.Error(w, 400, "name, valid email and password of at least 8 characters are required")
		return
	}
	hash, err := auth.HashPassword(input.Password)
	if err != nil {
		respond.Error(w, 500, "password processing failed")
		return
	}
	var user models.User
	err = a.db.QueryRow(r.Context(), `INSERT INTO users(name,email,password_hash) VALUES($1,LOWER($2),$3) RETURNING id,name,email,role,created_at`, input.Name, input.Email, hash).Scan(&user.ID, &user.Name, &user.Email, &user.Role, &user.CreatedAt)
	if err != nil {
		respond.Error(w, 409, "email already registered")
		return
	}
	token, _ := auth.CreateToken(user.ID, a.cfg.JWTSecret)
	respond.JSON(w, 201, map[string]any{"token": token, "user": user})
}

func (a *App) login(w http.ResponseWriter, r *http.Request) {
	var input credentials
	if respond.DecodeJSON(r, &input) != nil {
		respond.Error(w, 400, "invalid request")
		return
	}
	var user models.User
	var hash string
	err := a.db.QueryRow(r.Context(), `SELECT id,name,email,role,created_at,password_hash FROM users WHERE email=LOWER($1)`, input.Email).Scan(&user.ID, &user.Name, &user.Email, &user.Role, &user.CreatedAt, &hash)
	if err != nil || !auth.CheckPassword(hash, input.Password) {
		respond.Error(w, 401, "invalid email or password")
		return
	}
	token, _ := auth.CreateToken(user.ID, a.cfg.JWTSecret)
	respond.JSON(w, 200, map[string]any{"token": token, "user": user})
}

func (a *App) me(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := a.db.QueryRow(r.Context(), `SELECT id,name,email,role,created_at FROM users WHERE id=$1`, auth.UserID(r.Context())).Scan(&user.ID, &user.Name, &user.Email, &user.Role, &user.CreatedAt)
	if err != nil {
		respond.Error(w, 404, "user not found")
		return
	}
	respond.JSON(w, 200, user)
}

func (a *App) createProject(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	if respond.DecodeJSON(r, &input) != nil || strings.TrimSpace(input.Name) == "" {
		respond.Error(w, 400, "project name is required")
		return
	}
	var project models.Project
	err := a.db.QueryRow(r.Context(), `INSERT INTO projects(owner_id,name,description) VALUES($1,$2,$3) RETURNING id,owner_id,name,description,created_at`, auth.UserID(r.Context()), input.Name, input.Description).Scan(&project.ID, &project.OwnerID, &project.Name, &project.Description, &project.CreatedAt)
	if err != nil {
		respond.Error(w, 500, "could not create project")
		return
	}
	_, _ = a.db.Exec(r.Context(), `INSERT INTO project_members(project_id,user_id,member_role) VALUES($1,$2,'owner') ON CONFLICT DO NOTHING`, project.ID, project.OwnerID)
	respond.JSON(w, 201, project)
}

func (a *App) listProjects(w http.ResponseWriter, r *http.Request) {
	rows, err := a.db.Query(r.Context(), `SELECT p.id,p.owner_id,p.name,p.description,p.created_at FROM projects p JOIN project_members m ON m.project_id=p.id WHERE m.user_id=$1 ORDER BY p.created_at DESC`, auth.UserID(r.Context()))
	if err != nil {
		respond.Error(w, 500, "could not load projects")
		return
	}
	defer rows.Close()
	result := []models.Project{}
	for rows.Next() {
		var item models.Project
		_ = rows.Scan(&item.ID, &item.OwnerID, &item.Name, &item.Description, &item.CreatedAt)
		result = append(result, item)
	}
	respond.JSON(w, 200, result)
}

func (a *App) listReviews(w http.ResponseWriter, r *http.Request) {
	projectID := chi.URLParam(r, "projectID")
	if !a.isMember(r, projectID) {
		respond.Error(w, 403, "forbidden")
		return
	}
	rows, err := a.db.Query(r.Context(), `SELECT id,project_id,requested_by,source_type,source_name,branch_name,commit_sha,status,summary,files_scanned,critical_count,high_count,medium_count,low_count,error_message,created_at,completed_at FROM reviews WHERE project_id=$1 ORDER BY created_at DESC`, projectID)
	if err != nil {
		respond.Error(w, 500, "could not load reviews")
		return
	}
	defer rows.Close()
	result := []models.Review{}
	for rows.Next() {
		var x models.Review
		_ = rows.Scan(&x.ID, &x.ProjectID, &x.RequestedBy, &x.SourceType, &x.SourceName, &x.BranchName, &x.CommitSHA, &x.Status, &x.Summary, &x.FilesScanned, &x.CriticalCount, &x.HighCount, &x.MediumCount, &x.LowCount, &x.ErrorMessage, &x.CreatedAt, &x.CompletedAt)
		x.SourceName = a.decryptStoredValue(x.SourceName)
		if x.BranchName != nil {
			branch := a.decryptStoredValue(*x.BranchName)
			x.BranchName = &branch
		}
		result = append(result, x)
	}
	respond.JSON(w, 200, result)
}

func (a *App) uploadReview(w http.ResponseWriter, r *http.Request) {
	projectID := chi.URLParam(r, "projectID")
	if !a.isMember(r, projectID) {
		respond.Error(w, 403, "forbidden")
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, a.cfg.MaxUploadMB*1024*1024)
	if err := r.ParseMultipartForm(a.cfg.MaxUploadMB * 1024 * 1024); err != nil {
		respond.Error(w, 400, "invalid or oversized upload")
		return
	}
	file, header, err := r.FormFile("file")
	if err != nil {
		respond.Error(w, 400, "ZIP file is required")
		return
	}
	defer file.Close()
	if !strings.HasSuffix(strings.ToLower(header.Filename), ".zip") {
		respond.Error(w, 400, "only ZIP files are supported")
		return
	}
	reviewID := uuid.NewString()
	uploadDir := filepath.Join(a.cfg.WorkDir, "uploads")
	_ = os.MkdirAll(uploadDir, 0700)
	zipPath := filepath.Join(uploadDir, reviewID+".zip")
	output, err := os.OpenFile(zipPath, os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		respond.Error(w, 500, "could not store upload")
		return
	}
	_, err = io.Copy(output, file)
	output.Close()
	if err != nil {
		respond.Error(w, 500, "could not store upload")
		return
	}
	_, err = a.db.Exec(r.Context(), `INSERT INTO reviews(id,project_id,requested_by,source_type,source_name) VALUES($1,$2,$3,'upload',$4)`, reviewID, projectID, auth.UserID(r.Context()), a.encryptStoredValue(header.Filename))
	if err != nil {
		respond.Error(w, 500, "could not create review")
		return
	}
	go func() { defer os.Remove(zipPath); a.reviews.ProcessUpload(reviewID, zipPath) }()
	respond.JSON(w, 202, map[string]string{"reviewId": reviewID, "status": "queued"})
}

func (a *App) uploadSingleReview(w http.ResponseWriter, r *http.Request) {
	projectID := chi.URLParam(r, "projectID")
	if !a.isMember(r, projectID) {
		respond.Error(w, 403, "forbidden")
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, a.cfg.MaxUploadMB*1024*1024)
	if err := r.ParseMultipartForm(a.cfg.MaxUploadMB * 1024 * 1024); err != nil {
		respond.Error(w, 400, "invalid or oversized upload")
		return
	}
	file, header, err := r.FormFile("file")
	if err != nil {
		respond.Error(w, 400, "source file is required")
		return
	}
	defer file.Close()
	if !isReviewableFilename(header.Filename) {
		respond.Error(w, 400, "only source code files are supported (.js, .ts, .tsx, .py, .go, .java, .cs, .sql, .yml, .yaml, .json, .sh, .txt)")
		return
	}
	reviewID := uuid.NewString()
	uploadDir := filepath.Join(a.cfg.WorkDir, "uploads")
	_ = os.MkdirAll(uploadDir, 0700)
	tempPath := filepath.Join(uploadDir, reviewID+"-"+safeReviewFilename(header.Filename))
	output, err := os.OpenFile(tempPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		respond.Error(w, 500, "could not store upload")
		return
	}
	if _, err = io.Copy(output, file); err != nil {
		output.Close()
		respond.Error(w, 500, "could not store upload")
		return
	}
	output.Close()
	_, err = a.db.Exec(r.Context(), `INSERT INTO reviews(id,project_id,requested_by,source_type,source_name) VALUES($1,$2,$3,'upload-file',$4)`, reviewID, projectID, auth.UserID(r.Context()), a.encryptStoredValue(header.Filename))
	if err != nil {
		respond.Error(w, 500, "could not create review")
		return
	}
	go func() { defer os.Remove(tempPath); a.reviews.ProcessSingleFile(reviewID, tempPath, header.Filename) }()
	respond.JSON(w, 202, map[string]string{"reviewId": reviewID, "status": "queued"})
}

func (a *App) pasteReview(w http.ResponseWriter, r *http.Request) {
	projectID := chi.URLParam(r, "projectID")
	if !a.isMember(r, projectID) {
		respond.Error(w, 403, "forbidden")
		return
	}
	var input struct {
		Filename string `json:"filename"`
		Language string `json:"language"`
		Code     string `json:"code"`
	}
	if respond.DecodeJSON(r, &input) != nil || strings.TrimSpace(input.Code) == "" {
		respond.Error(w, 400, "code text is required")
		return
	}
	filename := strings.TrimSpace(input.Filename)
	if filename == "" {
		filename = "pasted-code" + languageToExtension(input.Language)
	}
	if !isReviewableFilename(filename) {
		filename = strings.TrimSuffix(filename, filepath.Ext(filename)) + languageToExtension(input.Language)
	}
	if !isReviewableFilename(filename) {
		filename += ".txt"
	}
	reviewID := uuid.NewString()
	uploadDir := filepath.Join(a.cfg.WorkDir, "uploads")
	_ = os.MkdirAll(uploadDir, 0700)
	tempPath := filepath.Join(uploadDir, reviewID+"-"+safeReviewFilename(filename))
	if err := os.WriteFile(tempPath, []byte(input.Code), 0600); err != nil {
		respond.Error(w, 500, "could not store pasted code")
		return
	}
	_, err := a.db.Exec(r.Context(), `INSERT INTO reviews(id,project_id,requested_by,source_type,source_name) VALUES($1,$2,$3,'paste',$4)`, reviewID, projectID, auth.UserID(r.Context()), a.encryptStoredValue(filename))
	if err != nil {
		respond.Error(w, 500, "could not create review")
		return
	}
	go func() { defer os.Remove(tempPath); a.reviews.ProcessSingleFile(reviewID, tempPath, filename) }()
	respond.JSON(w, 202, map[string]string{"reviewId": reviewID, "status": "queued"})
}

func (a *App) gitReview(w http.ResponseWriter, r *http.Request) {
	projectID := chi.URLParam(r, "projectID")
	if !a.isMember(r, projectID) {
		respond.Error(w, 403, "forbidden")
		return
	}
	var input struct {
		RepositoryURL       string `json:"repositoryUrl"`
		Branch              string `json:"branch"`
		FeatureRequirements string `json:"featureRequirements"`
		CommitMessage       string `json:"commitMessage"`
		AutoPush            bool   `json:"autoPush"`
	}
	if respond.DecodeJSON(r, &input) != nil {
		respond.Error(w, 400, "invalid request")
		return
	}
	if strings.TrimSpace(input.RepositoryURL) == "" {
		respond.Error(w, 400, "repository URL is required")
		return
	}
	reviewID := uuid.NewString()
	_, err := a.db.Exec(r.Context(), `INSERT INTO reviews(id,project_id,requested_by,source_type,source_name,branch_name) VALUES($1,$2,$3,'git',$4,$5)`, reviewID, projectID, auth.UserID(r.Context()), a.encryptStoredValue(input.RepositoryURL), a.encryptStoredValue(input.Branch))
	if err != nil {
		respond.Error(w, 500, "could not create review")
		return
	}
	if strings.TrimSpace(input.FeatureRequirements) == "" {
		go a.reviews.ProcessGit(reviewID, input.RepositoryURL, input.Branch)
	} else {
		go a.reviews.ProcessGitFeature(reviewID, input.RepositoryURL, input.Branch, input.FeatureRequirements, input.CommitMessage, input.AutoPush)
	}
	respond.JSON(w, 202, map[string]string{"reviewId": reviewID, "status": "queued"})
}

func (a *App) getReview(w http.ResponseWriter, r *http.Request) {
	var x models.Review
	err := a.db.QueryRow(r.Context(), `SELECT id,project_id,requested_by,source_type,source_name,branch_name,commit_sha,status,summary,files_scanned,critical_count,high_count,medium_count,low_count,error_message,created_at,completed_at FROM reviews WHERE id=$1`, chi.URLParam(r, "reviewID")).Scan(&x.ID, &x.ProjectID, &x.RequestedBy, &x.SourceType, &x.SourceName, &x.BranchName, &x.CommitSHA, &x.Status, &x.Summary, &x.FilesScanned, &x.CriticalCount, &x.HighCount, &x.MediumCount, &x.LowCount, &x.ErrorMessage, &x.CreatedAt, &x.CompletedAt)
	if err != nil {
		respond.Error(w, 404, "review not found")
		return
	}
	if !a.isMember(r, x.ProjectID) {
		respond.Error(w, 403, "forbidden")
		return
	}
	x.SourceName = a.decryptStoredValue(x.SourceName)
	if x.BranchName != nil {
		branch := a.decryptStoredValue(*x.BranchName)
		x.BranchName = &branch
	}
	respond.JSON(w, 200, x)
}
func (a *App) listFindings(w http.ResponseWriter, r *http.Request) {
	reviewID := chi.URLParam(r, "reviewID")
	var projectID string
	if a.db.QueryRow(r.Context(), `SELECT project_id FROM reviews WHERE id=$1`, reviewID).Scan(&projectID) != nil {
		respond.Error(w, 404, "review not found")
		return
	}
	if !a.isMember(r, projectID) {
		respond.Error(w, 403, "forbidden")
		return
	}
	rows, err := a.db.Query(r.Context(), `SELECT id,review_id,severity,category,rule_id,file_path,line_number,title,description,recommendation,code_excerpt,source,status,created_at FROM findings WHERE review_id=$1 ORDER BY CASE severity WHEN 'CRITICAL' THEN 1 WHEN 'HIGH' THEN 2 WHEN 'MEDIUM' THEN 3 WHEN 'LOW' THEN 4 ELSE 5 END,file_path,line_number`, reviewID)
	if err != nil {
		respond.Error(w, 500, "could not load findings")
		return
	}
	defer rows.Close()
	result := []models.Finding{}
	for rows.Next() {
		var x models.Finding
		_ = rows.Scan(&x.ID, &x.ReviewID, &x.Severity, &x.Category, &x.RuleID, &x.FilePath, &x.LineNumber, &x.Title, &x.Description, &x.Recommendation, &x.CodeExcerpt, &x.Source, &x.Status, &x.CreatedAt)
		result = append(result, x)
	}
	respond.JSON(w, 200, result)
}
func (a *App) isMember(r *http.Request, projectID string) bool {
	var exists bool
	_ = a.db.QueryRow(r.Context(), `SELECT EXISTS(SELECT 1 FROM project_members WHERE project_id=$1 AND user_id=$2)`, projectID, auth.UserID(r.Context())).Scan(&exists)
	return exists
}

func isReviewableFilename(name string) bool {
	ext := strings.ToLower(filepath.Ext(name))
	switch ext {
	case ".go", ".js", ".jsx", ".ts", ".tsx", ".py", ".java", ".cs", ".sql", ".yml", ".yaml", ".json", ".xml", ".tf", ".sh", ".txt":
		return true
	default:
		return false
	}
}

func safeReviewFilename(name string) string {
	base := filepath.Base(name)
	clean := strings.Map(func(r rune) rune {
		switch {
		case r >= 'a' && r <= 'z':
			return r
		case r >= 'A' && r <= 'Z':
			return r
		case r >= '0' && r <= '9':
			return r
		case r == '-', r == '_', r == '.', r == ' ':
			return r
		default:
			return '_'
		}
	}, base)
	if clean == "" || clean == "." || clean == ".." {
		return "snippet.txt"
	}
	return clean
}

func languageToExtension(language string) string {
	switch strings.ToLower(strings.TrimSpace(language)) {
	case "javascript", "js":
		return ".js"
	case "typescript", "ts":
		return ".ts"
	case "jsx":
		return ".jsx"
	case "tsx":
		return ".tsx"
	case "python", "py":
		return ".py"
	case "go", "golang":
		return ".go"
	case "java":
		return ".java"
	case "csharp", "c#", "cs":
		return ".cs"
	case "sql":
		return ".sql"
	case "yaml", "yml":
		return ".yml"
	case "json":
		return ".json"
	case "shell", "bash", "sh":
		return ".sh"
	default:
		return ".txt"
	}
}

func (a *App) githubWebhook(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(http.MaxBytesReader(w, r.Body, 2*1024*1024))
	if err != nil {
		respond.Error(w, 400, "invalid payload")
		return
	}
	if a.cfg.GitHubWebhookSecret == "" || !validSignature(body, r.Header.Get("X-Hub-Signature-256"), a.cfg.GitHubWebhookSecret) {
		respond.Error(w, 401, "invalid webhook signature")
		return
	}
	delivery := r.Header.Get("X-GitHub-Delivery")
	event := r.Header.Get("X-GitHub-Event")
	if delivery == "" {
		respond.Error(w, 400, "missing delivery id")
		return
	}
	tag, err := a.db.Exec(r.Context(), `INSERT INTO webhook_deliveries(delivery_id,event_name) VALUES($1,$2) ON CONFLICT DO NOTHING`, delivery, event)
	if err != nil {
		respond.Error(w, 500, "could not store delivery")
		return
	}
	if tag.RowsAffected() == 0 {
		respond.JSON(w, 200, map[string]string{"status": "duplicate"})
		return
	}
	var payload map[string]any
	_ = json.Unmarshal(body, &payload)
	respond.JSON(w, 202, map[string]string{"status": "accepted", "event": event, "message": "Webhook verified. Map the GitHub repository installation to a project before enabling automatic review creation."})
}
func validSignature(body []byte, signature, secret string) bool {
	if !strings.HasPrefix(signature, "sha256=") {
		return false
	}
	expected := hmac.New(sha256.New, []byte(secret))
	expected.Write(body)
	provided, err := hex.DecodeString(strings.TrimPrefix(signature, "sha256="))
	return err == nil && hmac.Equal(expected.Sum(nil), provided)
}
