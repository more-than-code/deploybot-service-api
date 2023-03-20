package repository

import (
	"context"
	"time"

	"github.com/more-than-code/deploybot-service-api/model"
	"github.com/more-than-code/deploybot-service-api/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (r *Repository) CreateTask(ctx context.Context, input *model.CreateTaskInput) (primitive.ObjectID, error) {
	coll := r.mongoClient.Database("pipeline").Collection("pipelines")
	filter := bson.M{"_id": input.PipelineId}

	doc := util.StructToBsonDoc(input.Task)
	if input.Task.Id.IsZero() {
		doc["id"] = primitive.NewObjectID()
	}

	doc["status"] = model.TaskPending
	doc["createdat"] = primitive.NewDateTimeFromTime(time.Now().UTC())

	update := bson.M{"$push": bson.M{"tasks": doc}}
	_, err := coll.UpdateOne(ctx, filter, update)

	return doc["id"].(primitive.ObjectID), err
}

func (r *Repository) GetTask(ctx context.Context, input *model.GetTaskInput) (*model.Task, error) {
	coll := r.mongoClient.Database("pipeline").Collection("pipelines")
	filter := bson.M{"_id": input.PipelineId, "tasks.id": input.Id}

	opts := options.FindOneOptions{Projection: bson.M{"tasks.$": 1}}
	var pipeline model.Pipeline
	err := coll.FindOne(ctx, filter, &opts).Decode(&pipeline)

	if err != nil {
		return nil, err
	}

	return &pipeline.Tasks[0], nil
}

func (r *Repository) GetTasks(ctx context.Context, input model.GetTasksInput) ([]model.Task, error) {
	coll := r.mongoClient.Database("pipeline").Collection("pipelines")
	filter := bson.M{"_id": input.PipelineId}

	if input.UpstreamTaskId != nil {
		filter["tasks.upstreamtaskid"] = input.UpstreamTaskId
	}

	var pipeline model.Pipeline
	err := coll.FindOne(ctx, filter).Decode(&pipeline)

	if err != nil {
		return nil, err
	}

	return pipeline.Tasks, nil
}

func (r *Repository) DeleteTask(ctx context.Context, input *model.DeleteTaskInput) error {
	coll := r.mongoClient.Database("pipeline").Collection("pipelines")
	filter := bson.M{"_id": input.PipelineId}
	update := bson.M{"$pull": bson.M{"tasks": bson.M{"id": input.Id}}}
	_, err := coll.UpdateOne(ctx, filter, update)

	return err
}

func (r *Repository) UpdateTask(ctx context.Context, input model.UpdateTaskInput) error {
	filter := bson.M{"_id": input.PipelineId, "tasks.id": input.Id}

	doc := bson.M{}
	doc["tasks.$.updatedat"] = primitive.NewDateTimeFromTime(time.Now().UTC())

	if input.Task.Name != nil {
		doc["tasks.$.name"] = input.Task.Name
	}
	if input.Task.ScheduledAt != nil {
		doc["tasks.$.scheduledat"] = input.Task.ScheduledAt
	}
	if input.Task.Config != nil {
		doc["tasks.$.config"] = input.Task.Config
	}
	if input.Task.Remarks != nil {
		doc["tasks.$.remarks"] = input.Task.Remarks
	}
	if input.Task.AutoRun != nil {
		doc["tasks.$.autorun"] = input.Task.AutoRun
	}
	if input.Task.StreamWebhook != nil {
		doc["tasks.$.streamwebhook"] = input.Task.StreamWebhook
	}
	if input.Task.UpstreamTaskId != nil {
		doc["tasks.$.upstreamtaskid"] = input.Task.UpstreamTaskId
	}
	if input.Task.Timeout != nil {
		doc["tasks.$.timeout"] = input.Task.Timeout
	}
	if input.Task.Type != nil {
		doc["tasks.$.type"] = input.Task.Type
	}

	update := bson.M{"$set": doc}

	coll := r.mongoClient.Database("pipeline").Collection("pipelines")
	_, err := coll.UpdateOne(ctx, filter, update)

	return err
}

func (r *Repository) UpdateTaskStatus(ctx context.Context, input *model.UpdateTaskStatusInput) error {
	filter := bson.M{"_id": input.PipelineId, "tasks.id": input.TaskId}

	doc := bson.M{"tasks.$.status": input.Task.Status}

	switch input.Task.Status {
	case model.TaskInProgress:
		doc["tasks.$.executedat"] = primitive.NewDateTimeFromTime(time.Now().UTC())
		doc["tasks.$.stoppedat"] = nil
	case model.TaskDone, model.TaskFailed, model.TaskCanceled:
		doc["tasks.$.stoppedat"] = primitive.NewDateTimeFromTime(time.Now().UTC())
	}

	update := bson.M{"$set": doc}

	coll := r.mongoClient.Database("pipeline").Collection("pipelines")
	_, err := coll.UpdateOne(ctx, filter, update)

	return err
}
