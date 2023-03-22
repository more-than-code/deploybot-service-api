package repository

import (
	"context"
	"testing"

	types "github.com/more-than-code/deploybot-service-api/deploybot-types"
)

func TestGetOrCreateUserBySubject(t *testing.T) {
	r, _ := NewRepository()

	user, err := r.GetOrCreateUserBySubject(context.TODO(), &types.Claims{Sub: "123456"})

	if err != nil {
		t.Error(err)
	}

	t.Log(user)
}
