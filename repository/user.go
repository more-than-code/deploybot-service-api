package repository

import (
	"context"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AuthenticationInput struct {
	Email    string
	Password string
}

type AuthenticationSsoInput struct {
	IdToken string
}

type AuthenticationOutput struct {
	UserId       primitive.ObjectID `json:"userId"`
	AccessToken  string             `json:"accessToken"`
	RefreshToken string             `json:"refreshToken"`
}

type User struct {
	Id           primitive.ObjectID `json:"id" bson:"_id"`
	Subject      string             `json:"subject"`
	Email        string             `json:"email"`
	ContactEmail string             `json:"contactEmail"`
	Password     string             `json:"password"`
	Name         string             `json:"name"`
	AvatarUrl    string             `json:"avatarUrl"`
	CreatedAt    primitive.DateTime `json:"createdAt"`
}

type CreateUserInput struct {
	Name             string
	Email            string
	Password         string
	VerificationCode string
}

type GetUsersInput struct {
	UserIds []primitive.ObjectID
}

type GetUsersOutput struct {
	Items      []User `json:"items"`
	TotalCount int    `json:"totalCount"`
}

type Claims struct {
	Iss           string
	Nbf           int64
	Aud           string
	Sub           string
	Email         string
	EmailVerified bool `json:"email_verified"`
	Azp           string
	Name          string
	Picture       string
	GivenName     string `json:"given_name"`
	Iat           int64
	Exp           int64
	Jti           string
}

func (r *Repository) CreateUser(ctx context.Context, input *CreateUserInput) error {
	hashedPassword, err := HashPassword(input.Password)

	if err != nil {
		return err
	}

	coll := r.mongoClient.Database("pipeline").Collection("users")

	doc := StructToBsonDoc(input)
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

func (r *Repository) GetUsers(ctx context.Context, input GetUsersInput) (*GetUsersOutput, error) {
	coll := r.mongoClient.Database("pipeline").Collection("users")

	filter := bson.M{"_id": bson.M{"$in": input.UserIds}}

	opts := options.FindOptions{Projection: bson.M{"password": 0}}
	cursor, err := coll.Find(ctx, filter, &opts)

	if err != nil {
		return nil, err
	}

	var output GetUsersOutput
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

func (r *Repository) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	coll := r.mongoClient.Database("pipeline").Collection("users")

	filter := bson.M{}
	filter["email"] = strings.ToLower(email)

	user := &User{}

	err := coll.FindOne(ctx, filter).Decode(user)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *Repository) GetOrCreateUserBySubject(ctx context.Context, claims *Claims) (*User, error) {
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

	user := User{}

	err := coll.FindOneAndUpdate(ctx, filter, update, &opts).Decode(&user)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *Repository) GetUserById(ctx context.Context, id primitive.ObjectID) (*User, error) {
	coll := r.mongoClient.Database("pipeline").Collection("users")

	user := &User{}

	opts := options.FindOneOptions{Projection: bson.M{"password": 0}}

	err := coll.FindOne(ctx, bson.M{"_id": id}, &opts).Decode(user)

	if err != nil {
		return nil, err
	}

	return user, nil
}
