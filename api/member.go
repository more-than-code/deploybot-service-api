package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/more-than-code/deploybot-service-api/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (a *Api) PostMember() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var input model.CreateMemberInput
		err := ctx.BindJSON(&input)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, PostMemberResponse{Code: CodeClientError, Msg: err.Error()})
			return
		}

		err = a.repo.CreateMember(ctx, input)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, PostMemberResponse{Code: CodeServerError, Msg: err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, PostMemberResponse{})
	}

}

func (a *Api) DeleteMember() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		pidStr := ctx.Param("pid")
		pid, _ := primitive.ObjectIDFromHex(pidStr)

		uidStr := ctx.Param("uid")
		uid, _ := primitive.ObjectIDFromHex(uidStr)

		err := a.repo.DeleteMember(ctx, model.DeleteMemberInput{ProjectId: pid, UserId: uid})

		if err != nil {
			ctx.JSON(http.StatusBadRequest, DeleteMemberResponse{Code: CodeClientError, Msg: err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, DeleteMemberResponse{})
	}
}

func (a *Api) PatchMember() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var member model.UpdateMember
		err := ctx.BindJSON(&member)

		pidStr := ctx.Param("pid")
		pid, _ := primitive.ObjectIDFromHex(pidStr)

		uidStr := ctx.Param("uid")
		uid, _ := primitive.ObjectIDFromHex(uidStr)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, PatchMemberResponse{Code: CodeClientError, Msg: err.Error()})
			return
		}

		err = a.repo.UpdateMember(ctx, model.UpdateMemberInput{ProjectId: pid, UserId: uid, Member: member})

		if err != nil {
			ctx.JSON(http.StatusBadRequest, PatchMemberResponse{Code: CodeServerError, Msg: err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, PatchMemberResponse{})
	}
}
