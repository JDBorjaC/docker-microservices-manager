package internal

type PullImageRequest struct {
	ImageId string `json:"imageId"`
}

type CreateMicroserviceRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Image       string `json:"image"`
	Code        string `json:"code"`
}

type MicroserviceResponse struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Image       string `json:"image"`
	Port        int    `json:"port"`
}
