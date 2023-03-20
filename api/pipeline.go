package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/more-than-code/deploybot-service-api/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (a *Api) PostPipeline() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var input model.CreatePipelineInput
		err := ctx.BindJSON(&input)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, PostPipelineResponse{Code: CodeClientError, Msg: err.Error()})
			return
		}

		id, err := a.repo.CreatePipeline(ctx, &input)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, PostPipelineResponse{Code: CodeServerError, Msg: err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, PostPipelineResponse{Payload: &PostPipelineResponsePayload{id}})
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

		output, err := a.repo.GetPipelines(ctx, model.GetPipelinesInput{RepoWatched: rw, BranchWatched: bw, AutoRun: ar, ProjectId: pid})

		if err != nil {
			ctx.JSON(http.StatusBadRequest, GetPipelinesResponse{Code: CodeClientError, Msg: err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, GetPipelinesResponse{Payload: output})
	}
}

func (a *Api) GetPipeline() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		name := ctx.Query("name")
		idStr := ctx.Query("id")
		id, _ := primitive.ObjectIDFromHex(idStr)

		input := model.GetPipelineInput{Name: name, Id: id}
		pl, err := a.repo.GetPipeline(ctx, input)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, GetPipelineResponse{Code: CodeClientError, Msg: err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, GetPipelineResponse{Payload: &GetPipelineResponsePayload{*pl}})
	}
}

func (a *Api) DeletePipeline() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")
		objId, _ := primitive.ObjectIDFromHex(id)

		err := a.repo.DeletePipeline(ctx, objId)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, DeletePipelineResponse{Code: CodeClientError, Msg: err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, DeletePipelineResponse{})
	}
}

func (a *Api) PatchPipeline() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var input model.UpdatePipelineInput
		err := ctx.BindJSON(&input)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, PatchPipelineResponse{Code: CodeClientError, Msg: err.Error()})
			return
		}

		err = a.repo.UpdatePipeline(ctx, input)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, PatchPipelineResponse{Code: CodeServerError, Msg: err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, PatchPipelineResponse{})
	}
}

func (a *Api) PutPipelineStatus() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var input model.UpdatePipelineStatusInput
		err := ctx.BindJSON(&input)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, PutPipelineStatusResponse{Code: CodeClientError, Msg: err.Error()})
			return
		}

		err = a.repo.UpdatePipelineStatus(ctx, input)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, PutPipelineStatusResponse{Code: CodeServerError, Msg: err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, PutPipelineStatusResponse{})
	}
}
