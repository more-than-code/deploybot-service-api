package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	types "github.com/more-than-code/deploybot-service-api/deploybot-types"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (a *Api) PostPipeline() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var input types.CreatePipelineInput
		err := ctx.BindJSON(&input)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, types.PostPipelineResponse{Code: types.CodeClientError, Msg: err.Error()})
			return
		}

		id, err := a.repo.CreatePipeline(ctx, &input)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, types.PostPipelineResponse{Code: types.CodeServerError, Msg: err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, types.PostPipelineResponse{Payload: &types.PostPipelineResponsePayload{id}})
	}

}

func (a *Api) GetPipelines() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		pidStr := ctx.Query("pid")
		pid, _ := primitive.ObjectIDFromHex(pidStr)

		repoWatched, exists := ctx.GetQuery("repoWatched")

		var rw *string
		if exists {
			rw = &repoWatched
		}

		branchWatched, exists := ctx.GetQuery("branchWatched")

		var bw *string
		if exists && branchWatched != "" {
			bw = &branchWatched
		}

		autoRun, exists := ctx.GetQuery("autoRun")
		var ar *bool
		if exists {
			cVal := false
			if autoRun == "true" {
				cVal = true
			}

			ar = &cVal
		}

		output, err := a.repo.GetPipelines(ctx, types.GetPipelinesInput{RepoWatched: rw, BranchWatched: bw, AutoRun: ar, ProjectId: pid})

		if err != nil {
			ctx.JSON(http.StatusBadRequest, types.GetPipelinesResponse{Code: types.CodeClientError, Msg: err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, types.GetPipelinesResponse{Payload: output})
	}
}

func (a *Api) GetPipeline() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		name := ctx.Query("name")
		idStr := ctx.Query("id")
		id, _ := primitive.ObjectIDFromHex(idStr)

		input := types.GetPipelineInput{Name: name, Id: id}
		pl, err := a.repo.GetPipeline(ctx, input)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, types.GetPipelineResponse{Code: types.CodeClientError, Msg: err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, types.GetPipelineResponse{Payload: &types.GetPipelineResponsePayload{*pl}})
	}
}

func (a *Api) DeletePipeline() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")
		objId, _ := primitive.ObjectIDFromHex(id)

		err := a.repo.DeletePipeline(ctx, objId)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, types.DeletePipelineResponse{Code: types.CodeClientError, Msg: err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, types.DeletePipelineResponse{})
	}
}

func (a *Api) PatchPipeline() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var input types.UpdatePipelineInput
		err := ctx.BindJSON(&input)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, types.PatchPipelineResponse{Code: types.CodeClientError, Msg: err.Error()})
			return
		}

		err = a.repo.UpdatePipeline(ctx, input)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, types.PatchPipelineResponse{Code: types.CodeServerError, Msg: err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, types.PatchPipelineResponse{})
	}
}

func (a *Api) PutPipelineStatus() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var input types.UpdatePipelineStatusInput
		err := ctx.BindJSON(&input)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, types.PutPipelineStatusResponse{Code: types.CodeClientError, Msg: err.Error()})
			return
		}

		err = a.repo.UpdatePipelineStatus(ctx, input)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, types.PutPipelineStatusResponse{Code: types.CodeServerError, Msg: err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, types.PutPipelineStatusResponse{})
	}
}
