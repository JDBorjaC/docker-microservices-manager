package main

import (
	"msmanager/orchestrator/internal"

	"github.com/gin-gonic/gin"
)

func main() {
	dockerClient, err := internal.NewDockerClient()
	if err != nil {
		panic(err)
	}

	service := internal.NewService(dockerClient)
	handler := internal.NewHandler(service)

	r := gin.Default()
	r.POST("/images/pull", handler.PullImage)
	r.Run(":8080")

}
