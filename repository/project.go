package repository

import (
	"context"
	"log"
	"time"

	types "github.com/more-than-code/deploybot-service-api/deploybot-types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Project struct {
	Id          primitive.ObjectID `json:"id" bson:"_id"`
	Name        string             `json:"name"`
	AvatarUrl   string             `json:"avatarUrl"`
	OwnerUserId primitive.ObjectID `json:"ownerUserId"`
	Members     []Member           `json:"members"`
	// Pipelines   []Pipeline         `json:"pipelines"`
	CreatedAt primitive.DateTime `json:"createdAt"`
	UpdatedAt primitive.DateTime `json:"updatedAt"`
}

type CreateProjectInput struct {
	Name   string
	UserId primitive.ObjectID
}

type UpdateProject struct {
	Name      *string `bson:",omitempty"`
	AvatarUrl *string `bson:",omitempty"`
}

type UpdateProjectInput struct {
	Id      primitive.ObjectID
	UserId  primitive.ObjectID
	Project UpdateProject
}

type DeleteProjectInput struct {
	Id     primitive.ObjectID
	UserId primitive.ObjectID
}

type GetProjectsInput struct {
	UserId primitive.ObjectID
}

type GetProjectInput struct {
	UserId primitive.ObjectID
	Id     primitive.ObjectID
}

type GetProjectsOutput struct {
	Items      []Project `json:"items"`
	TotalCount int64     `json:"totalCount"`
}

func (r *Repository) CreateProject(ctx context.Context, input *CreateProjectInput) (primitive.ObjectID, error) {
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

func (r *Repository) UpdateProject(ctx context.Context, input UpdateProjectInput) error {
	coll := r.mongoClient.Database("pipeline").Collection("projects")

	project := StructToBsonDoc(input.Project)

	update := bson.M{"$set": project}
	filter := bson.M{"_id": input.Id, "owneruserid": input.UserId}
	_, err := coll.UpdateOne(ctx, filter, update)

	if err != nil {
		log.Println(err)
		return nil
	}

	return nil
}

func (r *Repository) DeleteProject(ctx context.Context, input DeleteProjectInput) error {
	coll := r.mongoClient.Database("pipeline").Collection("projects")
	filter := bson.M{"_id": input.Id, "owneruserid": input.UserId}

	_, err := coll.DeleteOne(ctx, filter)

	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetProject(ctx context.Context, input GetProjectInput) (*Project, error) {
	filter := bson.M{"members.userid": bson.M{"$in": bson.A{input.UserId}}, "_id": input.Id}

	var project Project

	coll := r.mongoClient.Database("pipeline").Collection("projects")
	err := coll.FindOne(ctx, filter).Decode(&project)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &project, nil
}

func (r *Repository) GetProjects(ctx context.Context, input GetProjectsInput) (*GetProjectsOutput, error) {
	filter := bson.M{"members.userid": bson.M{"$in": bson.A{input.UserId}}}

	coll := r.mongoClient.Database("pipeline").Collection("projects")

	cursor, err := coll.Find(ctx, filter)

	if err != nil {
		return nil, err
	}

	var output GetProjectsOutput
	err = cursor.All(ctx, &output.Items)

	if err != nil {
		return nil, err
	}

	count, _ := coll.CountDocuments(ctx, filter)

	output.TotalCount = count

	return &output, nil
}
