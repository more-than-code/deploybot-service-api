package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/kelseyhightower/envconfig"
	"github.com/more-than-code/deploybot-service-api/api"
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

	api := api.NewApi()

	g.GET("/pipelines", api.GetPipelines())
	g.GET("/pipeline/:name", api.GetPipeline())
	g.POST("/pipeline", api.PostPipeline())
	g.PATCH("/pipeline", api.PatchPipeline())
	g.PUT("/pipelineStatus", api.PutPipelineStatus())

	g.GET("/task/:pid/:tid", api.GetTask())
	g.POST("/task", api.PostTask())
	g.PATCH("/task", api.PatchTask())
	g.PUT("/taskStatus", api.PutTaskStatus())

	g.Run(fmt.Sprintf(":%d", cfg.ServerPort))
}
