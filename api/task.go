package api

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	types "github.com/more-than-code/deploybot-service-api/deploybot-types"
	"github.com/more-than-code/deploybot-service-api/repository"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (a *Api) PostTask() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var input repository.CreateTaskInput
		err := ctx.BindJSON(&input)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, PostTaskResponse{Code: types.CodeClientError, Msg: err.Error()})
			return
		}

		id, err := a.repo.CreateTask(ctx, &input)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, PostTaskResponse{Code: types.CodeServerError, Msg: err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, PostTaskResponse{Payload: &PostTaskResponsePayload{Id: id}})
	}
}

func (a *Api) GetTask() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		pidStr := ctx.Query("pid")
		idStr := ctx.Query("id")

		pid, _ := primitive.ObjectIDFromHex(pidStr)
		id, _ := primitive.ObjectIDFromHex(idStr)

		input := repository.GetTaskInput{PipelineId: pid, Id: id}

		task, err := a.repo.GetTask(ctx, &input)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, GetTaskResponse{Code: types.CodeServerError, Msg: err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, GetTaskResponse{Payload: &GetTaskResponsePayload{Task: *task}})
	}
}

func (a *Api) DeleteTask() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var input repository.DeleteTaskInput
		err := ctx.BindJSON(&input)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, DeleteTaskResponse{Code: types.CodeServerError, Msg: err.Error()})
			return
		}

		err = a.repo.DeleteTask(ctx, &input)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, DeleteTaskResponse{Code: types.CodeServerError, Msg: err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, DeleteTaskResponse{})
	}
}

func (a *Api) PatchTask() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var input repository.UpdateTaskInput
		err := ctx.BindJSON(&input)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, PatchTaskResponse{Code: types.CodeClientError, Msg: err.Error()})
			return
		}

		err = a.repo.UpdateTask(ctx, input)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, PatchTaskResponse{Code: types.CodeServerError, Msg: err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, PatchTaskResponse{})
	}
}

func (a *Api) PutTaskStatus() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var input repository.UpdateTaskStatusInput
		err := ctx.BindJSON(&input)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, PutTaskStatusResponse{Code: types.CodeClientError, Msg: err.Error()})
			return
		}

		err = a.repo.UpdateTaskStatus(ctx, &input)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, PutTaskStatusResponse{Code: types.CodeServerError, Msg: err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, PutTaskStatusResponse{})

		go func() {
			pStatus := types.PipelineIdle

			if input.Task.Status == types.TaskDone {
				autoRun := true
				pl, _ := a.repo.GetPipeline(ctx, repository.GetPipelineInput{Id: input.PipelineId, TaskFilter: repository.TaskFilter{UpstreamTaskId: &input.TaskId, AutoRun: &autoRun}})

				if pl == nil || len(pl.Tasks) == 0 {
					a.repo.UpdatePipelineStatus(ctx, repository.UpdatePipelineStatusInput{PipelineId: input.PipelineId, Pipeline: struct{ Status string }{Status: types.PipelineIdle}})
					return
				}

				for _, t := range pl.Tasks {
					body, _ := json.Marshal(types.StreamWebhook{Payload: types.StreamWebhookPayload{PipelineId: types.ObjectId(pl.Id), TaskId: types.ObjectId(t.Id), Arguments: pl.Arguments}})

					req, _ := http.NewRequest("POST", t.StreamWebhook, bytes.NewReader(body))
					res, _ := http.DefaultClient.Do(req)

					if res != nil {
						log.Println(res.Status)
					}
				}
			} else if input.Task.Status == types.TaskInProgress {
				pStatus = types.PipelineBusy
			}

			a.repo.UpdatePipelineStatus(ctx, repository.UpdatePipelineStatusInput{PipelineId: input.PipelineId, Pipeline: struct{ Status string }{Status: pStatus}})
		}()
	}
}
