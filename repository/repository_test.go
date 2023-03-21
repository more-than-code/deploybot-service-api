package repository

import (
	"context"
	"testing"

	types "github.com/more-than-code/deploybot-service-api/deploybot-types"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestGetTasks(t *testing.T) {
	r, _ := NewRepository()

	tasks, err := r.GetPipelines(context.TODO(), types.GetPipelinesInput{})
	if err != nil {
		t.Fatal(err)
	}

	t.Log(tasks)
}

func TestGetPipeline(t *testing.T) {
	r, _ := NewRepository()

	pName := "geoy-webapp"
	aRun := true
	tId, _ := primitive.ObjectIDFromHex("6363bebf3ad85d86c5e2a5c8")

	pl, err := r.GetPipeline(context.TODO(), types.GetPipelineInput{Name: pName, TaskFilter: types.TaskFilter{UpstreamTaskId: &tId, AutoRun: &aRun}})

	if err != nil {
		t.Fatal(err)
	}

	if pl != nil {
		t.Log(pl.Tasks)
	}
}
