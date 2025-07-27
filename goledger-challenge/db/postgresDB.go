package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq" // PostgreSQL driver
)

func InitPostgres(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("DB open error: %w", err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS storage (
        id SERIAL PRIMARY KEY,
        value TEXT
    )`)
	if err != nil {
		return nil, fmt.Errorf("DB table create error: %w", err)
	}

	var count int
	err = db.QueryRow(`SELECT COUNT(*) FROM storage WHERE id = 1`).Scan(&count)
	if err != nil {
		return nil, fmt.Errorf("DB row count error: %w", err)
	}
	if count == 0 {
		_, err = db.Exec(`INSERT INTO storage (id, value) VALUES (1, '0')`)
		if err != nil {
			return nil, fmt.Errorf("DB row insert error: %w", err)
		}
	}

	return db, nil
}
