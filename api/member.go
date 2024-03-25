package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	types "github.com/more-than-code/deploybot-service-api/deploybot-types"
	"github.com/more-than-code/deploybot-service-api/repository"
)

func (a *Api) PostMember() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var input repository.CreateMemberInput
		err := ctx.BindJSON(&input)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, PostMemberResponse{Code: types.CodeClientError, Msg: err.Error()})
			return
		}

		err = a.repo.CreateMember(ctx, input)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, PostMemberResponse{Code: types.CodeServerError, Msg: err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, PostMemberResponse{})
	}

}

func (a *Api) DeleteMember() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		input := repository.DeleteMemberInput{}
		err := ctx.BindJSON(&input)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, DeleteMemberResponse{Code: types.CodeClientError, Msg: err.Error()})
			return
		}

		err = a.repo.DeleteMember(ctx, input)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, DeleteMemberResponse{Code: types.CodeClientError, Msg: err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, DeleteMemberResponse{})
	}
}

func (a *Api) PatchMember() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var input repository.UpdateMemberInput
		err := ctx.BindJSON(&input)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, PatchMemberResponse{Code: types.CodeClientError, Msg: err.Error()})
			return
		}

		err = a.repo.UpdateMember(ctx, input)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, PatchMemberResponse{Code: types.CodeServerError, Msg: err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, PatchMemberResponse{})
	}
}
