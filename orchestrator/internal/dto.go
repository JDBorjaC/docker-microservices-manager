package internal

type PullImageRequest struct {
	ImageId string `json:"imageId"`
}

type CreateMicroserviceRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Language    string `json:"language"`
	Code        string `json:"code"`
}

