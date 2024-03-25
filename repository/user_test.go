package repository

import (
	"context"
	"testing"
)

func TestGetOrCreateUserBySubject(t *testing.T) {
	r, _ := NewRepository()

	user, err := r.GetOrCreateUserBySubject(context.TODO(), &Claims{Sub: "123456"})

	if err != nil {
		t.Error(err)
	}

	t.Log(user)
}
