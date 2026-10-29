package generator

type GenerateProjectRequest struct {
	ProjectName        string `json:"projectName" binding:"required"`
	Requirements       string `json:"requirements" binding:"required"`
	FrontendTechnology string `json:"frontendTechnology" binding:"required"`
	BackendTechnology  string `json:"backendTechnology" binding:"required"`
	DatabaseTechnology string `json:"databaseTechnology" binding:"required"`
	GitURL             string `json:"gitUrl"`
	GitBranch          string `json:"gitBranch"`
}

type UpdateGeneratedFileRequest struct {
	Path    string `json:"path" binding:"required"`
	Content string `json:"content"`
}

type SyncGitRequest struct {
	GitURL    string `json:"gitUrl" binding:"required"`
	GitBranch string `json:"gitBranch"`
}

type GeneratedProjectResponse struct {
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
