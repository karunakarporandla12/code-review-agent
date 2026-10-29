package generator

import (
	"fmt"
	"strings"
)

func buildFrontendFiles(
	request GenerateProjectRequest,
	projectName string,
) ([]FileDefinition, error) {
	switch strings.ToLower(request.FrontendTechnology) {
	case "react":
		return buildReactFiles(request, projectName), nil
	case "angular":
		return buildAngularFiles(request, projectName), nil
	default:
		return nil, fmt.Errorf(
			"unsupported frontend technology: %s",
			request.FrontendTechnology,
		)
	}
}

func buildReactFiles(
	request GenerateProjectRequest,
	projectName string,
) []FileDefinition {
	packageJSON := fmt.Sprintf(`{
  "name": "%s-frontend",
  "private": true,
  "version": "1.0.0",
  "scripts": {
    "start": "react-scripts start",
    "build": "react-scripts build",
    "test": "react-scripts test"
  },
  "dependencies": {
    "@testing-library/jest-dom": "^6.6.3",
    "@testing-library/react": "^16.1.0",
    "@testing-library/user-event": "^14.5.2",
    "axios": "^1.7.9",
    "bootstrap": "^5.3.3",
    "react": "^18.3.1",
    "react-dom": "^18.3.1",
    "react-router-dom": "^7.1.1",
    "react-scripts": "5.0.1"
  }
}
`, projectName)

	appCode := fmt.Sprintf(`import React from "react";
import "./App.css";

function App() {
  return (
    <div className="app-shell">
      <header className="app-header">
        <h1>%s</h1>
        <p>%s</p>
      </header>

      <main className="app-content">
        <section className="app-card">
          <h2>Application Generated Successfully</h2>
          <p>
            Replace this initial page with modules based on the application
            requirements.
          </p>
        </section>
      </main>
    </div>
  );
}

export default App;
`, projectName, escapeJavaScriptText(request.Requirements))

	indexCode := `import React from "react";
import ReactDOM from "react-dom/client";
import "bootstrap/dist/css/bootstrap.min.css";
import "./index.css";
import App from "./App";

const root = ReactDOM.createRoot(
  document.getElementById("root")
);

root.render(
  <React.StrictMode>
    <App />
  </React.StrictMode>
);
`

	appCSS := `.app-shell {
  min-height: 100vh;
  background: #f5f7fb;
  color: #17233d;
}

.app-header {
  padding: 64px 24px;
  text-align: center;
  color: #ffffff;
  background: linear-gradient(135deg, #123b6d, #0a8f85);
}

.app-content {
  width: min(1100px, calc(100% - 32px));
  margin: 0 auto;
  padding: 48px 0;
}

.app-card {
  padding: 32px;
  background: #ffffff;
  border-radius: 16px;
  box-shadow: 0 16px 40px rgba(23, 35, 61, 0.1);
}
`

	dockerfile := `FROM node:20-alpine AS build
WORKDIR /app
COPY package.json ./
RUN npm install
COPY . .
RUN npm run build

FROM nginx:1.27-alpine
COPY --from=build /app/build /usr/share/nginx/html
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
`

	html := fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="utf-8" />
  <meta
    name="viewport"
    content="width=device-width, initial-scale=1"
  />
  <title>%s</title>
</head>
<body>
  <div id="root"></div>
</body>
</html>
`, projectName)

	return []FileDefinition{
		{
			Path:    "frontend/package.json",
			Content: packageJSON,
		},
		{
			Path:    "frontend/public/index.html",
			Content: html,
		},
		{
			Path:    "frontend/src/index.js",
			Content: indexCode,
		},
		{
			Path:    "frontend/src/index.css",
			Content: "body { margin: 0; font-family: Arial, sans-serif; }\n",
		},
		{
			Path:    "frontend/src/App.js",
			Content: appCode,
		},
		{
			Path:    "frontend/src/App.css",
			Content: appCSS,
		},
		{
			Path:    "frontend/Dockerfile",
			Content: dockerfile,
		},
	}
}

func buildAngularFiles(
	request GenerateProjectRequest,
	projectName string,
) []FileDefinition {
	packageJSON := fmt.Sprintf(`{
  "name": "%s-frontend",
  "version": "1.0.0",
  "scripts": {
    "start": "ng serve",
    "build": "ng build"
  },
  "dependencies": {
    "@angular/animations": "^19.0.0",
    "@angular/common": "^19.0.0",
    "@angular/compiler": "^19.0.0",
    "@angular/core": "^19.0.0",
    "@angular/forms": "^19.0.0",
    "@angular/platform-browser": "^19.0.0",
    "@angular/platform-browser-dynamic": "^19.0.0",
    "@angular/router": "^19.0.0",
    "rxjs": "^7.8.1",
    "tslib": "^2.8.1",
    "zone.js": "^0.15.0"
  },
  "devDependencies": {
    "@angular-devkit/build-angular": "^19.0.0",
    "@angular/cli": "^19.0.0",
    "@angular/compiler-cli": "^19.0.0",
    "typescript": "^5.6.3"
  }
}
`, projectName)

	component := fmt.Sprintf(`import { Component } from "@angular/core";
import { RouterOutlet } from "@angular/router";

@Component({
  selector: "app-root",
  standalone: true,
  imports: [RouterOutlet],
  templateUrl: "./app.component.html",
  styleUrl: "./app.component.css"
})
export class AppComponent {
  projectName = "%s";
  requirements = "%s";
}
`,
		escapeJavaScriptText(projectName),
		escapeJavaScriptText(request.Requirements),
	)

	template := `<div class="app-shell">
  <header class="app-header">
    <h1>{{ projectName }}</h1>
    <p>{{ requirements }}</p>
  </header>

  <main class="app-content">
    <section class="app-card">
      <h2>Application Generated Successfully</h2>
      <p>
        Replace this page with modules based on the application requirements.
      </p>
    </section>
  </main>
</div>

<router-outlet></router-outlet>
`

	mainCode := `import { bootstrapApplication } from "@angular/platform-browser";
import { AppComponent } from "./app/app.component";

bootstrapApplication(AppComponent).catch((error) => {
  console.error(error);
});
`

	dockerfile := `FROM node:20-alpine AS build
WORKDIR /app
COPY package.json ./
RUN npm install
COPY . .
RUN npm run build

FROM nginx:1.27-alpine
COPY --from=build /app/dist /usr/share/nginx/html
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
`

	angularJSON := fmt.Sprintf(`{
  "$schema": "./node_modules/@angular/cli/lib/config/schema.json",
  "version": 1,
  "projects": {
    "%s": {
      "projectType": "application",
      "root": "",
      "sourceRoot": "src",
      "architect": {
        "build": {
          "builder": "@angular-devkit/build-angular:application",
          "options": {
            "outputPath": "dist",
            "index": "src/index.html",
            "browser": "src/main.ts",
            "tsConfig": "tsconfig.app.json",
            "styles": [
              "src/styles.css"
            ]
          }
        },
        "serve": {
          "builder": "@angular-devkit/build-angular:dev-server"
        }
      }
    }
  }
}
`, projectName)

	return []FileDefinition{
		{
			Path:    "frontend/package.json",
			Content: packageJSON,
		},
		{
			Path:    "frontend/angular.json",
			Content: angularJSON,
		},
		{
			Path:    "frontend/src/main.ts",
			Content: mainCode,
		},
		{
			Path: "frontend/src/index.html",
			Content: fmt.Sprintf(
				"<!DOCTYPE html><html><head><title>%s</title></head><body><app-root></app-root></body></html>",
				projectName,
			),
		},
		{
			Path:    "frontend/src/styles.css",
			Content: "body { margin: 0; font-family: Arial, sans-serif; }\n",
		},
		{
			Path:    "frontend/src/app/app.component.ts",
			Content: component,
		},
		{
			Path:    "frontend/src/app/app.component.html",
			Content: template,
		},
		{
			Path: "frontend/src/app/app.component.css",
			Content: `.app-shell {
  min-height: 100vh;
  background: #f5f7fb;
}

.app-header {
  padding: 64px 24px;
  text-align: center;
  color: #ffffff;
  background: linear-gradient(135deg, #123b6d, #0a8f85);
}

.app-content {
  padding: 48px 24px;
}

.app-card {
  max-width: 1000px;
  margin: 0 auto;
  padding: 32px;
  background: #ffffff;
  border-radius: 16px;
}
`,
		},
		{
			Path:    "frontend/Dockerfile",
			Content: dockerfile,
		},
	}
}

func escapeJavaScriptText(value string) string {
	value = strings.ReplaceAll(value, `\`, `\\`)
	value = strings.ReplaceAll(value, `"`, `\"`)
	value = strings.ReplaceAll(value, "\r", " ")
	value = strings.ReplaceAll(value, "\n", " ")
	return value
}
