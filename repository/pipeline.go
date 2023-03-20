package repository

import (
	"context"
	"time"

	"github.com/more-than-code/deploybot-service-api/model"
	"github.com/more-than-code/deploybot-service-api/util"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"

	"go.mongodb.org/mongo-driver/bson"
)

func (r *Repository) CreatePipeline(ctx context.Context, input *model.CreatePipelineInput) (primitive.ObjectID, error) {
	doc := util.StructToBsonDoc(input)

	doc["createdat"] = primitive.NewDateTimeFromTime(time.Now().UTC())
	doc["status"] = model.PipelineIdle

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

func (r *Repository) GetPipelines(ctx context.Context, input model.GetPipelinesInput) (*model.GetPipelinesOutput, error) {
	coll := r.mongoClient.Database("pipeline").Collection("pipelines")

	filter := util.StructToBsonDoc(input)

	opts := options.Find().SetSort(bson.D{{"executedat", -1}})
	cursor, err := coll.Find(ctx, filter, opts)

	if err != nil {
		return nil, err
	}

	var output model.GetPipelinesOutput
	if err = cursor.All(ctx, &output.Items); err != nil {
		return nil, err
	}

	output.TotalCount = len(output.Items)

	return &output, nil
}

func (r *Repository) GetPipeline(ctx context.Context, input model.GetPipelineInput) (*model.Pipeline, error) {
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

	var pipeline model.Pipeline
	err := coll.FindOne(ctx, filter, &opts).Decode(&pipeline)

	if err != nil {
		return nil, err
	}

	return &pipeline, nil
}

func (r *Repository) UpdatePipeline(ctx context.Context, input model.UpdatePipelineInput) error {
	filter := bson.M{"_id": input.Id}

	doc := bson.M{}
	doc["updatedat"] = primitive.NewDateTimeFromTime(time.Now().UTC())

	if input.Pipeline.Name != nil {
		doc["name"] = input.Pipeline.Name
	}
	if input.Pipeline.ScheduledAt != nil {
		doc["scheduledat"] = input.Pipeline.ScheduledAt
	}
	if input.Pipeline.AutoRun != nil {
		doc["autorun"] = input.Pipeline.AutoRun
	}
	if input.Pipeline.Arguments != nil {
		doc["arguments"] = input.Pipeline.Arguments
	}
	if input.Pipeline.RepoWatched != nil {
		doc["repowatched"] = input.Pipeline.RepoWatched
	}
	if input.Pipeline.BranchWatched != nil {
		doc["branchwatched"] = input.Pipeline.BranchWatched
	}

	update := bson.M{"$set": doc}

	coll := r.mongoClient.Database("pipeline").Collection("pipelines")
	_, err := coll.UpdateOne(ctx, filter, update)

	return err
}

func (r *Repository) UpdatePipelineStatus(ctx context.Context, input model.UpdatePipelineStatusInput) error {
	filter := bson.M{"_id": input.PipelineId}

	doc := bson.M{"status": input.Pipeline.Status}

	switch input.Pipeline.Status {
	case model.PipelineBusy:
		doc["executedat"] = primitive.NewDateTimeFromTime(time.Now().UTC())
		doc["stoppedat"] = nil
	case model.PipelineIdle:
		doc["stoppedat"] = primitive.NewDateTimeFromTime(time.Now().UTC())
	}

	update := bson.M{"$set": doc}

	coll := r.mongoClient.Database("pipeline").Collection("pipelines")
	_, err := coll.UpdateOne(ctx, filter, update)

	return err
}
