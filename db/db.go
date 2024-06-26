package database

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

type SQLiteRepository struct {
	db *sql.DB
}

//go:generate mockery --name RepoOperations
type RepoOperations interface {
	CreateTable() error
	Append(msg string) error
	Fetch() (*[]string, error)
}

type Service struct {
	Storage RepoOperations
}

const fileName = "db/chat.db"

func NewSQLiteRepository() *SQLiteRepository {
	db, err := sql.Open("sqlite3", fileName)
	if err != nil {
		fmt.Println("cant open db", err)
		panic(err)
	}
	return &SQLiteRepository{
		db: db,
	}
}

func NewService(storage RepoOperations) *Service {
	return &Service{
		Storage: storage,
	}
}

func (r *SQLiteRepository) CreateTable() error {
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
