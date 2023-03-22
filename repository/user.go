package repository

import (
	"context"
	"strings"
	"time"

	types "github.com/more-than-code/deploybot-service-api/deploybot-types"
	"github.com/more-than-code/deploybot-service-api/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (r *Repository) CreateUser(ctx context.Context, input *types.CreateUserInput) error {
	hashedPassword, err := util.HashPassword(input.Password)

	if err != nil {
		return err
	}

	coll := r.mongoClient.Database("pipeline").Collection("users")

	doc := util.StructToBsonDoc(input)
	doc["name"] = input.Name
	doc["subject"] = input.Email
	doc["password"] = hashedPassword
	doc["createdat"] = primitive.NewDateTimeFromTime(time.Now().UTC())

	_, err = coll.InsertOne(ctx, doc)

	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetUsers(ctx context.Context, input types.GetUsersInput) (*types.GetUsersOutput, error) {
	coll := r.mongoClient.Database("pipeline").Collection("users")

	filter := bson.M{"_id": bson.M{"$in": input.UserIds}}

	opts := options.FindOptions{Projection: bson.M{"password": 0}}
	cursor, err := coll.Find(ctx, filter, &opts)

	if err != nil {
		return nil, err
	}

	var output types.GetUsersOutput
	if err = cursor.All(ctx, &output.Items); err != nil {
		return nil, err
	}

	output.TotalCount = len(output.Items)

	return &output, nil
}

func (r *Repository) DeleteUser(ctx context.Context, id primitive.ObjectID) error {
	coll := r.mongoClient.Database("pipeline").Collection("users")

	_, err := coll.DeleteOne(ctx, bson.M{"_id": id})

	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetUserByEmail(ctx context.Context, email string) (*types.User, error) {
	coll := r.mongoClient.Database("pipeline").Collection("users")

	filter := bson.M{}
	filter["email"] = strings.ToLower(email)

	user := &types.User{}

	err := coll.FindOne(ctx, filter).Decode(user)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *Repository) GetOrCreateUserBySubject(ctx context.Context, claims *types.Claims) (*types.User, error) {
	coll := r.mongoClient.Database("pipeline").Collection("users")

	filter := bson.M{}
	filter["subject"] = claims.Sub

	doc := bson.M{}
	doc["subject"] = claims.Sub
	doc["email"] = claims.Sub
	doc["contactemail"] = claims.Email
	doc["avatarurl"] = claims.Picture
	doc["name"] = claims.Name
	doc["createdat"] = primitive.NewDateTimeFromTime(time.Now().UTC())

	update := bson.M{"$set": doc}

	upsert := true
	after := options.After
	opts := options.FindOneAndUpdateOptions{Upsert: &upsert, ReturnDocument: &after}

	user := types.User{}

	err := coll.FindOneAndUpdate(ctx, filter, update, &opts).Decode(&user)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *Repository) GetUserById(ctx context.Context, id primitive.ObjectID) (*types.User, error) {
	coll := r.mongoClient.Database("pipeline").Collection("users")

	user := &types.User{}

	opts := options.FindOneOptions{Projection: bson.M{"password": 0}}

	err := coll.FindOne(ctx, bson.M{"_id": id}, &opts).Decode(user)

	if err != nil {
		return nil, err
	}

	return user, nil
}
