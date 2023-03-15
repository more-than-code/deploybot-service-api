package api

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/more-than-code/deploybot-service-api/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (a *Api) PostTask() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var input model.CreateTaskInput
		err := ctx.BindJSON(&input)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, PostTaskResponse{Code: CodeClientError, Msg: err.Error()})
			return
		}

		id, err := a.repo.CreateTask(ctx, &input)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, PostTaskResponse{Code: CodeServerError, Msg: err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, PostTaskResponse{Payload: PostTaskResponsePayload{id}})
	}
}

func (a *Api) GetTask() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		pidStr := ctx.Param("pid")
		idStr := ctx.Param("id")

		pid, _ := primitive.ObjectIDFromHex(pidStr)
		id, _ := primitive.ObjectIDFromHex(idStr)

		input := model.GetTaskInput{PipelineId: pid, Id: id}
		task, err := a.repo.GetTask(ctx, &input)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, GetTaskResponse{Code: CodeServerError, Msg: err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, GetTaskResponse{Payload: GetTaskResponsePayload{Task: task}})
	}
}

func (a *Api) DeleteTask() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		pidStr := ctx.Param("pid")
		idStr := ctx.Param("id")

		pid, _ := primitive.ObjectIDFromHex(pidStr)
		id, _ := primitive.ObjectIDFromHex(idStr)

		input := model.DeleteTaskInput{PipelineId: pid, Id: id}
		err := a.repo.DeleteTask(ctx, &input)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, DeleteTaskResponse{Code: CodeServerError, Msg: err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, DeleteTaskResponse{})
	}
}

func (a *Api) PatchTask() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var input model.UpdateTaskInput
		err := ctx.BindJSON(&input)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, PatchTaskResponse{Code: CodeClientError, Msg: err.Error()})
			return
		}

		err = a.repo.UpdateTask(ctx, input)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, PatchTaskResponse{Code: CodeServerError, Msg: err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, PatchTaskResponse{})
	}
}

func (a *Api) PutTaskStatus() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var input model.UpdateTaskStatusInput
		err := ctx.BindJSON(&input)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, PutTaskStatusResponse{Code: CodeClientError, Msg: err.Error()})
			return
		}

		err = a.repo.UpdateTaskStatus(ctx, &input)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, PutTaskStatusResponse{Code: CodeServerError, Msg: err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, PutTaskStatusResponse{})

		go func() {
			pStatus := model.PipelineIdle

			if input.Payload.Status == model.TaskDone {
				autoRun := true
				pl, _ := a.repo.GetPipeline(ctx, model.GetPipelineInput{Id: &input.PipelineId, TaskFilter: model.TaskFilter{UpstreamTaskId: &input.TaskId, AutoRun: &autoRun}})

				if pl == nil || len(pl.Tasks) == 0 {
					a.repo.UpdatePipelineStatus(ctx, model.UpdatePipelineStatusInput{PipelineId: input.PipelineId, Payload: model.UpdatePipelineStatusInputPayload{Status: model.PipelineIdle}})
					return
				}

				for _, t := range pl.Tasks {
					body, _ := json.Marshal(model.StreamWebhook{Payload: model.StreamWebhookPayload{PipelineId: pl.Id, TaskId: t.Id, Arguments: pl.Arguments}})

					req, _ := http.NewRequest("POST", t.StreamWebhook, bytes.NewReader(body))
					res, _ := http.DefaultClient.Do(req)

					if res != nil {
						log.Println(res.Status)
					}
				}
			} else if input.Payload.Status == model.TaskInProgress {
				pStatus = model.PipelineBusy
			}

			a.repo.UpdatePipelineStatus(ctx, model.UpdatePipelineStatusInput{PipelineId: input.PipelineId, Payload: model.UpdatePipelineStatusInputPayload{Status: pStatus}})
		}()
	}
}
