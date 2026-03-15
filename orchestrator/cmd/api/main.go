package main

import (
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	docs "msmanager/orchestrator/docs"
	"msmanager/orchestrator/internal"
)

func main() {

	repository, err := internal.NewRepository()
	if err != nil {
		panic(err)
	}
	defer repository.Close()

	dockerClient, err := internal.NewDockerClient()
	if err != nil {
		panic(err)
	}
	defer dockerClient.Close()

	service := internal.NewService(dockerClient, repository)
	handler := internal.NewHandler(service)

	r := gin.Default()
	docs.SwaggerInfo.BasePath = ""

	r.POST("/images/pull", handler.PullImage)
	r.POST("/microservices", handler.CreateMicroservice)
	r.GET("/microservices", handler.GetMicroservices)
	r.GET("/microservices/stream/:id", handler.StreamMicroserviceLogs)

	r.GET("/api/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	r.Run(":8080")

}
