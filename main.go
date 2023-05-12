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
		authorized.GET("/pipeline", api.GetPipeline())
		authorized.DELETE("/pipeline/:id", api.DeletePipeline())
		authorized.POST("/pipeline", api.PostPipeline())
		authorized.PATCH("/pipeline", api.PatchPipeline())
		authorized.PUT("/pipelineStatus", api.PutPipelineStatus())

		authorized.GET("/task", api.GetTask())
		authorized.DELETE("/task", api.DeleteTask())
		authorized.POST("/task", api.PostTask())
		authorized.PATCH("/task", api.PatchTask())
		authorized.PUT("/taskStatus", api.PutTaskStatus())

		authorized.GET("/projects", api.GetProjects())
		authorized.GET("/project/:id", api.GetProject())
		authorized.DELETE("/project/:id", api.DeleteProject())
		authorized.POST("/project", api.PostProject())
		authorized.PATCH("/project/:id", api.PatchProject())

		authorized.DELETE("/member", api.DeleteMember())
		authorized.POST("/member", api.PostMember())
		authorized.PATCH("/member", api.PatchMember())

		g.POST("/authenticate", api.Authenticate())
		g.POST("/authenticateSso", api.AuthenticateSso())
		authorized.POST("/user", api.PostUser())
		authorized.GET("/user", api.GetUser())
		authorized.GET("/users", api.GetUsers())
		authorized.DELETE("/user/:id", api.DeleteUser())
	}

	g.GET("/healthCheck", HealthCheckHandler())

	g.Run(fmt.Sprintf(":%d", cfg.ServerPort))
}

func HealthCheckHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {

	}
}
