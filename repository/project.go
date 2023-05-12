package repository

import (
	"context"
	"log"
	"time"

	types "github.com/more-than-code/deploybot-service-api/deploybot-types"
	"github.com/more-than-code/deploybot-service-api/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (r *Repository) CreateProject(ctx context.Context, input *types.CreateProjectInput) (primitive.ObjectID, error) {
	doc := bson.M{"owneruserid": input.UserId, "createdat": primitive.NewDateTimeFromTime(time.Now().UTC()),
		"name":    input.Name,
		"members": bson.A{bson.M{"userid": input.UserId, "role": types.RoleOwner, "createdat": primitive.NewDateTimeFromTime(time.Now().UTC())}}}

	coll := r.mongoClient.Database("pipeline").Collection("projects")

	res, err := coll.InsertOne(ctx, doc)

	if err != nil {
		log.Println(err)
		return primitive.NilObjectID, err
	}

	return res.InsertedID.(primitive.ObjectID), nil
}

func (r *Repository) UpdateProject(ctx context.Context, input types.UpdateProjectInput) error {
	coll := r.mongoClient.Database("pipeline").Collection("projects")

	project := util.StructToBsonDoc(input.Project)

	update := bson.M{"$set": project}
	filter := bson.M{"_id": input.Id, "owneruserid": input.UserId}
	_, err := coll.UpdateOne(ctx, filter, update)

	if err != nil {
		log.Println(err)
		return nil
	}

	return nil
}

func (r *Repository) DeleteProject(ctx context.Context, input types.DeleteProjectInput) error {
	coll := r.mongoClient.Database("pipeline").Collection("projects")
	filter := bson.M{"_id": input.Id, "owneruserid": input.UserId}

	_, err := coll.DeleteOne(ctx, filter)

	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetProject(ctx context.Context, input types.GetProjectInput) (*types.Project, error) {
	filter := bson.M{"members.userid": bson.M{"$in": bson.A{input.UserId}}, "_id": input.Id}

	var project types.Project

	coll := r.mongoClient.Database("pipeline").Collection("projects")
	err := coll.FindOne(ctx, filter).Decode(&project)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &project, nil
}

func (r *Repository) GetProjects(ctx context.Context, input types.GetProjectsInput) (*types.GetProjectsOutput, error) {
	filter := bson.M{"members.userid": bson.M{"$in": bson.A{input.UserId}}}

	coll := r.mongoClient.Database("pipeline").Collection("projects")

	cursor, err := coll.Find(ctx, filter)

	if err != nil {
		return nil, err
	}

	var output types.GetProjectsOutput
	err = cursor.All(ctx, &output.Items)

	if err != nil {
		return nil, err
	}

	count, _ := coll.CountDocuments(ctx, filter)

	output.TotalCount = count

	return &output, nil
}
