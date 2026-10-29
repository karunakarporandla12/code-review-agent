// package database

// import (
// 	"context"
// 	"fmt"

// 	"github.com/jackc/pgx/v5/pgxpool"
// )

// func Open(ctx context.Context, databaseURL string) (*pgxpool.Pool, error) {
// 	pool, err := pgxpool.New(ctx, databaseURL)
// 	if err != nil {
// 		return nil, fmt.Errorf("create PostgreSQL pool: %w", err)
// 	}
// 	if err := pool.Ping(ctx); err != nil {
// 		pool.Close()
// 		return nil, fmt.Errorf("connect PostgreSQL: %w", err)
// 	}
// 	return pool, nil
// }

// func EnsureSchema(ctx context.Context, db *pgxpool.Pool) error {
// 	if db == nil {
// 		return nil
// 	}

// 	statements := []string{
// 		`CREATE TABLE IF NOT EXISTS users (id UUID PRIMARY KEY DEFAULT gen_random_uuid(), name TEXT NOT NULL, email TEXT NOT NULL UNIQUE, password_hash TEXT NOT NULL, role TEXT NOT NULL DEFAULT 'user', created_at TIMESTAMPTZ NOT NULL DEFAULT NOW())`,
// 		`CREATE TABLE IF NOT EXISTS projects (id UUID PRIMARY KEY DEFAULT gen_random_uuid(), owner_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE, name TEXT NOT NULL, description TEXT DEFAULT '', created_at TIMESTAMPTZ NOT NULL DEFAULT NOW())`,
// 		`CREATE TABLE IF NOT EXISTS project_members (project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE, user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE, member_role TEXT NOT NULL DEFAULT 'member', PRIMARY KEY(project_id, user_id))`,
// 		`CREATE TABLE IF NOT EXISTS reviews (id UUID PRIMARY KEY DEFAULT gen_random_uuid(), project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE, requested_by UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE, source_type TEXT NOT NULL, source_name TEXT NOT NULL, branch_name TEXT, commit_sha TEXT, status TEXT NOT NULL DEFAULT 'queued', summary TEXT, files_scanned INTEGER NOT NULL DEFAULT 0, critical_count INTEGER NOT NULL DEFAULT 0, high_count INTEGER NOT NULL DEFAULT 0, medium_count INTEGER NOT NULL DEFAULT 0, low_count INTEGER NOT NULL DEFAULT 0, error_message TEXT DEFAULT '', created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(), completed_at TIMESTAMPTZ)`,
// 		`CREATE TABLE IF NOT EXISTS findings (id UUID PRIMARY KEY DEFAULT gen_random_uuid(), review_id UUID NOT NULL REFERENCES reviews(id) ON DELETE CASCADE, severity TEXT NOT NULL, category TEXT NOT NULL, rule_id TEXT NOT NULL, file_path TEXT NOT NULL, line_number INTEGER NOT NULL DEFAULT 0, title TEXT NOT NULL, description TEXT NOT NULL, recommendation TEXT NOT NULL, code_excerpt TEXT, source TEXT NOT NULL, status TEXT NOT NULL DEFAULT 'open', created_at TIMESTAMPTZ NOT NULL DEFAULT NOW())`,
// 		`CREATE TABLE IF NOT EXISTS webhook_deliveries (delivery_id TEXT PRIMARY KEY, event_name TEXT NOT NULL, created_at TIMESTAMPTZ NOT NULL DEFAULT NOW())`,
// 		`CREATE TABLE IF NOT EXISTS generated_projects (id UUID PRIMARY KEY DEFAULT gen_random_uuid(), owner_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE, project_name TEXT NOT NULL, slug TEXT NOT NULL, description TEXT DEFAULT '', requirements TEXT NOT NULL, root_path TEXT NOT NULL, git_url TEXT, git_branch TEXT DEFAULT 'main', status TEXT NOT NULL DEFAULT 'generated', created_at TIMESTAMPTZ NOT NULL DEFAULT NOW())`,
// 		`CREATE TABLE IF NOT EXISTS generated_project_files (id UUID PRIMARY KEY DEFAULT gen_random_uuid(), generated_project_id UUID NOT NULL REFERENCES generated_projects(id) ON DELETE CASCADE, file_name TEXT NOT NULL, file_path TEXT NOT NULL, file_content TEXT NOT NULL DEFAULT '', created_at TIMESTAMPTZ NOT NULL DEFAULT NOW())`,
// 	}

// 	for _, statement := range statements {
// 		if _, err := db.Exec(ctx, statement); err != nil {
// 			return fmt.Errorf("prepare schema: %w", err)
// 		}
// 	}
// 	return nil
// }

package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Open(ctx context.Context, databaseURL string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return nil, fmt.Errorf("create PostgreSQL pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("connect PostgreSQL: %w", err)
	}

	return pool, nil
}

func EnsureSchema(ctx context.Context, db *pgxpool.Pool) error {
	if db == nil {
		return nil
	}

	statements := []string{
		`CREATE EXTENSION IF NOT EXISTS pgcrypto`,

		`CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			name TEXT NOT NULL,
			email TEXT NOT NULL UNIQUE,
			password_hash TEXT NOT NULL,
			role TEXT NOT NULL DEFAULT 'user',
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,

		`CREATE TABLE IF NOT EXISTS projects (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			owner_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			name TEXT NOT NULL,
			description TEXT DEFAULT '',
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,

		`CREATE TABLE IF NOT EXISTS project_members (
			project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
			user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			member_role TEXT NOT NULL DEFAULT 'member',
			PRIMARY KEY(project_id, user_id)
		)`,

		`CREATE TABLE IF NOT EXISTS reviews (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
			requested_by UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			source_type TEXT NOT NULL,
			source_name TEXT NOT NULL,
			branch_name TEXT,
			commit_sha TEXT,
			status TEXT NOT NULL DEFAULT 'queued',
			summary TEXT,
			files_scanned INTEGER NOT NULL DEFAULT 0,
			critical_count INTEGER NOT NULL DEFAULT 0,
			high_count INTEGER NOT NULL DEFAULT 0,
			medium_count INTEGER NOT NULL DEFAULT 0,
			low_count INTEGER NOT NULL DEFAULT 0,
			error_message TEXT DEFAULT '',
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			completed_at TIMESTAMPTZ
		)`,

		`CREATE TABLE IF NOT EXISTS findings (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			review_id UUID NOT NULL REFERENCES reviews(id) ON DELETE CASCADE,
			severity TEXT NOT NULL,
			category TEXT NOT NULL,
			rule_id TEXT NOT NULL,
			file_path TEXT NOT NULL,
			line_number INTEGER NOT NULL DEFAULT 0,
			title TEXT NOT NULL,
			description TEXT NOT NULL,
			recommendation TEXT NOT NULL,
			code_excerpt TEXT,
			source TEXT NOT NULL,
			status TEXT NOT NULL DEFAULT 'open',
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,

		`CREATE TABLE IF NOT EXISTS webhook_deliveries (
			delivery_id TEXT PRIMARY KEY,
			event_name TEXT NOT NULL,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,

		`CREATE TABLE IF NOT EXISTS generated_projects (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			owner_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			project_name TEXT NOT NULL,
			slug TEXT NOT NULL,
			description TEXT DEFAULT '',
			requirements TEXT NOT NULL,
			frontend_technology TEXT NOT NULL DEFAULT 'react',
			backend_technology TEXT NOT NULL DEFAULT 'golang',
			database_technology TEXT NOT NULL DEFAULT 'postgresql',
			root_path TEXT NOT NULL,
			git_url TEXT,
			git_branch TEXT NOT NULL DEFAULT 'main',
			status TEXT NOT NULL DEFAULT 'generated',
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,

		`CREATE TABLE IF NOT EXISTS generated_project_files (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			generated_project_id UUID NOT NULL REFERENCES generated_projects(id) ON DELETE CASCADE,
			file_name TEXT NOT NULL,
			file_path TEXT NOT NULL,
			file_content TEXT NOT NULL DEFAULT '',
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			UNIQUE(generated_project_id, file_path)
		)`,

		`ALTER TABLE generated_projects
			ADD COLUMN IF NOT EXISTS frontend_technology TEXT NOT NULL DEFAULT 'react'`,

		`ALTER TABLE generated_projects
			ADD COLUMN IF NOT EXISTS backend_technology TEXT NOT NULL DEFAULT 'golang'`,

		`ALTER TABLE generated_projects
			ADD COLUMN IF NOT EXISTS database_technology TEXT NOT NULL DEFAULT 'postgresql'`,

		`ALTER TABLE generated_projects
			ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()`,

		`ALTER TABLE generated_project_files
			ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()`,

		`CREATE INDEX IF NOT EXISTS idx_projects_owner_id
			ON projects(owner_id)`,

		`CREATE INDEX IF NOT EXISTS idx_reviews_project_id
			ON reviews(project_id)`,

		`CREATE INDEX IF NOT EXISTS idx_findings_review_id
			ON findings(review_id)`,

		`CREATE INDEX IF NOT EXISTS idx_generated_projects_owner_id
			ON generated_projects(owner_id)`,

		`CREATE INDEX IF NOT EXISTS idx_generated_projects_created_at
			ON generated_projects(created_at DESC)`,

		`CREATE INDEX IF NOT EXISTS idx_generated_project_files_project_id
			ON generated_project_files(generated_project_id)`,
	}

	for _, statement := range statements {
		if _, err := db.Exec(ctx, statement); err != nil {
			return fmt.Errorf("prepare schema: %w", err)
		}
	}

	return nil
}
