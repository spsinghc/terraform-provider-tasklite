package task

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setupTestServer(t *testing.T, method string, response interface{}, statusCode int) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			t.Fatalf("Expected method %s, got %s", method, r.Method)
		}
		w.WriteHeader(statusCode)
		if response != nil {
			json.NewEncoder(w).Encode(response)
		}
	}))
}

func TestCreateTask(t *testing.T) {
	task := Task{Title: "Test Task"}
	taskResponse := Task{ID: 1, Title: "Test Task", Priority: 0, Complete: false}

	server := setupTestServer(t, http.MethodPost, taskResponse, http.StatusCreated)
	defer server.Close()

	createdTask, err := (NewClient(server.URL)).CreateTask(context.Background(), task)
	assert.NoError(t, err)
	assert.Equal(t, *createdTask, taskResponse)
}

func TestReadTask(t *testing.T) {
	taskResponse := Task{ID: 1, Title: "Test Task", Priority: 0, Complete: false}

	server := setupTestServer(t, http.MethodGet, taskResponse, http.StatusCreated)
	defer server.Close()

	task, err := (NewClient(server.URL)).ReadTask(context.Background(), 1)

	assert.NoError(t, err)
	assert.Equal(t, *task, taskResponse)
}

func TestUpdateTask(t *testing.T) {
	task := Task{ID: 1, Title: "Updated Task", Priority: 1, Complete: true}

	taskResponse := Task{ID: 1, Title: "Updated Task", Priority: 1, Complete: true}

	server := setupTestServer(t, http.MethodPut, taskResponse, http.StatusCreated)
	defer server.Close()

	updatedTask, err := (NewClient(server.URL)).UpdateTask(context.Background(), task)
	assert.NoError(t, err)

	if updatedTask.Title != taskResponse.Title {
		t.Errorf("Expected task title %s, got %s", taskResponse.Title, updatedTask.Title)
	}

	assert.Equal(t, *updatedTask, taskResponse)
}

func TestDeleteTask(t *testing.T) {
	server := setupTestServer(t, http.MethodDelete, nil, http.StatusOK)
	defer server.Close()

	err := (NewClient(server.URL)).DeleteTask(context.Background(), 1)
	assert.NoError(t, err)
}
