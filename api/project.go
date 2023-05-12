package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	types "github.com/more-than-code/deploybot-service-api/deploybot-types"
	"github.com/more-than-code/deploybot-service-api/util"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (a *Api) PostProject() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var input types.CreateProjectInput
		err := ctx.BindJSON(&input)

		input.UserId = util.GetUserFromContext(ctx).Id

		if err != nil {
			ctx.JSON(http.StatusBadRequest, types.PostProjectResponse{Code: types.CodeClientError, Msg: err.Error()})
			return
		}

		_, err = a.repo.CreateProject(ctx, &input)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, types.PostProjectResponse{Code: types.CodeServerError, Msg: err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, types.PostProjectResponse{})
	}

}

func (a *Api) GetProjects() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		output, err := a.repo.GetProjects(ctx, types.GetProjectsInput{UserId: util.GetUserFromContext(ctx).Id})

		if err != nil {
			ctx.JSON(http.StatusBadRequest, types.GetProjectsResponse{Code: types.CodeClientError, Msg: err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, types.GetProjectsResponse{Payload: output})
	}
}

func (a *Api) GetProject() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		idStr := ctx.Param("id")
		id, _ := primitive.ObjectIDFromHex(idStr)

		output, err := a.repo.GetProject(ctx, types.GetProjectInput{Id: id, UserId: util.GetUserFromContext(ctx).Id})

		if err != nil {
			ctx.JSON(http.StatusBadRequest, types.GetProjectsResponse{Code: types.CodeClientError, Msg: err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, types.GetProjectResponse{Payload: output})
	}
}

func (a *Api) DeleteProject() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("pid")
		objId, _ := primitive.ObjectIDFromHex(id)

		err := a.repo.DeleteProject(ctx, types.DeleteProjectInput{Id: objId, UserId: util.GetUserFromContext(ctx).Id})

		if err != nil {
			ctx.JSON(http.StatusBadRequest, types.DeleteProjectResponse{Code: types.CodeClientError, Msg: err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, types.DeleteProjectResponse{})
	}
}

func (a *Api) PatchProject() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var project types.UpdateProject
		err := ctx.BindJSON(&project)

		idStr := ctx.Param("id")
		id, _ := primitive.ObjectIDFromHex(idStr)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, types.PatchProjectResponse{Code: types.CodeClientError, Msg: err.Error()})
			return
		}

		err = a.repo.UpdateProject(ctx, types.UpdateProjectInput{Id: id, UserId: util.GetUserFromContext(ctx).Id, Project: project})

		if err != nil {
			ctx.JSON(http.StatusBadRequest, types.PatchProjectResponse{Code: types.CodeServerError, Msg: err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, types.PatchProjectResponse{})
	}
}
