package db

import "database/sql"

type SQLiteRepository struct {
	db *sql.DB
}

func NewSQLiteRepository(db *sql.DB) *SQLiteRepository {
	return &SQLiteRepository{
		db: db,
	}
}

func (r *SQLiteRepository) Migrate() error {
	query := `
    CREATE TABLE IF NOT EXISTS messages(
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        timestamp timestamp NOT NULL,
        text TEXT NOT NULL        
    );`

	_, err := r.db.Exec(query)
	return err
}

func (r *SQLiteRepository) Append(msg string) error {
	query := `INSERT INTO messages(timestamp, text) VALUES(datetime('now'),?);`
	_, err := r.db.Exec(query, msg)
	return err
}

func (r *SQLiteRepository) Fetch() (*[]string, error) {
	query := `SELECT text FROM messages LIMIT 100;`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}

	var messages []string
	for rows.Next() {
		var msg string
		if err := rows.Scan(&msg); err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}
	return &messages, nil
}
