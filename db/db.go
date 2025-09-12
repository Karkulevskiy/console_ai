package db

import (
	"context"
	"database/sql"
	"errors"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

const dbFile = "models.db"

var defaultModels = []string{
	// Google Gemini
	"googleai/gemini-2.5-flash",
	"googleai/gemini-1.5-pro",
	"googleai/gemini-pro-vision",
	"googleai/gemini-1.0-pro",
	"googleai/gemini-1.0-pro-vision",
	// OpenAI GPT
	"openai/gpt-4o",
	"openai/gpt-4-turbo",
	"openai/gpt-4",
	"openai/gpt-3.5-turbo",
	// Anthropic Claude
	"anthropic/claude-3-opus",
	"anthropic/claude-3-sonnet",
	"anthropic/claude-3-haiku",
	// Mistral
	"mistral/mistral-large",
	"mistral/mistral-medium",
	"mistral/mistral-small",
	// Llama
	"meta/llama-3-70b",
	"meta/llama-3-8b",
	// Cohere
	"cohere/command-r-plus",
	"cohere/command-r",
}

func seedModels(ctx context.Context) error {
	for _, name := range defaultModels {
		_ = AddModel(ctx, name)
	}
	return nil
}

func InitDB() error {
	if _, err := os.Stat(dbFile); errors.Is(err, os.ErrNotExist) {
		file, err := os.Create(dbFile)
		if err != nil {
			return err
		}
		file.Close()
	}

	var err error
	db, err = sql.Open("sqlite3", dbFile)
	if err != nil {
		return err
	}

	createTable := `
    CREATE TABLE IF NOT EXISTS models (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT UNIQUE NOT NULL
    );`
	_, err = db.Exec(createTable)
	if err != nil {
		return err
	}

	return seedModels(context.Background())
}

func GetAvailableModels(ctx context.Context) ([]string, error) {
	rows, err := db.QueryContext(ctx, "SELECT name FROM models")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var models []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		models = append(models, name)
	}
	return models, nil
}

func AddModel(ctx context.Context, name string) error {
	_, err := db.ExecContext(ctx, "INSERT INTO models(name) VALUES(?)", name)
	return err
}

func DeleteModel(ctx context.Context, name string) error {
	_, err := db.ExecContext(ctx, "DELETE FROM models WHERE name = ?", name)
	return err
}

func UpdateModel(ctx context.Context, oldName, newName string) error {
	_, err := db.ExecContext(ctx, "UPDATE models SET name = ? WHERE name = ?", newName, oldName)
	return err
}
