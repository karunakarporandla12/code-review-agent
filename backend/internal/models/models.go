package models

import "time"

type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"createdAt"`
}

type Project struct {
	ID          string    `json:"id"`
	OwnerID     string    `json:"ownerId"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
}

type Review struct {
	ID            string     `json:"id"`
	ProjectID     string     `json:"projectId"`
	RequestedBy   string     `json:"requestedBy"`
	SourceType    string     `json:"sourceType"`
	SourceName    string     `json:"sourceName"`
	BranchName    *string    `json:"branchName,omitempty"`
	CommitSHA     *string    `json:"commitSha,omitempty"`
	Status        string     `json:"status"`
	Summary       string     `json:"summary"`
	FilesScanned  int        `json:"filesScanned"`
	CriticalCount int        `json:"criticalCount"`
	HighCount     int        `json:"highCount"`
	MediumCount   int        `json:"mediumCount"`
	LowCount      int        `json:"lowCount"`
	ErrorMessage  string     `json:"errorMessage"`
	CreatedAt     time.Time  `json:"createdAt"`
	CompletedAt   *time.Time `json:"completedAt,omitempty"`
}

type Finding struct {
	ID             string    `json:"id"`
	ReviewID       string    `json:"reviewId"`
	Severity       string    `json:"severity"`
	Category       string    `json:"category"`
	RuleID         string    `json:"ruleId"`
	FilePath       string    `json:"filePath"`
	LineNumber     int       `json:"lineNumber"`
	Title          string    `json:"title"`
	Description    string    `json:"description"`
	Recommendation string    `json:"recommendation"`
	CodeExcerpt    string    `json:"codeExcerpt"`
	Source         string    `json:"source"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"createdAt"`
}
