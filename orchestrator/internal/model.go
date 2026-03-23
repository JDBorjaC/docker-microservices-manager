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
	Id          int       `json:"id" example:"1"`
	Name        string    `json:"name" example:"api-users"`
	Description string    `json:"description" example:"Servicio de autenticación"`
	Language    string    `json:"language" example:"express"`
	Image       string    `json:"image" example:"msm-api-users:latest"`
	ContainerId string    `json:"container_id" example:"a1b2c3d4e5f6"`
	Status      string    `json:"status" example:"running"`
	CreatedAt   time.Time `json:"created_at" example:"2025-01-01T15:04:05Z"`
}
