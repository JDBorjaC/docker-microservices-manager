package internal

import "database/sql"

type Repository struct {
	db *sql.DB
}

func NewRepository() (*Repository, error) {
	db, err := sql.Open("sqlite3", "./data.db")
	if err != nil {
		return nil, err
	}

	r := &Repository{db: db}

	if err := r.Init(); err != nil {
		db.Close()
		return nil, err
	}

	return r, nil
}

func (r *Repository) Close() error {
	return r.db.Close()
}

func (r *Repository) Init() error {
	_, err := r.db.Exec(`
    CREATE TABLE IF NOT EXISTS microservice (
        id INTEGER PRIMARY KEY,
        name TEXT UNIQUE,
        description TEXT,
        image TEXT,
        port INTEGER
    )
    `)
	return err
}
