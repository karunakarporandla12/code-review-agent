package generator

import (
	"fmt"
	"strings"
)

func buildBackendFiles(
	request GenerateProjectRequest,
	projectName string,
) ([]FileDefinition, error) {
	switch strings.ToLower(request.BackendTechnology) {
	case "golang":
		return buildGoBackendFiles(request, projectName), nil
	case "java":
		return buildJavaBackendFiles(request, projectName), nil
	default:
		return nil, fmt.Errorf(
			"unsupported backend technology: %s",
			request.BackendTechnology,
		)
	}
}

func buildGoBackendFiles(
	request GenerateProjectRequest,
	projectName string,
) []FileDefinition {
	moduleName := strings.ReplaceAll(projectName, "-", "")

	mainCode := fmt.Sprintf(`package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type HealthResponse struct {
	Status       string `+"`json:\"status\"`"+`
	Project      string `+"`json:\"project\"`"+`
	Requirements string `+"`json:\"requirements\"`"+`
	Database     string `+"`json:\"database\"`"+`
}

func main() {
	http.HandleFunc("/api/health", func(
		writer http.ResponseWriter,
		request *http.Request,
	) {
		writer.Header().Set("Content-Type", "application/json")

		response := HealthResponse{
			Status:       "running",
			Project:      "%s",
			Requirements: "%s",
			Database:     "%s",
		}

		if err := json.NewEncoder(writer).Encode(response); err != nil {
			http.Error(
				writer,
				"failed to encode response",
				http.StatusInternalServerError,
			)
		}
	})

	log.Println("API listening on http://localhost:8080")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
`,
		escapeGoString(projectName),
		escapeGoString(request.Requirements),
		escapeGoString(request.DatabaseTechnology),
	)

	goMod := fmt.Sprintf(`module %s

go 1.23
`, moduleName)

	dockerfile := `FROM golang:1.23-alpine AS build
WORKDIR /app
COPY go.mod ./
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/api

FROM alpine:3.21
WORKDIR /app
COPY --from=build /app/server /app/server
EXPOSE 8080
CMD ["/app/server"]
`

	envExample := buildEnvironmentExample(request.DatabaseTechnology)

	return []FileDefinition{
		{
			Path:    "backend/go.mod",
			Content: goMod,
		},
		{
			Path:    "backend/cmd/api/main.go",
			Content: mainCode,
		},
		{
			Path:    "backend/.env.example",
			Content: envExample,
		},
		{
			Path:    "backend/Dockerfile",
			Content: dockerfile,
		},
	}
}

func buildJavaBackendFiles(
	request GenerateProjectRequest,
	projectName string,
) []FileDefinition {
	artifactID := strings.ReplaceAll(projectName, "_", "-")
	className := toJavaClassName(projectName)
	packageName := "com.generated." +
		strings.ReplaceAll(strings.ToLower(className), "-", "")

	pomXML := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<project xmlns="http://maven.apache.org/POM/4.0.0">
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

  <build>
    <plugins>
      <plugin>
        <groupId>org.springframework.boot</groupId>
        <artifactId>spring-boot-maven-plugin</artifactId>
      </plugin>
    </plugins>
  </build>
</project>
`, artifactID)

	mainCode := fmt.Sprintf(`package %s;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;

@SpringBootApplication
public class %sApplication {
    public static void main(String[] args) {
        SpringApplication.run(%sApplication.class, args);
    }
}
`, packageName, className, className)

	controllerCode := fmt.Sprintf(`package %s.controller;

import java.util.Map;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RestController;

@RestController
public class HealthController {
    @GetMapping("/api/health")
    public Map<String, String> health() {
        return Map.of(
            "status", "running",
            "project", "%s",
            "requirements", "%s",
            "database", "%s"
        );
    }
}
`,
		packageName,
		escapeJavaString(projectName),
		escapeJavaString(request.Requirements),
		escapeJavaString(request.DatabaseTechnology),
	)

	packagePath := strings.ReplaceAll(packageName, ".", "/")

	dockerfile := `FROM maven:3.9-eclipse-temurin-21 AS build
WORKDIR /app
COPY pom.xml ./
COPY src ./src
RUN mvn clean package -DskipTests

FROM eclipse-temurin:21-jre
WORKDIR /app
COPY --from=build /app/target/*.jar /app/application.jar
EXPOSE 8080
CMD ["java", "-jar", "/app/application.jar"]
`

	return []FileDefinition{
		{
			Path:    "backend/pom.xml",
			Content: pomXML,
		},
		{
			Path: "backend/src/main/java/" +
				packagePath +
				"/" +
				className +
				"Application.java",
			Content: mainCode,
		},
		{
			Path: "backend/src/main/java/" +
				packagePath +
				"/controller/HealthController.java",
			Content: controllerCode,
		},
		{
			Path: "backend/src/main/resources/application.properties",
			Content: `server.port=8080
spring.application.name=generated-application
`,
		},
		{
			Path:    "backend/.env.example",
			Content: buildEnvironmentExample(request.DatabaseTechnology),
		},
		{
			Path:    "backend/Dockerfile",
			Content: dockerfile,
		},
	}
}

func buildEnvironmentExample(databaseTechnology string) string {
	switch databaseTechnology {
	case "postgresql":
		return `DATABASE_TYPE=postgresql
DATABASE_HOST=database
DATABASE_PORT=5432
DATABASE_NAME=appdb
DATABASE_USER=postgres
DATABASE_PASSWORD=postgres
`
	case "mysql":
		return `DATABASE_TYPE=mysql
DATABASE_HOST=database
DATABASE_PORT=3306
DATABASE_NAME=appdb
DATABASE_USER=root
DATABASE_PASSWORD=root
`
	case "mongodb":
		return `DATABASE_TYPE=mongodb
DATABASE_URL=mongodb://database:27017/appdb
`
	default:
		return "DATABASE_TYPE=none\n"
	}
}

func toJavaClassName(value string) string {
	parts := strings.FieldsFunc(value, func(character rune) bool {
		return character == '-' ||
			character == '_' ||
			character == ' '
	})

	result := ""

	for _, part := range parts {
		if part == "" {
			continue
		}

		result += strings.ToUpper(part[:1]) + part[1:]
	}

	if result == "" {
		return "GeneratedProject"
	}

	return result
}

func escapeGoString(value string) string {
	value = strings.ReplaceAll(value, `\`, `\\`)
	value = strings.ReplaceAll(value, `"`, `\"`)
	value = strings.ReplaceAll(value, "\r", " ")
	value = strings.ReplaceAll(value, "\n", " ")
	return value
}

func escapeJavaString(value string) string {
	value = strings.ReplaceAll(value, `\`, `\\`)
	value = strings.ReplaceAll(value, `"`, `\"`)
	value = strings.ReplaceAll(value, "\r", " ")
	value = strings.ReplaceAll(value, "\n", " ")
	return value
}
