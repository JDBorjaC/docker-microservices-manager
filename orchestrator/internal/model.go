package internal

type Microservice struct {
	Id          int
	Name        string
	Description string
	Image       string
	Port        int    // host port mapped to the container
	ContainerId string // Docker container ID
}
