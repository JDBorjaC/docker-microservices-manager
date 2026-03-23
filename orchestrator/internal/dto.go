package internal

type CreateMicroserviceRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Language    string `json:"language" binding:"required"`
	Code        string `json:"code" binding:"required"`
}

type UpdateMicroserviceRequest struct {
	Code        *string `json:"code"`
	Description *string `json:"description"`
}
