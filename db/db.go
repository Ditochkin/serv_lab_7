package db

import (
	"database/sql"
	"db_lab7/config"
	_ "github.com/mattn/go-sqlite3"
)

// Store ...
type Store struct {
	db  *sql.DB
	dsn string
}

// New ...
func New(config *config.Config) *Store {
	return &Store{
		dsn: config.DSN,
	}
}

// Open ...
func (s *Store) Open() error {
	db, err := sql.Open("sqlite3", s.dsn)
	if err != nil {
		return err
	}

	if err := db.Ping(); err != nil {
		return err
	}

	s.db = db

	return nil
}

// Query ...
func (s *Store) Query(querySTR string, args ...any) (*sql.Rows, error) {
	return s.db.Query(querySTR, args...)
}

// Exec ...
func (s *Store) Exec(querySTR string, args ...any) (sql.Result, error) {
	return s.db.Exec(querySTR, args...)
}

// Close ...
func (s *Store) Close() {
	s.db.Close()
}
