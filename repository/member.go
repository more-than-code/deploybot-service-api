package repository

import (
	"context"
	"errors"
	"time"

	types "github.com/more-than-code/deploybot-service-api/deploybot-types"
	"github.com/more-than-code/deploybot-service-api/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (r *Repository) CreateMember(ctx context.Context, input types.CreateMemberInput) error {
	member := util.StructToBsonDoc(input.Member)
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

func (r *Repository) DeleteMember(ctx context.Context, input types.DeleteMemberInput) error {
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

func (r *Repository) UpdateMember(ctx context.Context, input types.UpdateMemberInput) error {
	filter := bson.M{"_id": input.ProjectId, "members.userid": input.UserId, "owneruserid": bson.M{"$ne": input.UserId}}

	update := bson.M{"$set": bson.M{"members.$.role": input.Member, "members.$.updatedat": primitive.NewDateTimeFromTime(time.Now().UTC())}}

	coll := r.mongoClient.Database("pipeline").Collection("projects")
	_, err := coll.UpdateOne(ctx, filter, update)

	if err != nil {
		return err
	}

	return nil
}
