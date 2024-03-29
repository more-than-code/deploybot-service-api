package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"

	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/api/idtoken"
)

func StructToBsonDoc(source interface{}) bson.M {
	bytes, err := bson.Marshal(source)

	if err != nil {
		return nil
	}

	doc := bson.M{}
	err = bson.Unmarshal(bytes, &doc)

	if err != nil {
		return nil
	}

	return doc
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err
}

func GetUserFromContext(gc *gin.Context) User {
	param := gc.Param("user")

	var user User

	json.Unmarshal([]byte(param), &user)

	return user
}

func GetGoogleAuthClaims(clientId, token string) (*Claims, error) {
	payload, err := idtoken.Validate(context.Background(), token, clientId)

	if err != nil {
		return nil, err
	}

	fmt.Print(payload)

	var claims Claims

	bs, err := json.Marshal(payload.Claims)

	if err != nil {
		return nil, err
	}

	json.Unmarshal(bs, &claims)

	return &claims, nil
}
