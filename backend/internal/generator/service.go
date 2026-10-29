// package generator

// import (
//     "fmt"
//     "os"
//     "path/filepath"
//     "strings"
// )

// type ProjectResult struct {
//     Files []string
// }

// func GenerateProjectSkeleton(root, appName, requirements string) (ProjectResult, error) {
//     if strings.TrimSpace(root) == "" {
//         return ProjectResult{}, fmt.Errorf("root is required")
//     }
//     if err := os.MkdirAll(root, 0o755); err != nil {
//         return ProjectResult{}, err
//     }

//     summary := strings.TrimSpace(requirements)
//     if summary == "" {
//         summary = "Create a modern storefront application with authentication, catalog, cart, checkout, orders, and admin reporting."
//     }

//     files := []string{}

//     if err := os.WriteFile(filepath.Join(root, "README.md"), []byte(buildReadme(appName, summary)), 0o644); err != nil {
//         return ProjectResult{}, err
//     }
//     files = append(files, "README.md")

//     if err := os.MkdirAll(filepath.Join(root, "backend", "internal", "api"), 0o755); err != nil {
//         return ProjectResult{}, err
//     }
//     if err := os.MkdirAll(filepath.Join(root, "frontend", "src"), 0o755); err != nil {
//         return ProjectResult{}, err
//     }

//     backendFiles := map[string]string{
//         "backend/go.mod":                        buildGoMod(appName),
//         "backend/main.go":                      buildMainGo(appName),
//         "backend/internal/api/routes.go":       buildRoutesGo(),
//         "backend/internal/api/handlers.go":     buildHandlersGo(),
//         "backend/internal/store/store.go":      buildStoreGo(),
//         "backend/.env.example":                 buildEnvExample(),
//     }
//     for relativePath, content := range backendFiles {
//         fullPath := filepath.Join(root, relativePath)
//         if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
//             return ProjectResult{}, err
//         }
//         if err := os.WriteFile(fullPath, []byte(content), 0o644); err != nil {
//             return ProjectResult{}, err
//         }
//         files = append(files, relativePath)
//     }

//     frontendFiles := map[string]string{
//         "frontend/package.json":                 buildFrontendPackage(appName),
//         "frontend/tsconfig.json":                buildTsConfig(),
//         "frontend/src/App.tsx":                  buildFrontendApp(),
//         "frontend/src/main.tsx":                 buildFrontendMain(),
//         "frontend/src/index.css":                buildFrontendCSS(),
//         "frontend/.env.example":                 "REACT_APP_API_URL=http://localhost:8080\n",
//         "frontend/Dockerfile":                   "FROM node:20-alpine\nWORKDIR /app\nCOPY package*.json ./\nRUN npm install\nCOPY . .\nRUN npm run build\nEXPOSE 3000\nCMD [\"npm\", \"start\"]\n",
//         "frontend/nginx.conf":                  "server {\n  listen 80;\n  location / {\n    root /usr/share/nginx/html;\n    try_files $uri /index.html;\n  }\n}\n",
//         "frontend/.github/workflows/ci.yml":    "name: frontend-ci\non:\n  push:\n    branches: ['**']\n  pull_request:\n    branches: ['**']\njobs:\n  build:\n    runs-on: ubuntu-latest\n    steps:\n      - uses: actions/checkout@v4\n      - uses: actions/setup-node@v4\n        with:\n          node-version: 20\n      - run: npm install\n      - run: npm run build\n",
//     }
//     for relativePath, content := range frontendFiles {
//         fullPath := filepath.Join(root, relativePath)
//         if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
//             return ProjectResult{}, err
//         }
//         if err := os.WriteFile(fullPath, []byte(content), 0o644); err != nil {
//             return ProjectResult{}, err
//         }
//         files = append(files, relativePath)
//     }

//     if err := os.MkdirAll(filepath.Join(root, "kubernetes"), 0o755); err != nil {
//         return ProjectResult{}, err
//     }
//     if err := os.WriteFile(filepath.Join(root, "docker-compose.yml"), []byte(buildDockerCompose()), 0o644); err != nil {
//         return ProjectResult{}, err
//     }
//     files = append(files, "docker-compose.yml")

//     if err := os.WriteFile(filepath.Join(root, "kubernetes", "deployment.yaml"), []byte(buildKubernetesDeployment()), 0o644); err != nil {
//         return ProjectResult{}, err
//     }
//     files = append(files, "kubernetes/deployment.yaml")

//     return ProjectResult{Files: files}, nil
// }

// func buildReadme(appName, requirements string) string {
//     return fmt.Sprintf("# %s\n\nThis project was generated from the requirement:\n\n%s\n\n## Features\n- User authentication\n- Product catalog\n- Shopping cart\n- Secure checkout flow\n- Order management\n- Admin product management\n- Reports dashboard\n\n## Run locally\n\n```bash\ndocker-compose up --build\n```\n", appName, requirements)
// }

// func buildGoMod(appName string) string {
//     return "module " + sanitizeModuleName(appName) + "\n\ngo 1.23\n"
// }

// func buildMainGo(appName string) string {
//     return fmt.Sprintf("package main\n\nimport (\n    \"fmt\"\n    \"net/http\"\n)\n\nfunc main() {\n    http.HandleFunc(\"/health\", func(w http.ResponseWriter, r *http.Request) {\n        w.WriteHeader(http.StatusOK)\n        _, _ = w.Write([]byte(\"ok\"))\n    })\n\n    http.HandleFunc(\"/api/products\", func(w http.ResponseWriter, r *http.Request) {\n        w.Header().Set(\"Content-Type\", \"application/json\")\n        _, _ = w.Write([]byte(`[{\"id\":\"p1\",\"name\":\"Cappuccino\",\"price\":4.5}]`))\n    })\n\n    fmt.Println(\"%s backend running on :8080\")\n    _ = http.ListenAndServe(\":8080\", nil)\n}\n", appName)
// }

// func buildRoutesGo() string {
//     return "package api\n\nimport (\n    \"net/http\"\n)\n\nfunc RegisterRoutes() http.Handler {\n    mux := http.NewServeMux()\n    mux.HandleFunc(\"/health\", func(w http.ResponseWriter, r *http.Request) {\n        w.WriteHeader(http.StatusOK)\n        _, _ = w.Write([]byte(\"ok\"))\n    })\n    return mux\n}\n"
// }

// func buildHandlersGo() string {
//     return "package api\n\nimport \"net/http\"\n\nfunc HelloHandler(w http.ResponseWriter, r *http.Request) {\n    w.Header().Set(\"Content-Type\", \"application/json\")\n    _, _ = w.Write([]byte(`{\"message\":\"hello from generated backend\"}`))\n}\n"
// }

// func buildStoreGo() string {
//     return "package store\n\nvar Products = []map[string]any{\n    {\"id\": \"p1\", \"name\": \"Cappuccino\", \"price\": 4.5},\n    {\"id\": \"p2\", \"name\": \"Latte\", \"price\": 5.5},\n}\n"
// }

// func buildEnvExample() string {
//     return "PORT=8080\nJWT_SECRET=change-me\nDATABASE_URL=postgres://postgres:postgres@db:5432/app\n"
// }

// func buildFrontendPackage(appName string) string {
//     return fmt.Sprintf("{\n  \"name\": \"%s\",\n  \"version\": \"1.0.0\",\n  \"private\": true,\n  \"dependencies\": {\n    \"react\": \"^18.2.0\",\n    \"react-dom\": \"^18.2.0\",\n    \"axios\": \"^1.7.2\"\n  },\n  \"devDependencies\": {\n    \"@types/react\": \"^18.2.66\",\n    \"@types/react-dom\": \"^18.2.0\",\n    \"@vitejs/plugin-react\": \"^4.3.1\",\n    \"typescript\": \"^5.5.4\",\n    \"vite\": \"^5.4.0\"\n  },\n  \"scripts\": {\n    \"dev\": \"vite\",\n    \"build\": \"tsc && vite build\",\n    \"start\": \"vite --host 0.0.0.0 --port 3000\"\n  }\n}\n", sanitizeName(appName))
// }

// func buildTsConfig() string {
//     return "{\n  \"compilerOptions\": {\n    \"target\": \"ES2020\",\n    \"useDefineForClassFields\": true,\n    \"lib\": [\"DOM\", \"DOM.Iterable\", \"ES2020\"],\n    \"allowJs\": false,\n    \"skipLibCheck\": true,\n    \"esModuleInterop\": true,\n    \"allowSyntheticDefaultImports\": true,\n    \"strict\": true,\n    \"forceConsistentCasingInFileNames\": true,\n    \"module\": \"ESNext\",\n    \"moduleResolution\": \"Node\",\n    \"resolveJsonModule\": true,\n    \"isolatedModules\": true,\n    \"noEmit\": true,\n    \"jsx\": \"react-jsx\"\n  },\n  \"include\": [\"src\"],\n  \"references\": []\n}\n"
// }

// func buildFrontendApp() string {
//     return "import React from 'react';\n\nexport default function App() {\n  return (\n    <main style={{ fontFamily: 'sans-serif', padding: '2rem' }}>\n      <h1>Coffee Shop</h1>\n      <p>Welcome to the generated coffee shop app.</p>\n      <ul>\n        <li>User login</li>\n        <li>Products</li>\n        <li>Cart and checkout</li>\n        <li>Orders and admin dashboard</li>\n      </ul>\n    </main>\n  );\n}\n"
// }

// func buildFrontendMain() string {
//     return "import React from 'react';\nimport ReactDOM from 'react-dom/client';\nimport App from './App';\nimport './index.css';\n\nReactDOM.createRoot(document.getElementById('root')!).render(\n  <React.StrictMode>\n    <App />\n  </React.StrictMode>,\n);\n"
// }

// func buildFrontendCSS() string {
//     return "body { margin: 0; font-family: Arial, sans-serif; background: #f7f1ea; color: #2a1d12; }\nmain { max-width: 960px; margin: 0 auto; padding: 2rem; }\nul { line-height: 2; }\n"
// }

// func buildDockerCompose() string {
//     return "services:\n  backend:\n    build: ./backend\n    ports:\n      - \"8080:8080\"\n    environment:\n      - PORT=8080\n  frontend:\n    build: ./frontend\n    ports:\n      - \"3000:3000\"\n    depends_on:\n      - backend\n    environment:\n      - REACT_APP_API_URL=http://localhost:8080\n"
// }

// func buildKubernetesDeployment() string {
//     return "apiVersion: apps/v1\nkind: Deployment\nmetadata:\n  name: coffee-shop\nspec:\n  replicas: 2\n  selector:\n    matchLabels:\n      app: coffee-shop\n  template:\n    metadata:\n      labels:\n        app: coffee-shop\n    spec:\n      containers:\n        - name: coffee-shop\n          image: coffee-shop:latest\n          ports:\n            - containerPort: 3000\n---\napiVersion: v1\nkind: Service\nmetadata:\n  name: coffee-shop\nspec:\n  selector:\n    app: coffee-shop\n  ports:\n    - protocol: TCP\n      port: 80\n      targetPort: 3000\n"
// }

// func sanitizeName(input string) string {
//     value := strings.ToLower(strings.TrimSpace(input))
//     value = strings.ReplaceAll(value, " ", "-")
//     value = strings.Map(func(r rune) rune {
//         if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' || r == '_' {
//             return r
//         }
//         return '-'
//     }, value)
//     return strings.Trim(value, "-")
// }

// func sanitizeModuleName(input string) string {
//     return strings.ReplaceAll(strings.TrimSpace(input), " ", "-")
// }

package generator

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

type FileDefinition struct {
	Path    string
	Content string
}

type Service struct {
	BaseDirectory string
}

func NewService(baseDirectory string) *Service {
	return &Service{
		BaseDirectory: baseDirectory,
	}
}

func (service *Service) Generate(
	request GenerateProjectRequest,
	projectID string,
) (string, []FileDefinition, error) {
	projectName := sanitizeProjectName(request.ProjectName)

	if projectName == "" {
		return "", nil, errors.New("invalid project name")
	}

	if err := validateTechnologyStack(request); err != nil {
		return "", nil, err
	}

	rootPath := filepath.Join(
		service.BaseDirectory,
		projectName+"-"+projectID,
	)

	if err := os.MkdirAll(rootPath, 0755); err != nil {
		return "", nil, fmt.Errorf("create project directory: %w", err)
	}

	files := make([]FileDefinition, 0)

	frontendFiles, err := buildFrontendFiles(request, projectName)
	if err != nil {
		return "", nil, err
	}

	backendFiles, err := buildBackendFiles(request, projectName)
	if err != nil {
		return "", nil, err
	}

	files = append(files, frontendFiles...)
	files = append(files, backendFiles...)
	files = append(files, buildRootFiles(request, projectName)...)

	for _, file := range files {
		targetPath, err := safeJoin(rootPath, file.Path)
		if err != nil {
			return "", nil, err
		}

		if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
			return "", nil, fmt.Errorf(
				"create directory for %s: %w",
				file.Path,
				err,
			)
		}

		if err := os.WriteFile(
			targetPath,
			[]byte(file.Content),
			0644,
		); err != nil {
			return "", nil, fmt.Errorf(
				"write generated file %s: %w",
				file.Path,
				err,
			)
		}
	}

	sort.Slice(files, func(left int, right int) bool {
		return files[left].Path < files[right].Path
	})

	return rootPath, files, nil
}

func validateTechnologyStack(request GenerateProjectRequest) error {
	validFrontend := map[string]bool{
		"react":   true,
		"angular": true,
	}

	validBackend := map[string]bool{
		"golang": true,
		"java":   true,
	}

	validDatabase := map[string]bool{
		"postgresql": true,
		"mysql":      true,
		"mongodb":    true,
		"none":       true,
	}

	if !validFrontend[strings.ToLower(request.FrontendTechnology)] {
		return errors.New("unsupported frontend technology")
	}

	if !validBackend[strings.ToLower(request.BackendTechnology)] {
		return errors.New("unsupported backend technology")
	}

	if !validDatabase[strings.ToLower(request.DatabaseTechnology)] {
		return errors.New("unsupported database technology")
	}

	return nil
}

func sanitizeProjectName(value string) string {
	value = strings.TrimSpace(strings.ToLower(value))
	value = strings.ReplaceAll(value, "_", "-")
	value = strings.ReplaceAll(value, " ", "-")

	expression := regexp.MustCompile(`[^a-z0-9-]+`)
	value = expression.ReplaceAllString(value, "")

	multipleHyphens := regexp.MustCompile(`-+`)
	value = multipleHyphens.ReplaceAllString(value, "-")

	return strings.Trim(value, "-")
}

func safeJoin(rootPath string, relativePath string) (string, error) {
	cleanRoot, err := filepath.Abs(rootPath)
	if err != nil {
		return "", err
	}

	targetPath, err := filepath.Abs(
		filepath.Join(cleanRoot, filepath.Clean(relativePath)),
	)
	if err != nil {
		return "", err
	}

	expectedPrefix := cleanRoot + string(os.PathSeparator)

	if targetPath != cleanRoot &&
		!strings.HasPrefix(targetPath, expectedPrefix) {
		return "", errors.New("invalid generated file path")
	}

	return targetPath, nil
}
