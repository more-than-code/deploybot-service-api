package repository

import (
	"context"
	"time"

	types "github.com/more-than-code/deploybot-service-api/deploybot-types"
	"github.com/more-than-code/deploybot-service-api/util"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"

	"go.mongodb.org/mongo-driver/bson"
)

func (r *Repository) CreatePipeline(ctx context.Context, input *types.CreatePipelineInput) (primitive.ObjectID, error) {
	doc := util.StructToBsonDoc(input)

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

func (r *Repository) GetPipelines(ctx context.Context, input types.GetPipelinesInput) (*types.GetPipelinesOutput, error) {
	coll := r.mongoClient.Database("pipeline").Collection("pipelines")

	filter := util.StructToBsonDoc(input)

	opts := options.Find().SetSort(bson.D{{"executedat", -1}})
	cursor, err := coll.Find(ctx, filter, opts)

	if err != nil {
		return nil, err
	}

	var output types.GetPipelinesOutput
	if err = cursor.All(ctx, &output.Items); err != nil {
		return nil, err
	}

	output.TotalCount = len(output.Items)

	return &output, nil
}

func (r *Repository) GetPipeline(ctx context.Context, input types.GetPipelineInput) (*types.Pipeline, error) {
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

	var pipeline types.Pipeline
	err := coll.FindOne(ctx, filter, &opts).Decode(&pipeline)

	if err != nil {
		return nil, err
	}

	return &pipeline, nil
}

func (r *Repository) UpdatePipeline(ctx context.Context, input types.UpdatePipelineInput) error {
	filter := bson.M{"_id": input.Id}

	doc := util.StructToBsonDoc(input.Pipeline)
	doc["updatedat"] = primitive.NewDateTimeFromTime(time.Now().UTC())

	update := bson.M{"$set": doc}

	coll := r.mongoClient.Database("pipeline").Collection("pipelines")
	_, err := coll.UpdateOne(ctx, filter, update)

	return err
}

func (r *Repository) UpdatePipelineStatus(ctx context.Context, input types.UpdatePipelineStatusInput) error {
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
