package repository

import (
	"context"
	"errors"
	"time"

	types "github.com/more-than-code/deploybot-service-api/deploybot-types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Member struct {
	UserId    primitive.ObjectID `json:"userId"`
	Role      types.Role         `json:"role"`
	CreatedAt primitive.DateTime `json:"datetimeCreated"`
}

type CreateMemberInput struct {
	ProjectId primitive.ObjectID
	Member    struct {
		UserId primitive.ObjectID
		Role   types.Role
	}
}

type DeleteMemberInput struct {
	ProjectId primitive.ObjectID
	UserId    primitive.ObjectID
}

type UpdateMember struct {
	Role *types.Role
}
type UpdateMemberInput struct {
	ProjectId primitive.ObjectID
	UserId    primitive.ObjectID
	Member    UpdateMember
}

func (r *Repository) CreateMember(ctx context.Context, input CreateMemberInput) error {
	member := StructToBsonDoc(input.Member)
	member["createdat"] = primitive.NewDateTimeFromTime(time.Now().UTC())

	update := bson.M{"$push": bson.M{"members": member}}

	filter := bson.M{"_id": input.ProjectId, "members.userid": bson.M{"$ne": input.Member.UserId}}

	coll := r.mongoClient.Database("pipeline").Collection("projects")
	_, err := coll.UpdateOne(ctx, filter, update)

	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) DeleteMember(ctx context.Context, input DeleteMemberInput) error {
	update := bson.M{}
	update["$pull"] = bson.M{"members": bson.M{"userid": input.UserId}}

	filter := bson.M{"_id": input.ProjectId, "owneruserid": bson.M{"$ne": input.UserId}, "members.userid": input.UserId}

	coll := r.mongoClient.Database("pipeline").Collection("projects")
	res, err := coll.UpdateOne(ctx, filter, update)

	if err != nil {
		return err
	}

	if res.ModifiedCount == 0 {
		return errors.New("not deleted")
	}

	return nil
}

func (r *Repository) UpdateMember(ctx context.Context, input UpdateMemberInput) error {
	filter := bson.M{"_id": input.ProjectId, "members.userid": input.UserId, "owneruserid": bson.M{"$ne": input.UserId}}

	update := bson.M{"$set": bson.M{"members.$.role": input.Member, "members.$.updatedat": primitive.NewDateTimeFromTime(time.Now().UTC())}}

	coll := r.mongoClient.Database("pipeline").Collection("projects")
	_, err := coll.UpdateOne(ctx, filter, update)

	if err != nil {
		return err
	}

	return nil
}
