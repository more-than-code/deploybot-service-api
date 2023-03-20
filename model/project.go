package model

import "go.mongodb.org/mongo-driver/bson/primitive"

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

type GetProjectsOutput struct {
	Items      []Project `json:"items"`
	TotalCount int64     `json:"totalCount"`
}
