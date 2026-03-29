package main

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
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

	backgroundCtx := context.Background()
	// 1. Apagar/Prender issue: Sincronizar el estado de DB vs la realidad Docker al Boot
	if err := service.SyncStateOnStartup(backgroundCtx); err != nil {
		fmt.Printf("Error during startup sync: %v\n", err)
	}

	// 2. Encender la escucha infinita y pasiva de Eventos Docker
	service.StartReconciliationLoop(backgroundCtx)

	handler := internal.NewHandler(service)

	r := gin.Default()
	docs.SwaggerInfo.BasePath = ""

	r.Use(cors.Default())

	r.POST("/microservices", handler.CreateMicroservice)
	r.GET("/microservices", handler.GetMicroservices)
	r.GET("/microservices/:id", handler.GetMicroserviceByID)
	r.GET("/microservices/stream/:id", handler.StreamMicroserviceLogs)
	r.PUT("/microservices/:id", handler.UpdateMicroservice)
	r.PATCH("/microservices/stop/:id", handler.StopMicroservice)
	r.PATCH("/microservices/start/:id", handler.StartMicroservice)
	r.GET("/microservices/status/events", handler.StreamStatusUpdates)
	r.DELETE("/microservices/:id", handler.DeleteMicroservice)
	r.DELETE("/microservices", handler.DeleteAllMicroservices)

	r.GET("/api/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	r.Run(":8080")

}
