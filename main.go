package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/kelseyhightower/envconfig"
	"github.com/more-than-code/deploybot-service-api/api"
	"github.com/more-than-code/deploybot-service-api/middleware"
)

type Config struct {
	ServerPort int `envconfig:"SERVER_PORT"`
}

func main() {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		panic(err)
	}

	g := gin.Default()

	authorized := g.Group("/")

	authorized.Use(middleware.AuthRequired())
	{
		api := api.NewApi()
		authorized.GET("/pipelines", api.GetPipelines())
		authorized.GET("/pipeline/:name", api.GetPipeline())
		authorized.DELETE("/pipeline/:id", api.DeletePipeline())
		authorized.POST("/pipeline", api.PostPipeline())
		authorized.PATCH("/pipeline", api.PatchPipeline())
		authorized.PUT("/pipelineStatus", api.PutPipelineStatus())

		authorized.GET("/task/:pid/:id", api.GetTask())
		authorized.DELETE("/task/:pid/:id", api.DeleteTask())
		authorized.POST("/task", api.PostTask())
		authorized.PATCH("/task", api.PatchTask())
		authorized.PUT("/taskStatus", api.PutTaskStatus())

		g.POST("/authenticate", api.Authenticate())
		authorized.POST("/user", api.PostUser())
		authorized.GET("/user", api.GetUser())
		authorized.DELETE("/user/:id", api.DeleteUser())
	}

	g.GET("/healthCheck", HealthCheckHandler())

	g.Run(fmt.Sprintf(":%d", cfg.ServerPort))
}

func HealthCheckHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {

	}
}
