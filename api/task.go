package api

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	types "github.com/more-than-code/deploybot-service-api/deploybot-types"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (a *Api) PostTask() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var input types.CreateTaskInput
		err := ctx.BindJSON(&input)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, types.PostTaskResponse{Code: types.CodeClientError, Msg: err.Error()})
			return
		}

		id, err := a.repo.CreateTask(ctx, &input)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, types.PostTaskResponse{Code: types.CodeServerError, Msg: err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, types.PostTaskResponse{Payload: &types.PostTaskResponsePayload{Id: id}})
	}
}

func (a *Api) GetTask() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		pidStr := ctx.Query("pid")
		idStr := ctx.Query("id")

		pid, _ := primitive.ObjectIDFromHex(pidStr)
		id, _ := primitive.ObjectIDFromHex(idStr)

		input := types.GetTaskInput{PipelineId: pid, Id: id}

		task, err := a.repo.GetTask(ctx, &input)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, types.GetTaskResponse{Code: types.CodeServerError, Msg: err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, types.GetTaskResponse{Payload: &types.GetTaskResponsePayload{Task: *task}})
	}
}

func (a *Api) DeleteTask() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var input types.DeleteTaskInput
		err := ctx.BindJSON(&input)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, types.DeleteTaskResponse{Code: types.CodeServerError, Msg: err.Error()})
			return
		}

		err = a.repo.DeleteTask(ctx, &input)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, types.DeleteTaskResponse{Code: types.CodeServerError, Msg: err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, types.DeleteTaskResponse{})
	}
}

func (a *Api) PatchTask() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var input types.UpdateTaskInput
		err := ctx.BindJSON(&input)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, types.PatchTaskResponse{Code: types.CodeClientError, Msg: err.Error()})
			return
		}

		err = a.repo.UpdateTask(ctx, input)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, types.PatchTaskResponse{Code: types.CodeServerError, Msg: err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, types.PatchTaskResponse{})
	}
}

func (a *Api) PutTaskStatus() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var input types.UpdateTaskStatusInput
		err := ctx.BindJSON(&input)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, types.PutTaskStatusResponse{Code: types.CodeClientError, Msg: err.Error()})
			return
		}

		err = a.repo.UpdateTaskStatus(ctx, &input)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, types.PutTaskStatusResponse{Code: types.CodeServerError, Msg: err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, types.PutTaskStatusResponse{})

		go func() {
			pStatus := types.PipelineIdle

			if input.Task.Status == types.TaskDone {
				autoRun := true
				pl, _ := a.repo.GetPipeline(ctx, types.GetPipelineInput{Id: input.PipelineId, TaskFilter: types.TaskFilter{UpstreamTaskId: &input.TaskId, AutoRun: &autoRun}})

				if pl == nil || len(pl.Tasks) == 0 {
					a.repo.UpdatePipelineStatus(ctx, types.UpdatePipelineStatusInput{PipelineId: input.PipelineId, Pipeline: struct{ Status string }{Status: types.PipelineIdle}})
					return
				}

				for _, t := range pl.Tasks {
					body, _ := json.Marshal(types.StreamWebhook{Payload: types.StreamWebhookPayload{PipelineId: pl.Id, TaskId: t.Id, Arguments: pl.Arguments}})

					req, _ := http.NewRequest("POST", t.StreamWebhook, bytes.NewReader(body))
					res, _ := http.DefaultClient.Do(req)

					if res != nil {
						log.Println(res.Status)
					}
				}
			} else if input.Task.Status == types.TaskInProgress {
				pStatus = types.PipelineBusy
			}

			a.repo.UpdatePipelineStatus(ctx, types.UpdatePipelineStatusInput{PipelineId: input.PipelineId, Pipeline: struct{ Status string }{Status: pStatus}})
		}()
	}
}
