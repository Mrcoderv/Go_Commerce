package migrations

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"sort"
)

//go:embed *.sql
var MigrationFS embed.FS

// Run applies all embedded .sql migrations in alphabetical order.
func Run(ctx context.Context, db *sql.DB) error {
	entries, err := MigrationFS.ReadDir(".")
	if err != nil {
		return fmt.Errorf("read migrations: %w", err)
	}

	var files []string
	for _, entry := range entries {
		if !entry.IsDir() {
			files = append(files, entry.Name())
		}
	}
	sort.Strings(files)

	for _, file := range files {
		content, err := MigrationFS.ReadFile(file)
		if err != nil {
			return fmt.Errorf("read migration %s: %w", file, err)
		}
		if _, err := db.ExecContext(ctx, string(content)); err != nil {
			return fmt.Errorf("apply migration %s: %w", file, err)
		}
	}
	return nil
}
