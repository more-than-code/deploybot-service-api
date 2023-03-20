package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/more-than-code/deploybot-service-api/model"
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
		input := model.DeleteMemberInput{}
		err := ctx.BindJSON(&input)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, DeleteMemberResponse{Code: CodeClientError, Msg: err.Error()})
			return
		}

		err = a.repo.DeleteMember(ctx, input)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, DeleteMemberResponse{Code: CodeClientError, Msg: err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, DeleteMemberResponse{})
	}
}

func (a *Api) PatchMember() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var input model.UpdateMemberInput
		err := ctx.BindJSON(&input)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, PatchMemberResponse{Code: CodeClientError, Msg: err.Error()})
			return
		}

		err = a.repo.UpdateMember(ctx, input)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, PatchMemberResponse{Code: CodeServerError, Msg: err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, PatchMemberResponse{})
	}
}
