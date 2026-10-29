package app

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"code-review-agent/backend/internal/auth"
	"code-review-agent/backend/internal/respond"
	"code-review-agent/backend/internal/security"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type generateProjectRequest struct {
	ProjectName        string `json:"projectName"`
	Requirements       string `json:"requirements"`
	FrontendTechnology string `json:"frontendTechnology"`
	BackendTechnology  string `json:"backendTechnology"`
	DatabaseTechnology string `json:"databaseTechnology"`
	GitURL             string `json:"gitUrl"`
	GitBranch          string `json:"gitBranch"`
}

type generatedProjectResponse struct {
	ID                 string   `json:"id"`
	ProjectName        string   `json:"projectName"`
	Requirements       string   `json:"requirements"`
	FrontendTechnology string   `json:"frontendTechnology"`
	BackendTechnology  string   `json:"backendTechnology"`
	DatabaseTechnology string   `json:"databaseTechnology"`
	GitURL             string   `json:"gitUrl"`
	GitBranch          string   `json:"gitBranch"`
	RootPath           string   `json:"rootPath"`
	Files              []string `json:"files"`
}

type saveGeneratedFileRequest struct {
	Path    string `json:"path"`
	Content string `json:"content"`
}

type syncGeneratedGitRequest struct {
	GitURL    string `json:"gitUrl"`
	GitBranch string `json:"gitBranch"`
}

func (a *App) encryptStoredValue(value string) string {
	if value == "" || !a.cfg.EnableEncryption {
		return value
	}
	ciphertext, err := security.EncryptString(a.cfg.EncryptionKey, value)
	if err != nil {
		return value
	}
	return ciphertext
}

func (a *App) decryptStoredValue(value string) string {
	if value == "" || !a.cfg.EnableEncryption {
		return value
	}
	plaintext, err := security.DecryptString(a.cfg.EncryptionKey, value)
	if err != nil {
		return value
	}
	return plaintext
}

func (a *App) generateProject(w http.ResponseWriter, r *http.Request) {
	var input generateProjectRequest

	if respond.DecodeJSON(r, &input) != nil {
		respond.Error(w, http.StatusBadRequest, "invalid request")
		return
	}

	input.ProjectName = strings.TrimSpace(input.ProjectName)
	input.Requirements = strings.TrimSpace(input.Requirements)
	input.FrontendTechnology = normalizeGeneratedTechnology(input.FrontendTechnology)
	input.BackendTechnology = normalizeGeneratedTechnology(input.BackendTechnology)
	input.DatabaseTechnology = normalizeGeneratedTechnology(input.DatabaseTechnology)
	input.GitURL = strings.TrimSpace(input.GitURL)
	input.GitBranch = strings.TrimSpace(input.GitBranch)

	if input.Requirements == "" {
		respond.Error(w, http.StatusBadRequest, "requirements are required")
		return
	}

	input = inferGeneratedProjectSettings(input)

	if input.ProjectName == "" {
		respond.Error(w, http.StatusBadRequest, "project name is required")
		return
	}

	if input.FrontendTechnology == "" {
		input.FrontendTechnology = "react"
	}

	if input.BackendTechnology == "" {
		input.BackendTechnology = "golang"
	}

	if input.DatabaseTechnology == "" {
		input.DatabaseTechnology = "postgresql"
	}

	if input.GitBranch == "" {
		input.GitBranch = "main"
	}

	if !validFrontendTechnology(input.FrontendTechnology) {
		respond.Error(w, http.StatusBadRequest, "unsupported frontend technology")
		return
	}

	if !validBackendTechnology(input.BackendTechnology) {
		respond.Error(w, http.StatusBadRequest, "unsupported backend technology")
		return
	}

	if !validDatabaseTechnology(input.DatabaseTechnology) {
		respond.Error(w, http.StatusBadRequest, "unsupported database technology")
		return
	}

	projectID := uuid.NewString()
	projectFolderName := sanitizeGeneratedProjectName(input.ProjectName)

	rootPath := filepath.Join(
		a.cfg.WorkDir,
		"generated-projects",
		projectFolderName+"-"+projectID,
	)

	if err := os.MkdirAll(rootPath, 0755); err != nil {
		respond.Error(
			w,
			http.StatusInternalServerError,
			"could not create generated project folder",
		)
		return
	}

	files, err := createGeneratedProjectFiles(
		rootPath,
		input,
	)

	if err != nil {
		fmt.Println("GENERATED PROJECT DATABASE ERROR:", err)

		_ = os.RemoveAll(rootPath)

		respond.Error(
			w,
			http.StatusInternalServerError,
			"could not save generated project",
		)
		return
	}

	if len(files) == 0 {
		_ = os.RemoveAll(rootPath)
		respond.Error(
			w,
			http.StatusInternalServerError,
			"no project files were generated",
		)
		return
	}

	_, err = a.db.Exec(
		r.Context(),
		`
		INSERT INTO generated_projects (
			id,
			owner_id,
			project_name,
			slug,
			description,
			requirements,
			frontend_technology,
			backend_technology,
			database_technology,
			git_url,
			git_branch,
			root_path
		)
		VALUES (
			$1,
			$2,
			$3,
			$4,
			$5,
			$6,
			$7,
			$8,
			$9,
			$10,
			$11,
			$12
		)
		`,
		projectID,
		auth.UserID(r.Context()),
		input.ProjectName,
		sanitizeGeneratedProjectName(input.ProjectName),
		a.encryptStoredValue(input.Requirements),
		a.encryptStoredValue(input.Requirements),
		input.FrontendTechnology,
		input.BackendTechnology,
		input.DatabaseTechnology,
		a.encryptStoredValue(input.GitURL),
		a.encryptStoredValue(input.GitBranch),
		rootPath,
	)

	if err != nil {
		_ = os.RemoveAll(rootPath)

		respond.Error(
			w,
			http.StatusInternalServerError,
			"could not save generated project",
		)
		return
	}

	for _, filePath := range files {
		fullPath := filepath.Join(rootPath, filePath)

		content, readErr := os.ReadFile(fullPath)
		if readErr != nil {
			continue
		}

		_, _ = a.db.Exec(
			r.Context(),
			`
			INSERT INTO generated_project_files (
				id,
				generated_project_id,
				file_name,
				file_path,
				file_content
			)
			VALUES ($1, $2, $3, $4, $5)
			ON CONFLICT (generated_project_id, file_path)
			DO UPDATE SET
				file_name = EXCLUDED.file_name,
				file_content = EXCLUDED.file_content,
				updated_at = NOW()
			`,
			uuid.NewString(),
			projectID,
			filepath.Base(filePath),
			filePath,
			a.encryptStoredValue(string(content)),
		)
	}

	respond.JSON(
		w,
		http.StatusCreated,
		map[string]any{
			"id":       projectID,
			"root":     rootPath,
			"rootPath": rootPath,
			"files":    files,
		},
	)
}

func (a *App) listGeneratedProjects(
	w http.ResponseWriter,
	r *http.Request,
) {
	rows, err := a.db.Query(
		r.Context(),
		`
		SELECT
			id,
			project_name,
			requirements,
			frontend_technology,
			backend_technology,
			database_technology,
			git_url,
			git_branch,
			root_path
		FROM generated_projects
		WHERE owner_id = $1
		ORDER BY created_at DESC
		`,
		auth.UserID(r.Context()),
	)

	if err != nil {
		respond.Error(
			w,
			http.StatusInternalServerError,
			"could not load generated projects",
		)
		return
	}

	defer rows.Close()

	result := []generatedProjectResponse{}

	for rows.Next() {
		var project generatedProjectResponse

		err := rows.Scan(
			&project.ID,
			&project.ProjectName,
			&project.Requirements,
			&project.FrontendTechnology,
			&project.BackendTechnology,
			&project.DatabaseTechnology,
			&project.GitURL,
			&project.GitBranch,
			&project.RootPath,
		)

		if err != nil {
			continue
		}

		project.Requirements = a.decryptStoredValue(project.Requirements)
		project.GitURL = a.decryptStoredValue(project.GitURL)
		project.GitBranch = a.decryptStoredValue(project.GitBranch)

		files, loadErr := a.loadGeneratedProjectFiles(r, project.ID)
		if loadErr == nil {
			project.Files = files
		} else {
			project.Files = []string{}
		}

		result = append(result, project)
	}

	respond.JSON(w, http.StatusOK, result)
}

func (a *App) getGeneratedProject(
	w http.ResponseWriter,
	r *http.Request,
) {
	projectID := chi.URLParam(r, "generatedProjectID")

	var project generatedProjectResponse

	err := a.db.QueryRow(
		r.Context(),
		`
		SELECT
			id,
			project_name,
			requirements,
			frontend_technology,
			backend_technology,
			database_technology,
			git_url,
			git_branch,
			root_path
		FROM generated_projects
		WHERE id = $1
		AND owner_id = $2
		`,
		projectID,
		auth.UserID(r.Context()),
	).Scan(
		&project.ID,
		&project.ProjectName,
		&project.Requirements,
		&project.FrontendTechnology,
		&project.BackendTechnology,
		&project.DatabaseTechnology,
		&project.GitURL,
		&project.GitBranch,
		&project.RootPath,
	)

	if err != nil {
		respond.Error(
			w,
			http.StatusNotFound,
			"generated project not found",
		)
		return
	}

	files, err := a.loadGeneratedProjectFiles(
		r,
		projectID,
	)

	if err != nil {
		respond.Error(
			w,
			http.StatusInternalServerError,
			"could not load generated project files",
		)
		return
	}

	project.Requirements = a.decryptStoredValue(project.Requirements)
	project.GitURL = a.decryptStoredValue(project.GitURL)
	project.GitBranch = a.decryptStoredValue(project.GitBranch)
	project.Files = files

	respond.JSON(w, http.StatusOK, project)
}

func (a *App) getGeneratedProjectFile(
	w http.ResponseWriter,
	r *http.Request,
) {
	projectID := chi.URLParam(r, "generatedProjectID")
	filePath := strings.TrimSpace(r.URL.Query().Get("path"))

	if filePath == "" {
		respond.Error(
			w,
			http.StatusBadRequest,
			"file path is required",
		)
		return
	}

	var content string

	err := a.db.QueryRow(
		r.Context(),
		`
		SELECT f.file_content
		FROM generated_project_files f
		JOIN generated_projects p
			ON p.id = f.generated_project_id
		WHERE f.generated_project_id = $1
		AND f.file_path = $2
		AND p.owner_id = $3
		`,
		projectID,
		filePath,
		auth.UserID(r.Context()),
	).Scan(&content)

	if err != nil {
		var rootPath string
		rootErr := a.db.QueryRow(
			r.Context(),
			`
			SELECT root_path
			FROM generated_projects
			WHERE id = $1
			AND owner_id = $2
			`,
			projectID,
			auth.UserID(r.Context()),
		).Scan(&rootPath)
		if rootErr == nil {
			fullPath, safeErr := safeGeneratedFilePath(rootPath, filePath)
			if safeErr == nil {
				if diskContent, readErr := os.ReadFile(fullPath); readErr == nil {
					content = string(diskContent)
					respond.JSON(w, http.StatusOK, map[string]string{"content": a.decryptStoredValue(content)})
					return
				}
			}
		}
		respond.Error(
			w,
			http.StatusNotFound,
			"generated file not found",
		)
		return
	}

	respond.JSON(
		w,
		http.StatusOK,
		map[string]string{
			"content": a.decryptStoredValue(content),
		},
	)
}

func (a *App) saveGeneratedProjectFile(
	w http.ResponseWriter,
	r *http.Request,
) {
	projectID := chi.URLParam(r, "generatedProjectID")

	var input saveGeneratedFileRequest

	if respond.DecodeJSON(r, &input) != nil {
		respond.Error(
			w,
			http.StatusBadRequest,
			"invalid request",
		)
		return
	}

	input.Path = strings.TrimSpace(input.Path)

	if input.Path == "" {
		respond.Error(
			w,
			http.StatusBadRequest,
			"file path is required",
		)
		return
	}

	var rootPath string

	err := a.db.QueryRow(
		r.Context(),
		`
		SELECT root_path
		FROM generated_projects
		WHERE id = $1
		AND owner_id = $2
		`,
		projectID,
		auth.UserID(r.Context()),
	).Scan(&rootPath)

	if err != nil {
		respond.Error(
			w,
			http.StatusNotFound,
			"generated project not found",
		)
		return
	}

	fullPath, err := safeGeneratedFilePath(
		rootPath,
		input.Path,
	)

	if err != nil {
		respond.Error(
			w,
			http.StatusBadRequest,
			"invalid file path",
		)
		return
	}

	if err := os.WriteFile(
		fullPath,
		[]byte(input.Content),
		0644,
	); err != nil {
		respond.Error(
			w,
			http.StatusInternalServerError,
			"could not save generated file",
		)
		return
	}

	tag, err := a.db.Exec(
		r.Context(),
		`
		UPDATE generated_project_files
		SET
			file_content = $1,
			updated_at = NOW()
		WHERE generated_project_id = $2
		AND file_path = $3
		`,
		a.encryptStoredValue(input.Content),
		projectID,
		input.Path,
	)

	if err != nil {
		respond.Error(
			w,
			http.StatusInternalServerError,
			"could not update generated file",
		)
		return
	}

	if tag.RowsAffected() == 0 {
		_, insertErr := a.db.Exec(
			r.Context(),
			`
			INSERT INTO generated_project_files (
				id,
				generated_project_id,
				file_name,
				file_path,
				file_content
			)
			VALUES ($1, $2, $3, $4, $5)
			ON CONFLICT (generated_project_id, file_path)
			DO UPDATE SET
				file_name = EXCLUDED.file_name,
				file_content = EXCLUDED.file_content,
				updated_at = NOW()
			`,
			uuid.NewString(),
			projectID,
			filepath.Base(input.Path),
			input.Path,
			a.encryptStoredValue(input.Content),
		)
		if insertErr != nil {
			respond.Error(
				w,
				http.StatusInternalServerError,
				"could not update generated file",
			)
			return
		}
	}

	respond.JSON(
		w,
		http.StatusOK,
		map[string]string{
			"status": "saved",
		},
	)
}

func (a *App) syncGeneratedProjectGit(
	w http.ResponseWriter,
	r *http.Request,
) {
	projectID := chi.URLParam(r, "generatedProjectID")

	var input syncGeneratedGitRequest

	if respond.DecodeJSON(r, &input) != nil {
		respond.Error(
			w,
			http.StatusBadRequest,
			"invalid request",
		)
		return
	}

	input.GitURL = strings.TrimSpace(input.GitURL)
	input.GitBranch = strings.TrimSpace(input.GitBranch)

	if input.GitURL == "" {
		respond.Error(
			w,
			http.StatusBadRequest,
			"git URL is required",
		)
		return
	}

	if input.GitBranch == "" {
		input.GitBranch = "main"
	}

	var rootPath string

	err := a.db.QueryRow(
		r.Context(),
		`
		SELECT root_path
		FROM generated_projects
		WHERE id = $1
		AND owner_id = $2
		`,
		projectID,
		auth.UserID(r.Context()),
	).Scan(&rootPath)

	if err != nil {
		respond.Error(
			w,
			http.StatusNotFound,
			"generated project not found",
		)
		return
	}

	if err := runGeneratedGit(
		rootPath,
		"init",
	); err != nil {
		respond.Error(
			w,
			http.StatusInternalServerError,
			err.Error(),
		)
		return
	}

	_ = runGeneratedGit(
		rootPath,
		"remote",
		"remove",
		"origin",
	)

	if err := runGeneratedGit(
		rootPath,
		"remote",
		"add",
		"origin",
		input.GitURL,
	); err != nil {
		respond.Error(
			w,
			http.StatusInternalServerError,
			err.Error(),
		)
		return
	}

	if err := runGeneratedGit(
		rootPath,
		"checkout",
		"-B",
		input.GitBranch,
	); err != nil {
		respond.Error(
			w,
			http.StatusInternalServerError,
			err.Error(),
		)
		return
	}

	if err := runGeneratedGit(
		rootPath,
		"add",
		".",
	); err != nil {
		respond.Error(
			w,
			http.StatusInternalServerError,
			err.Error(),
		)
		return
	}

	_ = runGeneratedGit(
		rootPath,
		"commit",
		"-m",
		"Generated project",
	)

	if err := runGeneratedGit(
		rootPath,
		"push",
		"-u",
		"origin",
		input.GitBranch,
	); err != nil {
		respond.Error(
			w,
			http.StatusInternalServerError,
			err.Error(),
		)
		return
	}

	_, _ = a.db.Exec(
		r.Context(),
		`
		UPDATE generated_projects
		SET
			git_url = $1,
			git_branch = $2,
			updated_at = NOW()
		WHERE id = $3
		AND owner_id = $4
		`,
		a.encryptStoredValue(input.GitURL),
		a.encryptStoredValue(input.GitBranch),
		projectID,
		auth.UserID(r.Context()),
	)

	respond.JSON(
		w,
		http.StatusOK,
		map[string]string{
			"status": "pushed",
		},
	)
}

func (a *App) loadGeneratedProjectFiles(
	r *http.Request,
	projectID string,
) ([]string, error) {
	rows, err := a.db.Query(
		r.Context(),
		`
		SELECT file_path
		FROM generated_project_files
		WHERE generated_project_id = $1
		ORDER BY file_path
		`,
		projectID,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	files := []string{}

	for rows.Next() {
		var filePath string

		if err := rows.Scan(&filePath); err != nil {
			return nil, err
		}

		files = append(files, filePath)
	}

	if len(files) > 0 || rows.Err() != nil {
		return files, rows.Err()
	}

	var rootPath string

	projectRowErr := a.db.QueryRow(
		r.Context(),
		`
		SELECT root_path
		FROM generated_projects
		WHERE id = $1
		AND owner_id = $2
		`,
		projectID,
		auth.UserID(r.Context()),
	).Scan(&rootPath)

	if projectRowErr != nil {
		return files, rows.Err()
	}

	return collectGeneratedProjectFilesFromDisk(rootPath)
}

func createGeneratedProjectFiles(
	rootPath string,
	input generateProjectRequest,
) ([]string, error) {
	files := map[string]string{}

	files["README.md"] = createGeneratedReadme(input)

	switch input.FrontendTechnology {
	case "react":
		addGeneratedReactFiles(files, input)
	case "angular":
		addGeneratedAngularFiles(files, input)
	case "vue":
		addGeneratedVueFiles(files, input)
	case "nextjs":
		addGeneratedNextJSFiles(files, input)
	}

	switch input.BackendTechnology {
	case "golang":
		addGeneratedGoFiles(files, input)
	case "java":
		addGeneratedJavaFiles(files, input)
	case "node":
		addGeneratedNodeFiles(files, input)
	case "python":
		addGeneratedPythonFiles(files, input)
	case "dotnet":
		addGeneratedDotNetFiles(files, input)
	}

	addGeneratedDatabaseFiles(files, input)
	addGeneratedDockerFiles(files, input)
	addGeneratedKubernetesFiles(files, input)

	filePaths := make([]string, 0, len(files))

	for relativePath, content := range files {
		fullPath, err := safeGeneratedFilePath(
			rootPath,
			relativePath,
		)

		if err != nil {
			return nil, err
		}

		if err := os.MkdirAll(
			filepath.Dir(fullPath),
			0755,
		); err != nil {
			return nil, err
		}

		if err := os.WriteFile(
			fullPath,
			[]byte(content),
			0644,
		); err != nil {
			return nil, err
		}

		filePaths = append(
			filePaths,
			relativePath,
		)
	}

	return filePaths, nil
}

func createGeneratedReadme(
	input generateProjectRequest,
) string {
	return fmt.Sprintf(
		"# %s\n\nRequirements\n\n%s\n\nTechnology Stack\n\nFrontend: %s\nBackend: %s\nDatabase: %s\n",
		input.ProjectName,
		input.Requirements,
		input.FrontendTechnology,
		input.BackendTechnology,
		input.DatabaseTechnology,
	)
}

func inferGeneratedProjectSettings(input generateProjectRequest) generateProjectRequest {
	prompt := strings.ToLower(input.Requirements)

	if input.FrontendTechnology == "" || isDefaultTechnology(input.FrontendTechnology, "react") && strings.Contains(prompt, "angular") {
		input.FrontendTechnology = inferFrontendTechnologyFromPrompt(prompt)
	}

	if input.BackendTechnology == "" || isDefaultTechnology(input.BackendTechnology, "golang") && (strings.Contains(prompt, "java") || strings.Contains(prompt, "spring")) {
		input.BackendTechnology = inferBackendTechnologyFromPrompt(prompt)
	}

	if input.DatabaseTechnology == "" || isDefaultTechnology(input.DatabaseTechnology, "postgresql") && (strings.Contains(prompt, "mysql") || strings.Contains(prompt, "mongodb") || strings.Contains(prompt, "mongo") || strings.Contains(prompt, "postgres")) {
		input.DatabaseTechnology = inferDatabaseTechnologyFromPrompt(prompt)
	}

	if input.ProjectName == "" || isDefaultGeneratedProjectName(input.ProjectName) {
		if inferredName := inferProjectNameFromPrompt(prompt); inferredName != "" {
			input.ProjectName = inferredName
		}
	}

	return input
}

func inferFrontendTechnologyFromPrompt(prompt string) string {
	switch {
	case strings.Contains(prompt, "angular"):
		return "angular"
	case strings.Contains(prompt, "vue") || strings.Contains(prompt, "vuejs"):
		return "vue"
	case strings.Contains(prompt, "next") || strings.Contains(prompt, "nextjs"):
		return "nextjs"
	case strings.Contains(prompt, "react") || strings.Contains(prompt, "frontend"):
		return "react"
	default:
		return "react"
	}
}

func inferBackendTechnologyFromPrompt(prompt string) string {
	switch {
	case strings.Contains(prompt, "java") || strings.Contains(prompt, "spring"):
		return "java"
	case strings.Contains(prompt, "python") || strings.Contains(prompt, "fastapi") || strings.Contains(prompt, "django"):
		return "python"
	case strings.Contains(prompt, "node") || strings.Contains(prompt, "express") || strings.Contains(prompt, "javascript"):
		return "node"
	case strings.Contains(prompt, "golang") || strings.Contains(prompt, "go ") || strings.Contains(prompt, "go-") || strings.Contains(prompt, "go language"):
		return "golang"
	case strings.Contains(prompt, "backend"):
		return "golang"
	default:
		return "golang"
	}
}

func inferDatabaseTechnologyFromPrompt(prompt string) string {
	switch {
	case strings.Contains(prompt, "mongodb") || strings.Contains(prompt, "mongo db") || strings.Contains(prompt, "mongo"):
		return "mongodb"
	case strings.Contains(prompt, "mysql"):
		return "mysql"
	case strings.Contains(prompt, "sqlite"):
		return "sqlite"
	case strings.Contains(prompt, "postgres") || strings.Contains(prompt, "postgresql") || strings.Contains(prompt, "sql"):
		return "postgresql"
	default:
		return "postgresql"
	}
}

func inferProjectNameFromPrompt(prompt string) string {
	if prompt == "" {
		return ""
	}

	for _, keyword := range []string{"flower shop", "flower boutique", "cloth shop", "clothing shop", "clothing store", "fashion store", "fashion boutique", "coffee shop", "bakery shop", "grocery store", "book store", "restaurant"} {
		if strings.Contains(prompt, keyword) {
			return sanitizeGeneratedProjectName(strings.ReplaceAll(keyword, " ", "-"))
		}
	}

	normalized := strings.NewReplacer(
		"/", " ",
		"-", " ",
		"_", " ",
		".", " ",
		",", " ",
		"!", " ",
		"?", " ",
		";", " ",
		":", " ",
	).Replace(prompt)

	words := strings.Fields(normalized)
	for i := range words {
		if words[i] == "shop" || words[i] == "store" || words[i] == "boutique" || words[i] == "website" || words[i] == "app" || words[i] == "application" || words[i] == "dashboard" || words[i] == "portal" || words[i] == "platform" {
			start := i - 1
			if start < 0 {
				start = 0
			}
			end := i
			if end-start > 4 {
				end = start + 4
			}
			candidate := strings.Join(words[start:end], " ")
			if candidate != "" {
				return sanitizeGeneratedProjectName(candidate)
			}
		}
	}

	return ""
}

func isDefaultGeneratedProjectName(name string) bool {
	value := strings.TrimSpace(strings.ToLower(name))
	return value == "" || value == "coffee-shop" || value == "generated-project" || value == "generated-app"
}

func isDefaultTechnology(value string, defaultValue string) bool {
	return strings.TrimSpace(strings.ToLower(value)) == defaultValue
}

func addGeneratedReactFiles(
	files map[string]string,
	input generateProjectRequest,
) {
	files["frontend/package.json"] = fmt.Sprintf(
		`{
  "name": "%s-frontend",
  "version": "1.0.0",
  "private": true,
  "dependencies": {
    "react": "^18.3.1",
    "react-dom": "^18.3.1",
    "react-router-dom": "^6.28.0",
    "react-scripts": "5.0.1"
  },
  "scripts": {
    "start": "react-scripts start",
    "build": "react-scripts build"
  }
}`,
		sanitizeGeneratedProjectName(input.ProjectName),
	)

	files["frontend/public/index.html"] = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Generated Application</title>
</head>
<body>
  <div id="root"></div>
</body>
</html>`

	files["frontend/src/index.js"] = `import React from "react";
import ReactDOM from "react-dom/client";
import { BrowserRouter } from "react-router-dom";
import App from "./App";
import "./App.css";

const root = ReactDOM.createRoot(document.getElementById("root"));

root.render(
  <React.StrictMode>
    <BrowserRouter>
      <App />
    </BrowserRouter>
  </React.StrictMode>
);`

	files["frontend/src/App.js"] = `import { Routes, Route, Navigate } from "react-router-dom";
import LoginPage from "./pages/LoginPage";
import RegisterPage from "./pages/RegisterPage";
import DashboardPage from "./pages/DashboardPage";
import CatalogPage from "./pages/CatalogPage";
import CartPage from "./pages/CartPage";
import CheckoutPage from "./pages/CheckoutPage";
import AdminPage from "./pages/AdminPage";

function App() {
  return (
    <div className="app-shell">
      <Routes>
        <Route path="/" element={<Navigate to="/dashboard" replace />} />
        <Route path="/login" element={<LoginPage />} />
        <Route path="/register" element={<RegisterPage />} />
        <Route path="/dashboard" element={<DashboardPage />} />
        <Route path="/catalog" element={<CatalogPage />} />
        <Route path="/cart" element={<CartPage />} />
        <Route path="/checkout" element={<CheckoutPage />} />
        <Route path="/admin" element={<AdminPage />} />
      </Routes>
    </div>
  );
}

export default App;
`

	files["frontend/src/App.css"] = `:root {
  --bg: #f4f7fb;
  --panel: #ffffff;
  --primary: #0f74ff;
  --primary-dark: #0a54d8;
  --text: #1c2430;
  --muted: #5f6f86;
  --border: #dfe7f1;
  --success: #1e9a68;
  --warning: #ffb703;
}

* { box-sizing: border-box; }
body {
  margin: 0;
  font-family: Arial, sans-serif;
  background: var(--bg);
  color: var(--text);
}

a { text-decoration: none; }

.app-shell {
  min-height: 100vh;
}

.page {
  max-width: 1200px;
  margin: 0 auto;
  padding: 32px 24px 60px;
}

.navbar,
.card,
.table-card,
.form-card {
  background: var(--panel);
  border: 1px solid var(--border);
  border-radius: 16px;
  box-shadow: 0 8px 24px rgba(15, 23, 42, 0.04);
}

.navbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px;
  margin-bottom: 24px;
}

.nav-links {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
}

.nav-links a,
.button,
.link-button {
  border: none;
  background: var(--primary);
  color: white;
  border-radius: 10px;
  padding: 10px 14px;
  cursor: pointer;
  font-weight: 600;
}

.secondary {
  background: #edf3ff;
  color: var(--primary-dark);
}

.grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
  gap: 16px;
}

.card { padding: 20px; }

.form-card {
  max-width: 480px;
  margin: 40px auto;
  padding: 24px;
}

.input, select, textarea {
  width: 100%;
  padding: 12px 14px;
  border: 1px solid var(--border);
  border-radius: 10px;
  margin-top: 8px;
  margin-bottom: 16px;
  font-size: 14px;
}

.hero {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 16px;
  margin-bottom: 24px;
}

.price {
  font-weight: bold;
  color: var(--primary-dark);
}

.product-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
  gap: 18px;
}

.product-card {
  background: white;
  border: 1px solid var(--border);
  border-radius: 16px;
  padding: 16px;
}

.row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 12px;
}

.table-card {
  padding: 16px;
}

table {
  width: 100%;
  border-collapse: collapse;
}

th, td {
  padding: 12px;
  border-bottom: 1px solid var(--border);
  text-align: left;
}
`

	files["frontend/src/pages/LoginPage.jsx"] = `const LoginPage = () => (
  <div className="page">
    <div className="form-card">
      <h2>Login</h2>
      <label>Email
        <input className="input" type="email" placeholder="admin@example.com" />
      </label>
      <label>Password
        <input className="input" type="password" placeholder="••••••••" />
      </label>
      <button className="button" type="button">Sign in</button>
    </div>
  </div>
);

export default LoginPage;`

	files["frontend/src/pages/RegisterPage.jsx"] = `const RegisterPage = () => (
  <div className="page">
    <div className="form-card">
      <h2>Create account</h2>
      <label>Name
        <input className="input" type="text" placeholder="Full name" />
      </label>
      <label>Email
        <input className="input" type="email" placeholder="you@example.com" />
      </label>
      <label>Password
        <input className="input" type="password" placeholder="••••••••" />
      </label>
      <button className="button" type="button">Register</button>
    </div>
  </div>
);

export default RegisterPage;`

	files["frontend/src/pages/DashboardPage.jsx"] = `const stats = [
  { label: "Revenue", value: "$24.8k" },
  { label: "Orders", value: "324" },
  { label: "Customers", value: "1.2k" },
  { label: "Conversion", value: "6.4%" },
];

const DashboardPage = () => (
  <div className="page">
    <nav className="navbar">
      <strong>Dashboard</strong>
      <div className="nav-links">
        <a href="/catalog">Catalog</a>
        <a href="/cart">Cart</a>
        <a href="/admin">Admin</a>
      </div>
    </nav>

    <div className="grid">
      {stats.map((item) => (
        <div className="card" key={item.label}>
          <div>{item.label}</div>
          <h2>{item.value}</h2>
        </div>
      ))}
    </div>
  </div>
);

export default DashboardPage;`

	files["frontend/src/pages/CatalogPage.jsx"] = `const products = [
  { id: 1, name: "Premium Desk", price: "$299" },
  { id: 2, name: "Leather Chair", price: "$149" },
  { id: 3, name: "Office Lamp", price: "$89" },
];

const CatalogPage = () => (
  <div className="page">
    <div className="hero">
      <h2>Product catalog</h2>
      <button className="button" type="button">Add product</button>
    </div>

    <div className="product-grid">
      {products.map((product) => (
        <div key={product.id} className="product-card">
          <h3>{product.name}</h3>
          <p className="price">{product.price}</p>
          <button className="button secondary" type="button">Add to cart</button>
        </div>
      ))}
    </div>
  </div>
);

export default CatalogPage;`

	files["frontend/src/pages/CartPage.jsx"] = `const items = [
  { name: "Premium Desk", qty: 1, total: "$299" },
  { name: "Office Lamp", qty: 2, total: "$178" },
];

const CartPage = () => (
  <div className="page">
    <h2>Cart</h2>
    <div className="table-card">
      <table>
        <thead>
          <tr><th>Item</th><th>Qty</th><th>Total</th></tr>
        </thead>
        <tbody>
          {items.map((item) => (
            <tr key={item.name}>
              <td>{item.name}</td>
              <td>{item.qty}</td>
              <td>{item.total}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  </div>
);

export default CartPage;`

	files["frontend/src/pages/CheckoutPage.jsx"] = `const CheckoutPage = () => (
  <div className="page">
    <div className="form-card">
      <h2>Checkout</h2>
      <label>Card number
        <input className="input" type="text" placeholder="4242 4242 4242 4242" />
      </label>
      <label>Name on card
        <input className="input" type="text" placeholder="John Doe" />
      </label>
      <button className="button" type="button">Pay now</button>
    </div>
  </div>
);

export default CheckoutPage;`

	files["frontend/src/pages/AdminPage.jsx"] = `const AdminPage = () => (
  <div className="page">
    <h2>Admin</h2>
    <div className="card">
      <h3>Upload product</h3>
      <input className="input" type="text" placeholder="Product name" />
      <input className="input" type="number" placeholder="Price" />
      <button className="button" type="button">Save</button>
    </div>
  </div>
);

export default AdminPage;`
}

func addGeneratedAngularFiles(
	files map[string]string,
	input generateProjectRequest,
) {
	files["frontend/package.json"] = `{
  "name": "generated-angular-application",
  "version": "1.0.0",
  "scripts": {
    "start": "ng serve",
    "build": "ng build"
  },
  "dependencies": {
    "@angular/common": "^19.0.0",
    "@angular/compiler": "^19.0.0",
    "@angular/core": "^19.0.0",
    "@angular/platform-browser": "^19.0.0",
    "@angular/platform-browser-dynamic": "^19.0.0",
    "rxjs": "^7.8.1",
    "zone.js": "^0.15.0"
  },
  "devDependencies": {
    "@angular/cli": "^19.0.0",
    "@angular/compiler-cli": "^19.0.0",
    "typescript": "^5.6.0"
  }
}`

	files["frontend/src/app/app.component.ts"] = fmt.Sprintf(
		`import { Component } from "@angular/core";

@Component({
  selector: "app-root",
  standalone: true,
  template: `+"`"+`
    <main>
      <h1>%s</h1>
      <p>%s</p>
    </main>
  `+"`"+`
})
export class AppComponent {}
`,
		escapeGeneratedJS(input.ProjectName),
		escapeGeneratedJS(input.Requirements),
	)

	files["frontend/src/main.ts"] = `import { bootstrapApplication } from "@angular/platform-browser";
import { AppComponent } from "./app/app.component";

bootstrapApplication(AppComponent)
  .catch((error) => console.error(error));`
}

func addGeneratedVueFiles(
	files map[string]string,
	input generateProjectRequest,
) {
	files["frontend/package.json"] = `{
  "name": "generated-vue-app",
  "version": "1.0.0",
  "private": true,
  "scripts": {
    "dev": "vite",
    "build": "vite build"
  },
  "dependencies": {
    "vue": "^3.5.0"
  },
  "devDependencies": {
    "vite": "^5.4.0"
  }
}`

	files["frontend/index.html"] = `<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>Generated Vue App</title>
  </head>
  <body>
    <div id="app"></div>
    <script type="module" src="/src/main.js"></script>
  </body>
</html>`

	files["frontend/src/main.js"] = `import { createApp } from "vue";
import App from "./App.vue";
import "./style.css";

createApp(App).mount("#app");`

	files["frontend/src/App.vue"] = fmt.Sprintf(
		`<template>
  <main class="app-shell">
    <h1>%s</h1>
    <p>%s</p>
  </main>
</template>
`,
		escapeGeneratedJS(input.ProjectName),
		escapeGeneratedJS(input.Requirements),
	)

	files["frontend/src/style.css"] = `body { margin: 0; font-family: Arial, sans-serif; background: #f4f7fb; }
.app-shell { max-width: 960px; margin: 0 auto; padding: 80px 24px; color: #1f2937; }
.app-shell h1 { margin-bottom: 12px; }`
}

func addGeneratedNextJSFiles(
	files map[string]string,
	input generateProjectRequest,
) {
	files["frontend/package.json"] = `{
  "name": "generated-nextjs-app",
  "private": true,
  "scripts": {
    "dev": "next dev",
    "build": "next build",
    "start": "next start"
  },
  "dependencies": {
    "next": "14.2.15",
    "react": "18.3.1",
    "react-dom": "18.3.1"
  }
}`

	files["frontend/app/page.js"] = fmt.Sprintf(
		`export default function Page() {
  return (
    <main style={{ padding: "40px", fontFamily: "Arial, sans-serif" }}>
      <h1>%s</h1>
      <p>%s</p>
    </main>
  );
}
`,
		escapeGeneratedJS(input.ProjectName),
		escapeGeneratedJS(input.Requirements),
	)

	files["frontend/app/layout.js"] = `export const metadata = { title: "Generated App" };

export default function RootLayout({ children }) {
  return (
    <html lang="en">
      <body>{children}</body>
    </html>
  );
}`
}

func addGeneratedNodeFiles(
	files map[string]string,
	input generateProjectRequest,
) {
	files["backend/package.json"] = `{
  "name": "generated-node-api",
  "version": "1.0.0",
  "main": "server.js",
  "scripts": {
    "start": "node server.js"
  },
  "dependencies": {
    "express": "^4.19.2"
  }
}`

	files["backend/server.js"] = fmt.Sprintf(
		`const express = require("express");
const app = express();

app.use(express.json());

app.get("/api/health", (req, res) => {
  res.json({ status: "ok", app: "%s" });
});

app.listen(8080, () => {
  console.log("API running on http://localhost:8080");
});
`,
		escapeGeneratedJS(input.ProjectName),
	)
}

func addGeneratedPythonFiles(
	files map[string]string,
	input generateProjectRequest,
) {
	files["backend/requirements.txt"] = "fastapi==0.115.0\nuvicorn==0.30.6\n"
	files["backend/main.py"] = fmt.Sprintf(
		`from fastapi import FastAPI

app = FastAPI(title="%s")

@app.get("/api/health")
def health():
    return {"status": "ok", "app": "%s"}
`,
		escapeGeneratedJS(input.ProjectName),
		escapeGeneratedJS(input.ProjectName),
	)
}

func addGeneratedDotNetFiles(
	files map[string]string,
	input generateProjectRequest,
) {
	files["backend/generated-app.csproj"] = `
<Project Sdk="Microsoft.NET.Sdk.Web">
  <PropertyGroup>
    <TargetFramework>net8.0</TargetFramework>
    <ImplicitUsings>enable</ImplicitUsings>
    <Nullable>enable</Nullable>
  </PropertyGroup>
</Project>`

	files["backend/Program.cs"] = fmt.Sprintf(
		`var builder = WebApplication.CreateBuilder(args);
var app = builder.Build();

app.MapGet("/api/health", () => new { status = "ok", app = "%s" });

app.Run("http://0.0.0.0:8080");
`,
		escapeGeneratedJS(input.ProjectName),
	)
}

func addGeneratedGoFiles(
	files map[string]string,
	input generateProjectRequest,
) {
	projectModule := sanitizeGeneratedProjectName(
		input.ProjectName,
	)

	files["backend/go.mod"] = fmt.Sprintf(
		"module %s\n\ngo 1.23\n\nrequire (\n\tgithub.com/go-chi/chi/v5 v5.1.0\n\tgithub.com/google/uuid v1.6.0\n)\n",
		projectModule,
	)

	files["backend/cmd/api/main.go"] = "package main\n\nimport (\n\t\"log\"\n\t\"net/http\"\n\n\t\"github.com/go-chi/chi/v5\"\n)\n\nfunc main() {\n\tr := chi.NewRouter()\n\tr.Get(\"/api/health\", func(w http.ResponseWriter, r *http.Request) {\n\t\tw.Header().Set(\"Content-Type\", \"application/json\")\n\t\t_, _ = w.Write([]byte(\"{\\\"status\\\":\\\"ok\\\"}\"))\n\t})\n\n\tr.Route(\"/api\", func(r chi.Router) {\n\t\tr.Post(\"/auth/register\", func(w http.ResponseWriter, r *http.Request) {\n\t\t\tw.WriteHeader(http.StatusCreated)\n\t\t\t_, _ = w.Write([]byte(\"{\\\"status\\\":\\\"created\\\"}\"))\n\t\t})\n\t\tr.Post(\"/auth/login\", func(w http.ResponseWriter, r *http.Request) {\n\t\t\t_, _ = w.Write([]byte(\"{\\\"status\\\":\\\"logged-in\\\"}\"))\n\t\t})\n\t\tr.Get(\"/products\", func(w http.ResponseWriter, r *http.Request) {\n\t\t\t_, _ = w.Write([]byte(\"[{\\\"id\\\":\\\"1\\\",\\\"name\\\":\\\"Premium Desk\\\",\\\"price\\\":299}]\"))\n\t\t})\n\t\tr.Get(\"/dashboard/summary\", func(w http.ResponseWriter, r *http.Request) {\n\t\t\t_, _ = w.Write([]byte(\"{\\\"revenue\\\":24800,\\\"orders\\\":324,\\\"customers\\\":1200}\"))\n\t\t})\n\t})\n\n\tlog.Println(\"API running on http://localhost:8080\")\n\tlog.Fatal(http.ListenAndServe(\":8080\", r))\n}\n"

	files["backend/internal/models/models.go"] = "package models\n\ntype User struct {\n\tID    string `json:\"id\"`\n\tName  string `json:\"name\"`\n\tEmail string `json:\"email\"`\n}\n\ntype Product struct {\n\tID    string `json:\"id\"`\n\tName  string `json:\"name\"`\n\tPrice float64 `json:\"price\"`\n}\n\ntype Order struct {\n\tID      string `json:\"id\"`\n\tUserID  string `json:\"userId\"`\n\tTotal   float64 `json:\"total\"`\n\tStatus  string `json:\"status\"`\n}\n"

	files["backend/internal/auth/auth.go"] = "package auth\n\nimport \"fmt\"\n\nfunc GenerateToken(email string) string {\n\treturn fmt.Sprintf(\"token-%s\", email)\n}\n"

	files["backend/internal/handlers/products.go"] = "package handlers\n\nimport (\n\t\"encoding/json\"\n\t\"net/http\"\n)\n\ntype ProductResponse struct {\n\tID    string `json:\"id\"`\n\tName  string `json:\"name\"`\n\tPrice float64 `json:\"price\"`\n}\n\nfunc ListProducts(w http.ResponseWriter, r *http.Request) {\n\tproducts := []ProductResponse{{ID: \"1\", Name: \"Premium Desk\", Price: 299}}\n\t_ = json.NewEncoder(w).Encode(products)\n}\n"

}

func addGeneratedJavaFiles(
	files map[string]string,
	input generateProjectRequest,
) {
	artifact := sanitizeGeneratedProjectName(
		input.ProjectName,
	)

	files["backend/pom.xml"] = fmt.Sprintf(
		`<?xml version="1.0" encoding="UTF-8"?>
<project xmlns="http://maven.apache.org/POM/4.0.0"
 xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
 xsi:schemaLocation="
 http://maven.apache.org/POM/4.0.0
 https://maven.apache.org/xsd/maven-4.0.0.xsd">

  <modelVersion>4.0.0</modelVersion>

  <parent>
    <groupId>org.springframework.boot</groupId>
    <artifactId>spring-boot-starter-parent</artifactId>
    <version>3.4.1</version>
  </parent>

  <groupId>com.generated</groupId>
  <artifactId>%s</artifactId>
  <version>1.0.0</version>

  <properties>
    <java.version>21</java.version>
  </properties>

  <dependencies>
    <dependency>
      <groupId>org.springframework.boot</groupId>
      <artifactId>spring-boot-starter-web</artifactId>
    </dependency>
  </dependencies>

</project>`,
		artifact,
	)

	files["backend/src/main/java/com/generated/Application.java"] = `package com.generated;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;

@SpringBootApplication
public class Application {

    public static void main(String[] args) {
        SpringApplication.run(
            Application.class,
            args
        );
    }
}`
}

func addGeneratedDatabaseFiles(
	files map[string]string,
	input generateProjectRequest,
) {
	switch input.DatabaseTechnology {
	case "postgresql":
		files["database/schema.sql"] = `CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY,
    name VARCHAR(200) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    role VARCHAR(50) DEFAULT 'customer',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS products (
    id UUID PRIMARY KEY,
    name VARCHAR(200) NOT NULL,
    description TEXT,
    price NUMERIC(10,2) NOT NULL,
    stock INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS orders (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id),
    total NUMERIC(10,2) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS order_items (
    id UUID PRIMARY KEY,
    order_id UUID NOT NULL REFERENCES orders(id),
    product_id UUID NOT NULL REFERENCES products(id),
    quantity INTEGER NOT NULL DEFAULT 1,
    unit_price NUMERIC(10,2) NOT NULL
);

CREATE TABLE IF NOT EXISTS carts (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id),
    product_id UUID NOT NULL REFERENCES products(id),
    quantity INTEGER NOT NULL DEFAULT 1
);`

	case "mysql":
		files["database/schema.sql"] = `CREATE TABLE IF NOT EXISTS users (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(200) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(50) DEFAULT 'customer',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS products (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(200) NOT NULL,
    description TEXT,
    price DECIMAL(10,2) NOT NULL,
    stock INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS orders (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL,
    total DECIMAL(10,2) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS order_items (
    id VARCHAR(36) PRIMARY KEY,
    order_id VARCHAR(36) NOT NULL,
    product_id VARCHAR(36) NOT NULL,
    quantity INT NOT NULL DEFAULT 1,
    unit_price DECIMAL(10,2) NOT NULL
);`

	case "mongodb":
		files["database/README.md"] = `Database: MongoDB

Collections:
- users
- products
- orders
- carts

Configure a MongoDB connection string through an environment variable.
`
	case "sqlite":
		files["database/schema.sql"] = `CREATE TABLE IF NOT EXISTS users (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    role TEXT DEFAULT 'customer',
    created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS products (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT,
    price REAL NOT NULL,
    stock INTEGER NOT NULL DEFAULT 0,
    created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS orders (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    total REAL NOT NULL,
    status TEXT NOT NULL DEFAULT 'pending',
    created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);`
	case "none":
		files["database/README.md"] = `Database: none

No database is configured for this project.
`
	}
}

func addGeneratedDockerFiles(
	files map[string]string,
	input generateProjectRequest,
) {
	if input.FrontendTechnology == "react" || input.FrontendTechnology == "vue" || input.FrontendTechnology == "nextjs" {
		files["frontend/Dockerfile"] = `FROM node:20-alpine AS build

WORKDIR /app

COPY package.json ./

RUN npm install

COPY . .

RUN npm run build

FROM nginx:alpine

COPY --from=build /app/build /usr/share/nginx/html

EXPOSE 80`
	}

	if input.BackendTechnology == "golang" {
		files["backend/Dockerfile"] = `FROM golang:1.23-alpine AS build

WORKDIR /app

COPY . .

RUN go build -o server ./cmd/api

FROM alpine:latest

WORKDIR /app

COPY --from=build /app/server /app/server

EXPOSE 8080

CMD ["/app/server"]`
	}

	if input.BackendTechnology == "java" {
		files["backend/Dockerfile"] = `FROM maven:3.9-eclipse-temurin-21 AS build

WORKDIR /app

COPY . .

RUN mvn clean package -DskipTests

FROM eclipse-temurin:21-jre

WORKDIR /app

COPY --from=build /app/target/*.jar app.jar

EXPOSE 8080

CMD ["java", "-jar", "app.jar"]`
	}

	if input.BackendTechnology == "node" {
		files["backend/Dockerfile"] = `FROM node:20-alpine

WORKDIR /app

COPY package.json ./
RUN npm install

COPY . .

EXPOSE 8080

CMD ["npm", "start"]`
	}

	if input.BackendTechnology == "python" {
		files["backend/Dockerfile"] = `FROM python:3.12-slim

WORKDIR /app

COPY requirements.txt ./
RUN pip install -r requirements.txt

COPY . .

EXPOSE 8080

CMD ["uvicorn", "main:app", "--host", "0.0.0.0", "--port", "8080"]`
	}

	if input.BackendTechnology == "dotnet" {
		files["backend/Dockerfile"] = `FROM mcr.microsoft.com/dotnet/sdk:8.0 AS build
WORKDIR /src
COPY . .
RUN dotnet publish -c Release -o /app/publish

FROM mcr.microsoft.com/dotnet/aspnet:8.0
WORKDIR /app
COPY --from=build /app/publish .
EXPOSE 8080
ENTRYPOINT ["dotnet", "generated-app.dll"]`
	}
}

func addGeneratedKubernetesFiles(
	files map[string]string,
	input generateProjectRequest,
) {
	projectName := sanitizeGeneratedProjectName(
		input.ProjectName,
	)

	files["k8s/backend-deployment.yaml"] = fmt.Sprintf(
		`apiVersion: apps/v1
kind: Deployment
metadata:
  name: %s-backend
spec:
  replicas: 1
  selector:
    matchLabels:
      app: %s-backend
  template:
    metadata:
      labels:
        app: %s-backend
    spec:
      containers:
        - name: backend
          image: %s-backend:latest
          ports:
            - containerPort: 8080
`,
		projectName,
		projectName,
		projectName,
		projectName,
	)
}

func validFrontendTechnology(
	technology string,
) bool {
	switch normalizeGeneratedTechnology(technology) {
	case "react", "angular", "vue", "nextjs":
		return true

	default:
		return false
	}
}

func validBackendTechnology(
	technology string,
) bool {
	switch normalizeGeneratedTechnology(technology) {
	case "golang", "java", "node", "python", "dotnet":
		return true

	default:
		return false
	}
}

func validDatabaseTechnology(
	technology string,
) bool {
	switch normalizeGeneratedTechnology(technology) {
	case "postgresql", "mysql", "mongodb", "sqlite", "none":
		return true

	default:
		return false
	}
}

func normalizeGeneratedTechnology(value string) string {
	normalized := strings.ToLower(strings.TrimSpace(value))
	if normalized == "" {
		return ""
	}

	normalized = strings.ReplaceAll(normalized, "_", " ")
	normalized = strings.ReplaceAll(normalized, "-", " ")
	normalized = strings.ReplaceAll(normalized, ".", " ")
	normalized = strings.Join(strings.Fields(normalized), " ")

	switch {
	case strings.Contains(normalized, "react"):
		return "react"
	case strings.Contains(normalized, "angular"):
		return "angular"
	case strings.Contains(normalized, "vue"):
		return "vue"
	case strings.Contains(normalized, "next"):
		return "nextjs"
	case strings.Contains(normalized, "mongodb") || strings.Contains(normalized, "mongo db") || strings.Contains(normalized, "mongo"):
		return "mongodb"
	case strings.Contains(normalized, "mysql"):
		return "mysql"
	case strings.Contains(normalized, "sqlite"):
		return "sqlite"
	case strings.Contains(normalized, "postgres") || strings.Contains(normalized, "sql"):
		return "postgresql"
	case strings.Contains(normalized, "dotnet") || strings.Contains(normalized, "asp") || strings.Contains(normalized, "aspnet") || strings.Contains(normalized, "c sharp"):
		return "dotnet"
	case strings.Contains(normalized, "python") || strings.Contains(normalized, "fastapi") || strings.Contains(normalized, "django"):
		return "python"
	case strings.Contains(normalized, "node") || strings.Contains(normalized, "express"):
		return "node"
	case strings.Contains(normalized, "spring") || strings.Contains(normalized, "java"):
		return "java"
	case strings.Contains(normalized, "golang") || strings.Contains(normalized, "go lang") || strings.Contains(normalized, "go language") || (strings.Contains(normalized, "go") && !strings.Contains(normalized, "mongo") && !strings.Contains(normalized, "sql") && !strings.Contains(normalized, "postgres")):
		return "golang"
	case strings.Contains(normalized, "none"):
		return "none"
	default:
		return normalized
	}
}

func collectGeneratedProjectFilesFromDisk(rootPath string) ([]string, error) {
	if rootPath == "" {
		return []string{}, nil
	}

	files := []string{}

	if err := filepath.Walk(rootPath, func(path string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if info.IsDir() {
			return nil
		}

		relPath, err := filepath.Rel(rootPath, path)
		if err != nil {
			return err
		}

		if relPath == "." || relPath == "" {
			return nil
		}

		files = append(files, filepath.ToSlash(relPath))
		return nil
	}); err != nil {
		return nil, err
	}

	return files, nil
}

func sanitizeGeneratedProjectName(
	name string,
) string {
	value := strings.ToLower(
		strings.TrimSpace(name),
	)

	value = strings.ReplaceAll(
		value,
		"_",
		"-",
	)

	value = strings.ReplaceAll(
		value,
		" ",
		"-",
	)

	expression := regexp.MustCompile(
		`[^a-z0-9-]+`,
	)

	value = expression.ReplaceAllString(
		value,
		"",
	)

	multipleHyphens := regexp.MustCompile(
		`-+`,
	)

	value = multipleHyphens.ReplaceAllString(
		value,
		"-",
	)

	value = strings.Trim(
		value,
		"-",
	)

	if value == "" {
		return "generated-project"
	}

	return value
}

func safeGeneratedFilePath(
	rootPath string,
	relativePath string,
) (string, error) {
	rootAbsolute, err := filepath.Abs(rootPath)

	if err != nil {
		return "", err
	}

	targetAbsolute, err := filepath.Abs(
		filepath.Join(
			rootAbsolute,
			filepath.Clean(relativePath),
		),
	)

	if err != nil {
		return "", err
	}

	expectedPrefix :=
		rootAbsolute + string(os.PathSeparator)

	if targetAbsolute != rootAbsolute &&
		!strings.HasPrefix(
			targetAbsolute,
			expectedPrefix,
		) {
		return "", fmt.Errorf(
			"invalid generated file path",
		)
	}

	return targetAbsolute, nil
}

func escapeGeneratedJS(value string) string {
	value = strings.ReplaceAll(
		value,
		`\`,
		`\\`,
	)

	value = strings.ReplaceAll(
		value,
		`"`,
		`\"`,
	)

	value = strings.ReplaceAll(
		value,
		"\r",
		" ",
	)

	value = strings.ReplaceAll(
		value,
		"\n",
		" ",
	)

	return value
}

func runGeneratedGit(
	directory string,
	arguments ...string,
) error {
	command := exec.Command(
		"git",
		arguments...,
	)

	command.Dir = directory

	output, err := command.CombinedOutput()

	if err != nil {
		message := strings.TrimSpace(
			string(output),
		)

		if message == "" {
			message = err.Error()
		}

		return fmt.Errorf("%s", message)
	}

	return nil
}
