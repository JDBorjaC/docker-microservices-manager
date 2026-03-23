package internal

import "time"

const (
	StatusCreated = "created"
	StatusRunning = "running"
	StatusFailed  = "failed"
	StatusStopped = "stopped"

	LangFlask   = "flask"
	LangExpress = "express"
	LangGin     = "gin"
	LangCargo   = "cargo"
)

type Microservice struct {
	Id          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Language    string    `json:"language"`
	Image       string    `json:"image"`
	ContainerId string    `json:"container_id"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}
