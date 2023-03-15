package util

import (
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
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
