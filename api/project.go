package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	types "github.com/more-than-code/deploybot-service-api/deploybot-types"
	"github.com/more-than-code/deploybot-service-api/repository"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (a *Api) PostProject() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var input repository.CreateProjectInput
		err := ctx.BindJSON(&input)

		input.UserId = repository.GetUserFromContext(ctx).Id

		if err != nil {
			ctx.JSON(http.StatusBadRequest, PostProjectResponse{Code: types.CodeClientError, Msg: err.Error()})
			return
		}

		_, err = a.repo.CreateProject(ctx, &input)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, PostProjectResponse{Code: types.CodeServerError, Msg: err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, PostProjectResponse{})
	}

}

func (a *Api) GetProjects() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		output, err := a.repo.GetProjects(ctx, repository.GetProjectsInput{UserId: repository.GetUserFromContext(ctx).Id})

		if err != nil {
			ctx.JSON(http.StatusBadRequest, GetProjectsResponse{Code: types.CodeClientError, Msg: err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, GetProjectsResponse{Payload: output})
	}
}

func (a *Api) GetProject() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		idStr := ctx.Param("id")
		id, _ := primitive.ObjectIDFromHex(idStr)

		output, err := a.repo.GetProject(ctx, repository.GetProjectInput{Id: id, UserId: repository.GetUserFromContext(ctx).Id})

		if err != nil {
			ctx.JSON(http.StatusBadRequest, GetProjectsResponse{Code: types.CodeClientError, Msg: err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, GetProjectResponse{Payload: output})
	}
}

func (a *Api) DeleteProject() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("pid")
		objId, _ := primitive.ObjectIDFromHex(id)

		err := a.repo.DeleteProject(ctx, repository.DeleteProjectInput{Id: objId, UserId: repository.GetUserFromContext(ctx).Id})

		if err != nil {
			ctx.JSON(http.StatusBadRequest, DeleteProjectResponse{Code: types.CodeClientError, Msg: err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, DeleteProjectResponse{})
	}
}

func (a *Api) PatchProject() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var project repository.UpdateProject
		err := ctx.BindJSON(&project)

		idStr := ctx.Param("id")
		id, _ := primitive.ObjectIDFromHex(idStr)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, PatchProjectResponse{Code: types.CodeClientError, Msg: err.Error()})
			return
		}

		err = a.repo.UpdateProject(ctx, repository.UpdateProjectInput{Id: id, UserId: repository.GetUserFromContext(ctx).Id, Project: project})

		if err != nil {
			ctx.JSON(http.StatusBadRequest, PatchProjectResponse{Code: types.CodeServerError, Msg: err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, PatchProjectResponse{})
	}
}
