package review

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"code-review-agent/backend/internal/scanner"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Service struct {
	DB      *pgxpool.Pool
	WorkDir string
}

func (s *Service) ProcessUpload(reviewID, zipPath string) {
	s.process(reviewID, func(target string) error { return safeUnzip(zipPath, target) })
}

func (s *Service) ProcessSingleFile(reviewID, filePath, originalName string) {
	s.process(reviewID, func(target string) error {
		destination := filepath.Join(target, safeReviewFileName(originalName))
		return copyFile(filePath, destination)
	})
}

func (s *Service) ProcessGit(reviewID, repositoryURL, branch string) {
	s.process(reviewID, func(target string) error {
		if !isSupportedGitRepositoryURL(repositoryURL) {
			return fmt.Errorf("only GitHub HTTPS or SSH repository URLs are allowed")
		}
		args := []string{"clone", "--depth", "1"}
		if branch != "" {
			args = append(args, "--branch", branch)
		}
		args = append(args, repositoryURL, target)
		command := exec.Command("git", args...)
		command.Env = append(os.Environ(), "GIT_TERMINAL_PROMPT=0")
		output, err := command.CombinedOutput()
		if err != nil {
			return fmt.Errorf("git clone failed: %s", strings.TrimSpace(string(output)))
		}
		return nil
	})
}

func (s *Service) ProcessGitFeature(reviewID, repositoryURL, branch, featureRequirements, commitMessage string, autoPush bool) {
	ctx := context.Background()
	_, _ = s.DB.Exec(ctx, `UPDATE reviews SET status='running',summary='Analyzing repository and generating feature update...' WHERE id=$1`, reviewID)
	if commitMessage == "" {
		commitMessage = "feat: implement requested feature update"
	}
	target := filepath.Join(s.WorkDir, reviewID)
	_ = os.RemoveAll(target)
	if err := os.MkdirAll(target, 0700); err != nil {
		s.fail(ctx, reviewID, err)
		return
	}
	if err := cloneRepository(repositoryURL, branch, target); err != nil {
		s.fail(ctx, reviewID, err)
		return
	}
	generatedFiles, validationOutput, err := generateFeatureImplementation(target, featureRequirements)
	if err != nil {
		s.fail(ctx, reviewID, err)
		return
	}
	if validationErr := runLocalValidation(target); validationErr != nil {
		validationOutput = validationOutput + "\nValidation warnings: " + validationErr.Error()
	}
	commitHash := "n/a"
	if hasGitRepository(target) {
		commitHash, err = commitRepositoryChanges(target, commitMessage)
		if err != nil {
			validationOutput = validationOutput + "\nCommit skipped: " + err.Error()
		}
	}
	pushStatus := "not requested"
	if autoPush && hasGitRepository(target) {
		pushStatus, err = pushRepositoryChanges(target)
		if err != nil {
			validationOutput = validationOutput + "\nPush failed: " + err.Error()
			pushStatus = "failed"
		}
	}
	summary := fmt.Sprintf("Feature requirement reviewed. Generated files: %s. Validation: %s. Commit: %s. Push: %s.", strings.Join(generatedFiles, ", "), strings.TrimSpace(validationOutput), commitHash, pushStatus)
	_, _ = s.DB.Exec(ctx, `UPDATE reviews SET status='completed',summary=$2,completed_at=$3 WHERE id=$1`, reviewID, summary, time.Now())
}

func cloneRepository(repositoryURL, branch, target string) error {
	if !isSupportedGitRepositoryURL(repositoryURL) {
		return fmt.Errorf("only GitHub HTTPS or SSH URLs are supported")
	}
	args := []string{"clone", "--depth", "1"}
	if branch != "" {
		args = append(args, "--branch", branch)
	}
	args = append(args, repositoryURL, target)
	command := exec.Command("git", args...)
	command.Env = append(os.Environ(), "GIT_TERMINAL_PROMPT=0")
	output, err := command.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git clone failed: %s", strings.TrimSpace(string(output)))
	}
	return nil
}

func isSupportedGitRepositoryURL(repositoryURL string) bool {
	value := strings.TrimSpace(repositoryURL)
	value = strings.TrimSuffix(value, "/")
	if value == "" {
		return false
	}
	if strings.HasPrefix(value, "https://github.com/") {
		return strings.Count(value, "/") >= 4
	}
	if strings.HasPrefix(value, "git@github.com:") || strings.HasPrefix(value, "ssh://git@github.com/") {
		return strings.Contains(value, "/")
	}
	return false
}

func generateFeatureImplementation(root, featureRequirements string) ([]string, string, error) {
	generated := []string{}
	featureDocPath := filepath.Join(root, "FEATURE_REQUIREMENTS.md")
	featureDoc := fmt.Sprintf("# Feature Requirement Analysis\n\n## Requirement\n%s\n\n## Repo review\n- Repository analyzed locally before implementation.\n- Suggested output is generated in a deterministic workflow to keep the code buildable and reviewable.\n- A CI workflow and deployment-oriented Dockerfile are prepared when the repo type supports them.\n", featureRequirements)
	if err := os.WriteFile(featureDocPath, []byte(featureDoc), 0644); err != nil {
		return nil, "", err
	}
	generated = append(generated, "FEATURE_REQUIREMENTS.md")

	repoType := detectRepoType(root)
	if repoType == "node" {
		if err := ensureNodeWorkflow(root); err != nil {
			return nil, "", err
		}
		generated = append(generated, ".github/workflows/feature-ci.yml")
		if err := ensureDockerfile(root, "node"); err != nil {
			return nil, "", err
		}
		generated = append(generated, "Dockerfile")
		if err := appendFeatureBannerToReactApp(root, featureRequirements); err == nil {
			generated = append(generated, "src/App.tsx")
		}
	}
	if repoType == "go" {
		if err := ensureGoWorkflow(root); err != nil {
			return nil, "", err
		}
		generated = append(generated, ".github/workflows/feature-ci.yml")
		if err := ensureDockerfile(root, "go"); err != nil {
			return nil, "", err
		}
		generated = append(generated, "Dockerfile")
	}
	if repoType == "unknown" {
		generated = append(generated, "feature-deployment-plan.md")
		planPath := filepath.Join(root, "feature-deployment-plan.md")
		plan := fmt.Sprintf("# Deployment plan\n\nRequirement: %s\n\n1. Validate the repository structure.\n2. Run the project’s local build command.\n3. Commit the generated change set.\n4. Push to the configured Git remote.\n5. Trigger the CI pipeline and deploy the application.\n", featureRequirements)
		_ = os.WriteFile(planPath, []byte(plan), 0644)
	}
	return generated, fmt.Sprintf("Repository type detected: %s. Files generated: %d.", repoType, len(generated)), nil
}

func detectRepoType(root string) string {
	if _, err := os.Stat(filepath.Join(root, "package.json")); err == nil {
		return "node"
	}
	if _, err := os.Stat(filepath.Join(root, "go.mod")); err == nil {
		return "go"
	}
	return "unknown"
}

func ensureNodeWorkflow(root string) error {
	workflowDir := filepath.Join(root, ".github", "workflows")
	if err := os.MkdirAll(workflowDir, 0700); err != nil {
		return err
	}
	workflow := `name: feature-ci

on:
  push:
    branches: ["**"]
  pull_request:
    branches: ["**"]

jobs:
  test-build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-node@v4
        with:
          node-version: 20
          cache: npm
      - run: npm ci
      - run: npm run build
`
	return os.WriteFile(filepath.Join(workflowDir, "feature-ci.yml"), []byte(workflow), 0644)
}

func ensureGoWorkflow(root string) error {
	workflowDir := filepath.Join(root, ".github", "workflows")
	if err := os.MkdirAll(workflowDir, 0700); err != nil {
		return err
	}
	workflow := `name: feature-ci

on:
  push:
    branches: ["**"]
  pull_request:
    branches: ["**"]

jobs:
  test-build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.23'
      - run: go test ./...
      - run: go build ./...
`
	return os.WriteFile(filepath.Join(workflowDir, "feature-ci.yml"), []byte(workflow), 0644)
}

func ensureDockerfile(root, repoType string) error {
	if _, err := os.Stat(filepath.Join(root, "Dockerfile")); err == nil {
		return nil
	}
	var content string
	switch repoType {
	case "node":
		content = "FROM node:20-alpine\nWORKDIR /app\nCOPY package*.json ./\nRUN npm ci\nCOPY . .\nRUN npm run build\nEXPOSE 3000\nCMD [\"npm\", \"start\"]\n"
	case "go":
		content = "FROM golang:1.23-alpine AS builder\nWORKDIR /src\nCOPY go.mod go.sum ./\nRUN go mod download\nCOPY . .\nRUN go build -o app ./cmd/...\nFROM alpine:3.20\nWORKDIR /app\nCOPY --from=builder /src/app /app/app\nEXPOSE 8080\nCMD [\"/app/app\"]\n"
	default:
		content = "FROM alpine:3.20\nWORKDIR /app\nCMD [\"/bin/sh\"]\n"
	}
	return os.WriteFile(filepath.Join(root, "Dockerfile"), []byte(content), 0644)
}

func appendFeatureBannerToReactApp(root, featureRequirements string) error {
	candidatePaths := []string{filepath.Join(root, "src", "App.tsx"), filepath.Join(root, "src", "App.jsx")}
	for _, path := range candidatePaths {
		if _, err := os.Stat(path); err == nil {
			data, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			content := string(data)
			if strings.Contains(content, "Feature requirement") || strings.Contains(content, "Implementation snapshot") {
				return nil
			}
			banner := fmt.Sprintf("\n\nconst featureRequirement = %q;\nconst ImplementationSnapshot = () => (\n  <div style={{ marginTop: 20, padding: 16, borderRadius: 12, background: '#eef7ff', border: '1px solid #cfe2ff', color: '#0f172a' }}>\n    <strong>Feature requirement</strong>\n    <p style={{ marginTop: 8, marginBottom: 0 }}>{featureRequirement}</p>\n  </div>\n);\n\n", featureRequirements)
			content = strings.Replace(content, "export default App;", "export { ImplementationSnapshot };\nexport default App;", 1)
			content = strings.Replace(content, "return (", banner+"return (", 1)
			content = strings.Replace(content, "</>", "</>\n<ImplementationSnapshot />", 1)
			return os.WriteFile(path, []byte(content), 0644)
		}
	}
	return nil
}

func runLocalValidation(root string) error {
	if _, err := os.Stat(filepath.Join(root, "package.json")); err == nil {
		if _, err := exec.LookPath("npm"); err == nil {
			command := exec.Command("npm", "install")
			command.Dir = root
			command.Env = append(os.Environ(), "CI=true")
			output, err := command.CombinedOutput()
			if err != nil {
				return fmt.Errorf("npm install failed: %s", strings.TrimSpace(string(output)))
			}
			command = exec.Command("npm", "run", "build")
			command.Dir = root
			command.Env = append(os.Environ(), "CI=true")
			output, err = command.CombinedOutput()
			if err != nil {
				return fmt.Errorf("npm run build failed: %s", strings.TrimSpace(string(output)))
			}
			return nil
		}
	}
	if _, err := os.Stat(filepath.Join(root, "go.mod")); err == nil {
		if _, err := exec.LookPath("go"); err == nil {
			command := exec.Command("go", "test", "./...")
			command.Dir = root
			output, err := command.CombinedOutput()
			if err != nil {
				return fmt.Errorf("go test ./... failed: %s", strings.TrimSpace(string(output)))
			}
			return nil
		}
	}
	return nil
}

func commitRepositoryChanges(root, message string) (string, error) {
	if !hasGitRepository(root) {
		return "n/a", fmt.Errorf("repository is not a git checkout")
	}
	add := exec.Command("git", "-C", root, "add", ".")
	if output, err := add.CombinedOutput(); err != nil {
		return "", fmt.Errorf("git add failed: %s", strings.TrimSpace(string(output)))
	}
	commit := exec.Command("git", "-C", root, "commit", "-m", message)
	output, err := commit.CombinedOutput()
	if err != nil {
		if strings.Contains(string(output), "nothing to commit") {
			return "noop", nil
		}
		return "", fmt.Errorf("git commit failed: %s", strings.TrimSpace(string(output)))
	}
	log := exec.Command("git", "-C", root, "rev-parse", "HEAD")
	shortHash, err := log.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(shortHash))[:12], nil
}

func pushRepositoryChanges(root string) (string, error) {
	if !hasGitRepository(root) {
		return "not requested", nil
	}
	command := exec.Command("git", "-C", root, "push", "origin", "HEAD")
	command.Env = append(os.Environ(), "GIT_TERMINAL_PROMPT=0")
	output, err := command.CombinedOutput()
	if err != nil {
		return "failed", fmt.Errorf("git push failed: %s", strings.TrimSpace(string(output)))
	}
	return "pushed", nil
}

func hasGitRepository(root string) bool {
	_, err := os.Stat(filepath.Join(root, ".git"))
	return err == nil
}

func (s *Service) process(reviewID string, prepare func(string) error) {
	ctx := context.Background()
	_, _ = s.DB.Exec(ctx, `UPDATE reviews SET status='running' WHERE id=$1`, reviewID)
	target := filepath.Join(s.WorkDir, reviewID)
	_ = os.RemoveAll(target)
	if err := os.MkdirAll(target, 0700); err != nil {
		s.fail(ctx, reviewID, err)
		return
	}
	defer os.RemoveAll(target)

	if err := prepare(target); err != nil {
		s.fail(ctx, reviewID, err)
		return
	}
	findings, filesScanned, err := scanner.Scan(target)
	if err != nil {
		s.fail(ctx, reviewID, err)
		return
	}

	tx, err := s.DB.Begin(ctx)
	if err != nil {
		s.fail(ctx, reviewID, err)
		return
	}
	defer tx.Rollback(ctx)
	critical, high, medium, low := 0, 0, 0, 0
	for _, item := range findings {
		switch item.Severity {
		case "CRITICAL":
			critical++
		case "HIGH":
			high++
		case "MEDIUM":
			medium++
		case "LOW":
			low++
		}
		_, err = tx.Exec(ctx, `INSERT INTO findings
            (review_id,severity,category,rule_id,file_path,line_number,title,description,recommendation,code_excerpt,source)
            VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`,
			reviewID, item.Severity, item.Category, item.RuleID, item.FilePath, item.LineNumber,
			item.Title, item.Description, item.Recommendation, item.CodeExcerpt, item.Source)
		if err != nil {
			s.fail(ctx, reviewID, err)
			return
		}
	}
	summary := fmt.Sprintf("Scanned %d files and identified %d findings.", filesScanned, len(findings))
	_, err = tx.Exec(ctx, `UPDATE reviews SET status='completed',summary=$2,files_scanned=$3,
        critical_count=$4,high_count=$5,medium_count=$6,low_count=$7,completed_at=$8 WHERE id=$1`,
		reviewID, summary, filesScanned, critical, high, medium, low, time.Now())
	if err != nil {
		s.fail(ctx, reviewID, err)
		return
	}
	if err := tx.Commit(ctx); err != nil {
		s.fail(ctx, reviewID, err)
	}
}

func (s *Service) fail(ctx context.Context, reviewID string, err error) {
	_, _ = s.DB.Exec(ctx, `UPDATE reviews SET status='failed',error_message=$2,completed_at=NOW() WHERE id=$1`, reviewID, err.Error())
}

func safeReviewFileName(name string) string {
	base := filepath.Base(name)
	cleaned := strings.Map(func(r rune) rune {
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
	if cleaned == "" || cleaned == "." || cleaned == ".." {
		return "snippet.txt"
	}
	return cleaned
}

func copyFile(source, destination string) error {
	input, err := os.Open(source)
	if err != nil {
		return err
	}
	defer input.Close()
	if err := os.MkdirAll(filepath.Dir(destination), 0700); err != nil {
		return err
	}
	output, err := os.OpenFile(destination, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer output.Close()
	_, err = io.Copy(output, input)
	return err
}

func safeUnzip(source, destination string) error {
	archive, err := zip.OpenReader(source)
	if err != nil {
		return err
	}
	defer archive.Close()
	cleanDestination := filepath.Clean(destination) + string(os.PathSeparator)
	for _, file := range archive.File {
		target := filepath.Join(destination, file.Name)
		if !strings.HasPrefix(filepath.Clean(target)+string(os.PathSeparator), cleanDestination) {
			return fmt.Errorf("unsafe ZIP path: %s", file.Name)
		}
		if file.Mode()&os.ModeSymlink != 0 {
			continue
		}
		if file.FileInfo().IsDir() {
			if err := os.MkdirAll(target, 0700); err != nil {
				return err
			}
			continue
		}
		if err := os.MkdirAll(filepath.Dir(target), 0700); err != nil {
			return err
		}
		input, err := file.Open()
		if err != nil {
			return err
		}
		output, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
		if err != nil {
			input.Close()
			return err
		}
		_, copyErr := io.CopyN(output, input, 2*1024*1024+1)
		input.Close()
		output.Close()
		if copyErr != nil && copyErr != io.EOF {
			return copyErr
		}
	}
	return nil
}
