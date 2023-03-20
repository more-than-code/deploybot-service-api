package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Pipeline struct {
	Id            primitive.ObjectID `json:"id" bson:"_id"`
	Name          string             `json:"name"`
	CreatedAt     primitive.DateTime `json:"createAt"`
	UpdatedAt     primitive.DateTime `json:"updateAt"`
	ExecutedAt    primitive.DateTime `json:"executedAt"`
	StoppedAt     primitive.DateTime `json:"stoppedAt"`
	ScheduledAt   primitive.DateTime `json:"scheduledAt"`
	Status        string             `json:"status"`
	Arguments     []string           `json:"arguments"`
	Tasks         []Task             `json:"tasks"`
	RepoWatched   string             `json:"repoWatched"`
	BranchWatched string             `json:"branchWatched"`
	AutoRun       bool               `json:"autoRun"`
	ProjectId     primitive.ObjectID `json:"projectId"`
}

type CreatePipelineInput struct {
	Name          string
	Arguments     []string
	RepoWatched   string
	BranchWatched string
	AutoRun       bool
	ProjectId     primitive.ObjectID
}

type TaskFilter struct {
	UpstreamTaskId *primitive.ObjectID
	AutoRun        *bool
}

type GetPipelineInput struct {
	Id         primitive.ObjectID
	Name       string
	TaskFilter TaskFilter
}

type GetPipelinesInput struct {
	RepoWatched   *string `bson:",omitempty"`
	BranchWatched *string `bson:",omitempty"`
	AutoRun       *bool   `bson:",omitempty"`
	ProjectId     primitive.ObjectID
}

type GetPipelinesOutput struct {
	TotalCount int        `json:"totalCount"`
	Items      []Pipeline `json:"items"`
}

type PipelineUpdate struct {
	Name          *string             `bson:",omitempty"`
	ScheduledAt   *primitive.DateTime `bson:",omitempty"`
	Arguments     []string            `bson:",omitempty"`
	RepoWatched   *string             `bson:",omitempty"`
	BranchWatched *string             `bson:",omitempty"`
	AutoRun       *bool               `bson:",omitempty"`
}

type UpdatePipelineInput struct {
	Id       primitive.ObjectID
	Pipeline PipelineUpdate
}

type UpdatePipelineStatusInput struct {
	PipelineId primitive.ObjectID
	Pipeline   struct {
		Status string
	}
}
