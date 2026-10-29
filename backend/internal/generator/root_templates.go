package generator

import "fmt"

func buildRootFiles(
	request GenerateProjectRequest,
	projectName string,
) []FileDefinition {
	readme := fmt.Sprintf(`# %s

## Requirements

%s

## Technology Stack

Frontend: %s
Backend: %s
Database: %s

## Project Structure

frontend
backend
k8s
docker-compose.yml

## Local Development

Run the frontend and backend according to their individual README files.
`,
		projectName,
		request.Requirements,
		request.FrontendTechnology,
		request.BackendTechnology,
		request.DatabaseTechnology,
	)

	dockerCompose := buildDockerCompose(request)

	return []FileDefinition{
		{
			Path:    "README.md",
			Content: readme,
		},
		{
			Path:    "docker-compose.yml",
			Content: dockerCompose,
		},
		{
			Path: "k8s/frontend-deployment.yaml",
			Content: fmt.Sprintf(`apiVersion: apps/v1
kind: Deployment
metadata:
  name: %s-frontend
spec:
  replicas: 1
  selector:
    matchLabels:
      app: %s-frontend
  template:
    metadata:
      labels:
        app: %s-frontend
    spec:
      containers:
        - name: frontend
          image: %s-frontend:latest
          imagePullPolicy: IfNotPresent
          ports:
            - containerPort: 80
---
apiVersion: v1
kind: Service
metadata:
  name: %s-frontend-service
spec:
  selector:
    app: %s-frontend
  ports:
    - port: 80
      targetPort: 80
`,
				projectName,
				projectName,
				projectName,
				projectName,
				projectName,
				projectName,
			),
		},
		{
			Path: "k8s/backend-deployment.yaml",
			Content: fmt.Sprintf(`apiVersion: apps/v1
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
          imagePullPolicy: IfNotPresent
          ports:
            - containerPort: 8080
---
apiVersion: v1
kind: Service
metadata:
  name: %s-backend-service
spec:
  selector:
    app: %s-backend
  ports:
    - port: 8080
      targetPort: 8080
`,
				projectName,
				projectName,
				projectName,
				projectName,
				projectName,
				projectName,
			),
		},
	}
}

func buildDockerCompose(request GenerateProjectRequest) string {
	databaseService := ""

	switch request.DatabaseTechnology {
	case "postgresql":
		databaseService = `  database:
    image: postgres:17
    environment:
      POSTGRES_DB: appdb
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
    ports:
      - "5432:5432"
    volumes:
      - database_data:/var/lib/postgresql/data
`
	case "mysql":
		databaseService = `  database:
    image: mysql:8.4
    environment:
      MYSQL_DATABASE: appdb
      MYSQL_ROOT_PASSWORD: root
    ports:
      - "3306:3306"
    volumes:
      - database_data:/var/lib/mysql
`
	case "mongodb":
		databaseService = `  database:
    image: mongo:8
    ports:
      - "27017:27017"
    volumes:
      - database_data:/data/db
`
	}

	volumes := ""

	if request.DatabaseTechnology != "none" {
		volumes = `
volumes:
  database_data:
`
	}

	return fmt.Sprintf(`services:
  frontend:
    build:
      context: ./frontend
    ports:
      - "3000:80"
    depends_on:
      - backend

  backend:
    build:
      context: ./backend
    ports:
      - "8080:8080"
%s%s`,
		databaseService,
		volumes,
	)
}
