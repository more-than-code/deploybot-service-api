package api

import (
	types "github.com/more-than-code/deploybot-service-api/deploybot-types"
	"github.com/more-than-code/deploybot-service-api/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PostPipelineResponsePayload struct {
	Id primitive.ObjectID
}
type PostPipelineResponse struct {
	Code    int                          `json:"code"`
	Msg     string                       `json:"msg"`
	Payload *PostPipelineResponsePayload `json:"payload"`
}

type GetPipelinesResponse struct {
	Code    int                            `json:"code"`
	Msg     string                         `json:"msg"`
	Payload *repository.GetPipelinesOutput `json:"payload"`
}

type GetPipelineResponsePayload struct {
	Pipeline repository.Pipeline
}

type GetPipelineResponse struct {
	Code    int                         `json:"code"`
	Msg     string                      `json:"msg"`
	Payload *GetPipelineResponsePayload `json:"payload"`
}

type DeletePipelineResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

type PatchPipelineResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

type PutPipelineStatusResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

type PostTaskResponsePayload struct {
	Id primitive.ObjectID
}
type PostTaskResponse struct {
	Code    int                      `json:"code"`
	Msg     string                   `json:"msg"`
	Payload *PostTaskResponsePayload `json:"payload"`
}

type GetTaskResponsePayload struct {
	Task types.Task `json:"task"`
}
type GetTaskResponse struct {
	Code    int                     `json:"code"`
	Msg     string                  `json:"msg"`
	Payload *GetTaskResponsePayload `json:"payload"`
}

type DeleteTaskResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

type PatchTaskResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

type PutTaskStatusResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

type AuthenticationResponse struct {
	Msg     string                           `json:"msg"`
	Code    int                              `json:"code"`
	Payload *repository.AuthenticationOutput `json:"payload"`
}

type PostUserResponse struct {
	Msg  string `json:"msg"`
	Code int    `json:"code"`
}

type DeleteUserResponse struct {
	Msg  string `json:"msg"`
	Code int    `json:"code"`
}

type GetUserResponse struct {
	Msg     string           `json:"msg"`
	Code    int              `json:"code"`
	Payload *repository.User `json:"payload"`
}

type GetUsersResponse struct {
	Msg     string                     `json:"msg"`
	Code    int                        `json:"code"`
	Payload *repository.GetUsersOutput `json:"payload"`
}

type GetProjectsResponse struct {
	Msg     string                        `json:"msg"`
	Code    int                           `json:"code"`
	Payload *repository.GetProjectsOutput `json:"payload"`
}

type GetProjectResponse struct {
	Msg     string              `json:"msg"`
	Code    int                 `json:"code"`
	Payload *repository.Project `json:"payload"`
}

type DeleteProjectResponse struct {
	Msg  string `json:"msg"`
	Code int    `json:"code"`
}

type PostProjectResponse struct {
	Msg  string `json:"msg"`
	Code int    `json:"code"`
}

type PatchProjectResponse struct {
	Msg  string `json:"msg"`
	Code int    `json:"code"`
}

type DeleteMemberResponse struct {
	Msg  string `json:"msg"`
	Code int    `json:"code"`
}

type PostMemberResponse struct {
	Msg  string `json:"msg"`
	Code int    `json:"code"`
}

type PatchMemberResponse struct {
	Msg  string `json:"msg"`
	Code int    `json:"code"`
}

type GetDiskInfoResponse struct {
	Msg     string          `json:"msg"`
	Code    int             `json:"code"`
	Payload *types.DiskInfo `json:"payload"`
}
