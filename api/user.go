package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	types "github.com/more-than-code/deploybot-service-api/deploybot-types"
	"github.com/more-than-code/deploybot-service-api/repository"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (a *Api) Authenticate() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var input repository.AuthenticationInput
		err := ctx.BindJSON(&input)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, PostPipelineResponse{Code: types.CodeClientError, Msg: err.Error()})
			return
		}

		user, err := a.repo.GetUserByEmail(ctx, input.Email)

		res := AuthenticationResponse{}

		if user == nil {
			log.Println(err)
			res.Code = types.CodeWrongEmailOrPassword
			res.Msg = types.MsgWrongEmailOrPassword
			ctx.JSON(http.StatusBadRequest, res)
			return
		}

		err = repository.CheckPasswordHash(input.Password, user.Password)
		if err != nil {
			log.Println(err)
			res.Code = types.CodeWrongEmailOrPassword
			res.Msg = types.MsgWrongEmailOrPassword
			ctx.JSON(http.StatusBadRequest, res)
			return
		}

		partialUser := &repository.User{Id: user.Id}
		bytes, _ := json.Marshal(partialUser)

		output := repository.AuthenticationOutput{}
		at, err := a.atHelper.Authenticate(string(bytes))

		if err != nil {
			log.Println(err)
			res.Code = types.CodeAuthenticationFailure
			res.Msg = types.MsgAuthenticationFailure
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

func (a *Api) AuthenticateSso() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var input repository.AuthenticationSsoInput
		err := ctx.BindJSON(&input)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, PostPipelineResponse{Code: types.CodeClientError, Msg: err.Error()})
			return
		}

		claims, err := repository.GetGoogleAuthClaims(a.googleClientId, input.IdToken)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, PostPipelineResponse{Code: types.CodeClientError, Msg: err.Error()})
			return
		}

		user, err := a.repo.GetOrCreateUserBySubject(ctx, claims)

		res := AuthenticationResponse{}

		if err != nil {
			log.Println(err)
			res.Code = types.CodeClientError
			res.Msg = err.Error()
			ctx.JSON(http.StatusBadRequest, res)
			return
		}

		partialUser := &repository.User{Id: user.Id}
		bytes, _ := json.Marshal(partialUser)

		output := repository.AuthenticationOutput{}
		at, err := a.atHelper.Authenticate(string(bytes))

		if err != nil {
			log.Println(err)
			res.Code = types.CodeAuthenticationFailure
			res.Msg = types.MsgAuthenticationFailure
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
		var input repository.CreateUserInput
		err := ctx.BindJSON(&input)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, PostUserResponse{Code: types.CodeClientError, Msg: err.Error()})
			return
		}

		// TODO: validate verification code instead of hard-coding
		if input.VerificationCode != "1235" {
			ctx.JSON(http.StatusBadRequest, PostUserResponse{Code: types.CodeWrongVerificationCode, Msg: types.MsgWrongVerificationCode})
		}

		err = a.repo.CreateUser(ctx, &input)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, PostUserResponse{Code: types.CodeServerError, Msg: err.Error()})
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
			ctx.JSON(http.StatusBadRequest, DeleteUserResponse{Code: types.CodeClientError, Msg: err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, DeleteUserResponse{})
	}
}

func (a *Api) GetUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		idStr := ctx.Query("id")

		var id primitive.ObjectID
		if idStr == "" {
			id = repository.GetUserFromContext(ctx).Id
		} else {
			id, _ = primitive.ObjectIDFromHex(idStr)
		}
		user, err := a.repo.GetUserById(ctx, id)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, GetUserResponse{Code: types.CodeClientError, Msg: err.Error()})
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

		output, err := a.repo.GetUsers(ctx, repository.GetUsersInput{UserIds: uidList})

		if err != nil {
			ctx.JSON(http.StatusBadRequest, GetUsersResponse{Code: types.CodeClientError, Msg: err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, GetUsersResponse{Payload: output})
	}
}
