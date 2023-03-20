package repository

import (
	"context"
	"strings"
	"time"

	"github.com/more-than-code/deploybot-service-api/model"
	"github.com/more-than-code/deploybot-service-api/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (r *Repository) CreateUser(ctx context.Context, input *model.CreateUserInput) error {
	hashedPassword, err := util.HashPassword(input.Password)

	if err != nil {
		return err
	}

	coll := r.mongoClient.Database("pipeline").Collection("users")

	doc := util.StructToBsonDoc(input)
	doc["password"] = hashedPassword
	doc["createdat"] = primitive.NewDateTimeFromTime(time.Now().UTC())

	_, err = coll.InsertOne(ctx, doc)

	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetUsers(ctx context.Context, input model.GetUsersInput) (*model.GetUsersOutput, error) {
	coll := r.mongoClient.Database("pipeline").Collection("users")

	filter := bson.M{"_id": bson.M{"$in": input.UserIds}}

	opts := options.FindOptions{Projection: bson.M{"password": 0}}
	cursor, err := coll.Find(ctx, filter, &opts)

	if err != nil {
		return nil, err
	}

	var output model.GetUsersOutput
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

func (r *Repository) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	coll := r.mongoClient.Database("pipeline").Collection("users")

	filter := bson.M{}
	filter["email"] = strings.ToLower(email)

	user := &model.User{}

	err := coll.FindOne(ctx, filter).Decode(user)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *Repository) GetUserById(ctx context.Context, id primitive.ObjectID) (*model.User, error) {
	coll := r.mongoClient.Database("pipeline").Collection("users")

	user := &model.User{}

	opts := options.FindOneOptions{Projection: bson.M{"password": 0}}

	err := coll.FindOne(ctx, bson.M{"_id": id}, &opts).Decode(user)

	if err != nil {
		return nil, err
	}

	return user, nil
}
