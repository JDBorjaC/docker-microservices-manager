package internal

import (
	"database/sql"

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
        image TEXT,
        container_id TEXT,
        status TEXT DEFAULT 'created',
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP
    )
    `)
	return err
}

func (r *Repository) InsertMicroservice(ms *Microservice) error {
	query := `INSERT INTO microservice (name, description, image, container_id, status, created_at) VALUES (?, ?, ?, ?, ?, ?)`
	result, err := r.db.Exec(query, ms.Name, ms.Description, ms.Image, ms.ContainerId, ms.Status, ms.CreatedAt)
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

func (r *Repository) UpdateMicroserviceContainerID(id int, containerId string) error {
	_, err := r.db.Exec(`UPDATE microservice SET container_id = ? WHERE id = ?`, containerId, id)
	return err
}

func (r *Repository) UpdateMicroserviceStatus(id int, status string) error {
	_, err := r.db.Exec(`UPDATE microservice SET status = ? WHERE id = ?`, status, id)
	return err
}

func (r *Repository) GetAllMicroservices() ([]Microservice, error) {
	rows, err := r.db.Query(`SELECT id, name, description, image, container_id, status, created_at FROM microservice ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var microservices []Microservice
	for rows.Next() {
		var ms Microservice
		var containerId sql.NullString
		err := rows.Scan(&ms.Id, &ms.Name, &ms.Description, &ms.Image, &containerId, &ms.Status, &ms.CreatedAt)
		if err != nil {
			return nil, err
		}
		ms.ContainerId = containerId.String
		microservices = append(microservices, ms)
	}
	return microservices, nil
}

func (r *Repository) GetMicroserviceContainerID(id int) (string, error) {
	var containerId sql.NullString
	err := r.db.QueryRow(`SELECT container_id FROM microservice WHERE id = ?`, id).Scan(&containerId)
	if err != nil {
		return "", err
	}
	return containerId.String, nil
}

func (r *Repository) GetMicroserviceByName(name string) (*Microservice, error) {
	var ms Microservice
	var containerId sql.NullString
	err := r.db.QueryRow(`SELECT id, name, description, image, container_id, status, created_at FROM microservice WHERE name = ?`, name).Scan(&ms.Id, &ms.Name, &ms.Description, &ms.Image, &containerId, &ms.Status, &ms.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Not found
		}
		return nil, err
	}
	ms.ContainerId = containerId.String
	return &ms, nil
}
