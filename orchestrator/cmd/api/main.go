package main

import (
	"msmanager/orchestrator/internal"

	"github.com/gin-gonic/gin"

	docs "msmanager/orchestrator/docs"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	dockerClient, err := internal.NewDockerClient()
	if err != nil {
		panic(err)
	}
	defer dockerClient.Close()

	service := internal.NewService(dockerClient)
	handler := internal.NewHandler(service)

	r := gin.Default()
	docs.SwaggerInfo.BasePath = "/api"

	r.POST("/images/pull", handler.PullImage)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	r.Run(":8080")

}
