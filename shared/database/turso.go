package database

import (
	"database/sql"
	"fmt"

	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

// OpenTurso abre la conexión usando DSN libsql://...turso.io?authToken=...
func OpenTurso(url, token string) (*sql.DB, error) {
	if url == "" || token == "" {
		return nil, fmt.Errorf("missing Turso url or token")
	}
	dsn := fmt.Sprintf("%s?authToken=%s", url, token)
	db, err := sql.Open("libsql", dsn)
	if err != nil {
		return nil, fmt.Errorf("open libsql: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ping libsql: %w", err)
	}
	return db, nil
}

// InitSchema crea la tabla books con ISBN único
func InitSchema(db *sql.DB) error {
	ddl := `
	CREATE TABLE IF NOT EXISTS books (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		author TEXT NOT NULL,
		year INTEGER NOT NULL,
		genre TEXT NOT NULL,
		isbn TEXT NOT NULL UNIQUE,
		created_at TIMESTAMP NOT NULL,
		updated_at TIMESTAMP NOT NULL
	);`
	_, err := db.Exec(ddl)
	return err
}
