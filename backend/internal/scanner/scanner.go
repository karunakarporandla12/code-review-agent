package scanner

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type Finding struct {
	Severity       string
	Category       string
	RuleID         string
	FilePath       string
	LineNumber     int
	Title          string
	Description    string
	Recommendation string
	CodeExcerpt    string
	Source         string
}

type Rule struct {
	severity       string
	category       string
	id             string
	title          string
	pattern        *regexp.Regexp
	description    string
	recommendation string
}

var rules = []Rule{
	{"CRITICAL", "Security", "SECRET-001", "Possible private key", regexp.MustCompile(`BEGIN (RSA |EC |OPENSSH )?PRIVATE KEY`), "A private key marker appears in source control.", "Remove the key, rotate it immediately, and use a secret manager."},
	{"HIGH", "Security", "SECRET-002", "Possible hardcoded secret", regexp.MustCompile(`(?i)(api[_-]?key|secret|password|token)\s*[:=]\s*["'][^"']{8,}["']`), "A credential-like value appears to be hardcoded.", "Move the value to an approved secret manager or environment variable and rotate exposed credentials."},
	{"HIGH", "Security", "SQL-001", "Possible SQL injection", regexp.MustCompile(`(?i)(SELECT|INSERT|UPDATE|DELETE).*(\+|fmt\.Sprintf)`), "SQL appears to be constructed through string concatenation or formatting.", "Use a parameterized query and validate input."},
	{"HIGH", "Security", "CMD-001", "Possible command injection", regexp.MustCompile(`(?i)(exec\.Command|os\.system|child_process\.exec).*\+`), "A system command may contain concatenated input.", "Use fixed command arguments, allowlists, and avoid shell interpretation."},
	{"MEDIUM", "Quality", "ERROR-001", "Ignored error", regexp.MustCompile(`(?i)\b(err|error)\s*:=.*\n?`), "Review whether the returned error is handled on the following lines.", "Return, wrap, log, or intentionally document ignored errors."},
	{"MEDIUM", "Debugging", "DEBUG-001", "Debug statement", regexp.MustCompile(`\b(console\.log|fmt\.Println|debugger;)`), "Debug output is present in application code.", "Remove debug output or replace it with structured logging at the correct level."},
	{"LOW", "Maintainability", "TODO-001", "Unresolved TODO", regexp.MustCompile(`(?i)\b(TODO|FIXME)\b`), "An unresolved work marker remains in the code.", "Create a tracked work item or complete the change before release."},
}

var allowedExt = map[string]bool{
	".go": true, ".js": true, ".jsx": true, ".ts": true, ".tsx": true,
	".py": true, ".java": true, ".cs": true, ".sql": true, ".yml": true,
	".yaml": true, ".json": true, ".xml": true, ".tf": true, ".sh": true, ".txt": true,
}

func Scan(root string) ([]Finding, int, error) {
	findings := []Finding{}
	filesScanned := 0
	err := filepath.Walk(root, func(path string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if info.IsDir() {
			name := info.Name()
			if name == ".git" || name == "node_modules" || name == "vendor" || name == "build" || name == "dist" {
				return filepath.SkipDir
			}
			return nil
		}
		if info.Size() > 1024*1024 || !allowedExt[strings.ToLower(filepath.Ext(path))] {
			return nil
		}
		filesScanned++
		relative, _ := filepath.Rel(root, path)
		file, err := os.Open(path)
		if err != nil {
			return nil
		}
		defer file.Close()
		lineNo := 0
		input := bufio.NewScanner(file)
		input.Buffer(make([]byte, 64*1024), 1024*1024)
		for input.Scan() {
			lineNo++
			line := input.Text()
			for _, rule := range rules {
				if rule.pattern.MatchString(line) {
					findings = append(findings, Finding{
						Severity: rule.severity, Category: rule.category, RuleID: rule.id,
						FilePath: filepath.ToSlash(relative), LineNumber: lineNo, Title: rule.title,
						Description: rule.description, Recommendation: rule.recommendation,
						CodeExcerpt: strings.TrimSpace(line), Source: "built-in",
					})
				}
			}
		}
		return nil
	})
	return findings, filesScanned, err
}
