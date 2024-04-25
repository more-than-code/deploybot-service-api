package repository

import (
	"context"
	"time"

	types "github.com/more-than-code/deploybot-service-api/deploybot-types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"

	"go.mongodb.org/mongo-driver/bson"
)

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
	Labels        map[string]*string `json:"labels"`
	Tasks         []types.Task       `json:"tasks"`
	RepoWatched   string             `json:"repoWatched"`
	BranchWatched string             `json:"branchWatched"`
	AutoRun       bool               `json:"autoRun"`
	ProjectId     primitive.ObjectID `json:"projectId"`
}

type CreatePipelineInput struct {
	Name          string
	Arguments     []string
	Labels        map[string]*string
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
	Labels        map[string]*string  `bson:",omitempty"`
	RepoWatched   *string             `bson:",omitempty"`
	BranchWatched *string             `bson:",omitempty"`
	AutoRun       *bool               `bson:",omitempty"`
	ProjectId     *primitive.ObjectID `bson:",omitempty"`
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

func (r *Repository) CreatePipeline(ctx context.Context, input *CreatePipelineInput) (primitive.ObjectID, error) {
	doc := StructToBsonDoc(input)

	doc["createdat"] = primitive.NewDateTimeFromTime(time.Now().UTC())
	doc["status"] = types.PipelineIdle

	coll := r.mongoClient.Database("pipeline").Collection("pipelines")
	result, err := coll.InsertOne(ctx, doc)

	if err != nil {
		return primitive.NilObjectID, err
	}

	return result.InsertedID.(primitive.ObjectID), nil
}

func (r *Repository) DeletePipeline(ctx context.Context, id primitive.ObjectID) error {
	coll := r.mongoClient.Database("pipeline").Collection("pipelines")
	_, err := coll.DeleteOne(ctx, bson.M{"_id": id})

	return err
}

func (r *Repository) GetPipelines(ctx context.Context, input GetPipelinesInput) (*GetPipelinesOutput, error) {
	coll := r.mongoClient.Database("pipeline").Collection("pipelines")

	filter := StructToBsonDoc(input)

	opts := options.Find().SetSort(bson.D{{"executedat", -1}})
	cursor, err := coll.Find(ctx, filter, opts)

	if err != nil {
		return nil, err
	}

	var output GetPipelinesOutput
	if err = cursor.All(ctx, &output.Items); err != nil {
		return nil, err
	}

	output.TotalCount = len(output.Items)

	return &output, nil
}

func (r *Repository) GetPipeline(ctx context.Context, input GetPipelineInput) (*Pipeline, error) {
	coll := r.mongoClient.Database("pipeline").Collection("pipelines")

	filter := bson.M{}
	if !input.Id.IsZero() {
		filter["_id"] = input.Id
	}
	if input.Name != "" {
		filter["name"] = input.Name
	}

	taskFilter := bson.M{}
	if input.TaskFilter.UpstreamTaskId != nil {
		taskFilter["upstreamtaskid"] = input.TaskFilter.UpstreamTaskId

	}
	if input.TaskFilter.AutoRun != nil {
		taskFilter["autorun"] = input.TaskFilter.AutoRun
	}

	opts := options.FindOneOptions{}
	if len(taskFilter) > 0 {
		opts.SetProjection(bson.M{"arguments": 1, "tasks": bson.M{"$elemMatch": taskFilter}})
	}

	var pipeline Pipeline
	err := coll.FindOne(ctx, filter, &opts).Decode(&pipeline)

	if err != nil {
		return nil, err
	}

	return &pipeline, nil
}

func (r *Repository) UpdatePipeline(ctx context.Context, input UpdatePipelineInput) error {
	filter := bson.M{"_id": input.Id}

	doc := StructToBsonDoc(input.Pipeline)
	doc["updatedat"] = primitive.NewDateTimeFromTime(time.Now().UTC())

	update := bson.M{"$set": doc}

	coll := r.mongoClient.Database("pipeline").Collection("pipelines")
	_, err := coll.UpdateOne(ctx, filter, update)

	return err
}

func (r *Repository) UpdatePipelineStatus(ctx context.Context, input UpdatePipelineStatusInput) error {
	filter := bson.M{"_id": input.PipelineId}

	doc := bson.M{"status": input.Pipeline.Status}

	switch input.Pipeline.Status {
	case types.PipelineBusy:
		doc["executedat"] = primitive.NewDateTimeFromTime(time.Now().UTC())
		doc["stoppedat"] = nil
	case types.PipelineIdle:
		doc["stoppedat"] = primitive.NewDateTimeFromTime(time.Now().UTC())
	}

	update := bson.M{"$set": doc}

	coll := r.mongoClient.Database("pipeline").Collection("pipelines")
	_, err := coll.UpdateOne(ctx, filter, update)

	return err
}
