package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/more-than-code/deploybot-service-api/model"
	"github.com/more-than-code/deploybot-service-api/util"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (a *Api) PostProject() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var input model.CreateProjectInput
		err := ctx.BindJSON(&input)

		input.UserId = util.GetUserFromContext(ctx).Id

		if err != nil {
			ctx.JSON(http.StatusBadRequest, PostProjectResponse{Code: CodeClientError, Msg: err.Error()})
			return
		}

		_, err = a.repo.CreateProject(ctx, &input)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, PostProjectResponse{Code: CodeServerError, Msg: err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, PostProjectResponse{})
	}

}

func (a *Api) GetProjects() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		output, err := a.repo.GetProjects(ctx, model.GetProjectsInput{UserId: util.GetUserFromContext(ctx).Id})

		if err != nil {
			ctx.JSON(http.StatusBadRequest, GetProjectsResponse{Code: CodeClientError, Msg: err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, GetProjectsResponse{Payload: output})
	}
}

func (a *Api) DeleteProject() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("pid")
		objId, _ := primitive.ObjectIDFromHex(id)

		err := a.repo.DeleteProject(ctx, model.DeleteProjectInput{Id: objId, UserId: util.GetUserFromContext(ctx).Id})

		if err != nil {
			ctx.JSON(http.StatusBadRequest, DeleteProjectResponse{Code: CodeClientError, Msg: err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, DeleteProjectResponse{})
	}
}

func (a *Api) PatchProject() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var project model.UpdateProject
		err := ctx.BindJSON(&project)

		idStr := ctx.Param("id")
		id, _ := primitive.ObjectIDFromHex(idStr)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, PatchProjectResponse{Code: CodeClientError, Msg: err.Error()})
			return
		}

		err = a.repo.UpdateProject(ctx, model.UpdateProjectInput{Id: id, UserId: util.GetUserFromContext(ctx).Id, Project: project})

		if err != nil {
			ctx.JSON(http.StatusBadRequest, PatchProjectResponse{Code: CodeServerError, Msg: err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, PatchProjectResponse{})
	}
}
