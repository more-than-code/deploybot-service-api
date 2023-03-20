package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type AuthenticationInput struct {
	Email    string
	Password string
}

type AuthenticationOutput struct {
	UserId       primitive.ObjectID `json:"userId"`
	AccessToken  string             `json:"accessToken"`
	RefreshToken string             `json:"refreshToken"`
}

type User struct {
	Id        primitive.ObjectID `json:"id" bson:"_id"`
	Email     string             `json:"email"`
	Password  string             `json:"password"`
	CreatedAt primitive.DateTime `json:"createdAt"`
	AvatarUrl string             `json:"avatarUrl"`
}

type CreateUserInput struct {
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
