package internal

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

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
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT UNIQUE,
        description TEXT,
        language TEXT,
        image TEXT,
        container_id TEXT,
        status TEXT DEFAULT 'created',
		code TEXT,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP
    )
    `)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) InsertMicroservice(ms *Microservice) error {
	query := `INSERT INTO microservice (name, description, language, image, container_id, status, code, created_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	result, err := r.db.Exec(query, ms.Name, ms.Description, ms.Language, ms.Image, ms.ContainerId, ms.Status, ms.Code, ms.CreatedAt)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	ms.Id = int(id)
	return nil
}

func (r *Repository) GetAllMicroservices() ([]Microservice, error) {
	rows, err := r.db.Query(`SELECT id, name, description, language, image, container_id, status, created_at FROM microservice ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var microservices []Microservice
	for rows.Next() {
		var ms Microservice
		var containerId sql.NullString
		err := rows.Scan(&ms.Id, &ms.Name, &ms.Description, &ms.Language, &ms.Image, &containerId, &ms.Status, &ms.CreatedAt)
		if err != nil {
			return nil, err
		}
		ms.ContainerId = containerId.String
		microservices = append(microservices, ms)
	}
	return microservices, nil
}

var allowedColumns = map[string]bool{
	"id": true, "name": true, "container_id": true, "status": true,
	"description": true, "image": true, "language": true, "code": true,
}

func (r *Repository) GetMicroserviceBy(column string, value any) (*Microservice, error) {
	if !allowedColumns[column] {
		return nil, fmt.Errorf("invalid column: %s", column)
	}

	var ms Microservice
	var containerId sql.NullString
	query := fmt.Sprintf(
		`SELECT id, name, description, language, image, container_id, status, code, created_at FROM microservice WHERE %s = ?`, column,
	)
	err := r.db.QueryRow(query, value).Scan(
		&ms.Id, &ms.Name, &ms.Description, &ms.Language, &ms.Image, &containerId, &ms.Status, &ms.Code, &ms.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	ms.ContainerId = containerId.String
	return &ms, nil
}

// Updates multiple allowed columns for a given microservice ID in a single query.
func (r *Repository) UpdateMicroservice(id int, updates map[string]any) error {
	if len(updates) == 0 {
		return fmt.Errorf("no updates provided")
	}

	query := "UPDATE microservice SET "
	var args []any
	var setClauses []string

	for column, value := range updates {
		if !allowedColumns[column] {
			return fmt.Errorf("invalid column: %s", column)
		}
		setClauses = append(setClauses, fmt.Sprintf("%s = ?", column))
		args = append(args, value)
	}

	query += strings.Join(setClauses, ", ") + " WHERE id = ?"
	args = append(args, id)

	_, err := r.db.Exec(query, args...)
	return err
}

func (r *Repository) DeleteMicroservice(id int) error {
	_, err := r.db.Exec(`DELETE FROM microservice WHERE id = ?`, id)
	return err
}
