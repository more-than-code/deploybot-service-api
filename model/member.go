package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Member struct {
	User      User               `json:"user"`
	Role      Role               `json:"role"`
	CreatedAt primitive.DateTime `json:"datetimeCreated"`
}

type CreateMemberInput struct {
	ProjectId primitive.ObjectID
	Member    struct {
		UserId primitive.ObjectID
		Role   Role
	}
}

type DeleteMemberInput struct {
	ProjectId primitive.ObjectID
	UserId    primitive.ObjectID
}

type UpdateMember struct {
	Role *Role
}
type UpdateMemberInput struct {
	ProjectId primitive.ObjectID
	UserId    primitive.ObjectID
	Member    UpdateMember
}
