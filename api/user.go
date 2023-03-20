package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/more-than-code/deploybot-service-api/model"
	"github.com/more-than-code/deploybot-service-api/util"
)

func (a *Api) Authenticate() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var input model.AuthenticationInput
		err := ctx.BindJSON(&input)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, PostPipelineResponse{Code: CodeClientError, Msg: err.Error()})
			return
		}

		user, err := a.repo.GetUserByEmail(ctx, input.Email)

		res := AuthenticationResponse{}

		if user == nil {
			log.Println(err)
			res.Code = model.CodeWrongEmailOrPassword
			res.Msg = model.MsgWrongEmailOrPassword
			ctx.JSON(http.StatusBadRequest, res)
			return
		}

		err = util.CheckPasswordHash(input.Password, user.Password)
		if err != nil {
			log.Println(err)
			res.Code = model.CodeWrongEmailOrPassword
			res.Msg = model.MsgWrongEmailOrPassword
			ctx.JSON(http.StatusBadRequest, res)
			return
		}

		partialUser := &model.User{Id: user.Id}
		bytes, _ := json.Marshal(partialUser)

		output := model.AuthenticationOutput{}
		at, err := a.atHelper.Authenticate(string(bytes))

		if err != nil {
			log.Println(err)
			res.Code = model.CodeAuthenticationFailure
			res.Msg = model.MsgAuthenticationFailure
			ctx.JSON(http.StatusBadRequest, res)
			return
		}

		output.AccessToken = at

		rt, _ := a.rtHelper.Authenticate(user.Id.Hex())
		output.RefreshToken = rt
		output.UserId = user.Id

		res.Payload = &output

		ctx.JSON(http.StatusOK, res)
	}

}

func (a *Api) PostUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var input model.CreateUserInput
		err := ctx.BindJSON(&input)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, PostUserResponse{Code: CodeClientError, Msg: err.Error()})
			return
		}

		// TODO: validate verification code instead of hard-coding
		if input.VerificationCode != "1235" {
			ctx.JSON(http.StatusBadRequest, PostUserResponse{Code: model.CodeWrongVerificationCode, Msg: model.MsgWrongVerificationCode})
		}

		err = a.repo.CreateUser(ctx, &input)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, PostUserResponse{Code: CodeServerError, Msg: err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, PostUserResponse{})
	}

}

func (a *Api) DeleteUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")
		objId, _ := primitive.ObjectIDFromHex(id)

		err := a.repo.DeleteUser(ctx, objId)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, DeleteUserResponse{Code: CodeClientError, Msg: err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, DeleteUserResponse{})
	}
}

func (a *Api) GetUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")

		objId, _ := primitive.ObjectIDFromHex(id)
		user, err := a.repo.GetUserById(ctx, objId)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, GetUserResponse{Code: CodeClientError, Msg: err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, GetUserResponse{Payload: user})
	}
}

func (a *Api) GetUsers() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		uidStrList := ctx.QueryArray("uid")

		uidList := make([]primitive.ObjectID, len(uidStrList))
		for _, us := range uidStrList {
			uid, _ := primitive.ObjectIDFromHex(us)
			uidList = append(uidList, uid)
		}

		output, err := a.repo.GetUsers(ctx, model.GetUsersInput{UserIds: uidList})

		if err != nil {
			ctx.JSON(http.StatusBadRequest, GetUsersResponse{Code: CodeClientError, Msg: err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, GetUsersResponse{Payload: output})
	}
}
