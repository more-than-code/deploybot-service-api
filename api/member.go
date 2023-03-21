package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	types "github.com/more-than-code/deploybot-service-api/deploybot-types"
)

func (a *Api) PostMember() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var input types.CreateMemberInput
		err := ctx.BindJSON(&input)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, types.PostMemberResponse{Code: types.CodeClientError, Msg: err.Error()})
			return
		}

		err = a.repo.CreateMember(ctx, input)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, types.PostMemberResponse{Code: types.CodeServerError, Msg: err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, types.PostMemberResponse{})
	}

}

func (a *Api) DeleteMember() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		input := types.DeleteMemberInput{}
		err := ctx.BindJSON(&input)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, types.DeleteMemberResponse{Code: types.CodeClientError, Msg: err.Error()})
			return
		}

		err = a.repo.DeleteMember(ctx, input)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, types.DeleteMemberResponse{Code: types.CodeClientError, Msg: err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, types.DeleteMemberResponse{})
	}
}

func (a *Api) PatchMember() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var input types.UpdateMemberInput
		err := ctx.BindJSON(&input)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, types.PatchMemberResponse{Code: types.CodeClientError, Msg: err.Error()})
			return
		}

		err = a.repo.UpdateMember(ctx, input)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, types.PatchMemberResponse{Code: types.CodeServerError, Msg: err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, types.PatchMemberResponse{})
	}
}
