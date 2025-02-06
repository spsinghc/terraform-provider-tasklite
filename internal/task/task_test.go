package task

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
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

func TestParseResponseError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = io.WriteString(w, `Bad Request`)
	}))

	defer server.Close()
	c := NewClient(server.URL)
	req, _ := http.NewRequest(http.MethodGet, server.URL, nil)
	resp, _ := c.HTTPClient.Do(req)

	var task Task
	err := c.parseResponse(resp, &task)
	assert.Error(t, err)
	e := fmt.Errorf("HTTP %d: %s", http.StatusBadRequest, "Bad Request")
	assert.Equal(t, e, err)
}

func TestParseResponseSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, `{"id":1,"title":"Test task", "priority":1, "complete":true}`)
	}))
	defer server.Close()

	c := NewClient(server.URL)
	req, _ := http.NewRequest(http.MethodGet, server.URL, nil)
	resp, _ := c.HTTPClient.Do(req)

	var task Task
	err := c.parseResponse(resp, &task)
	assert.NoError(t, err)

	assert.Equal(t, Task{
		ID:       1,
		Title:    "Test task",
		Priority: 1,
		Complete: true,
	}, task)
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
	server := setupTestServer(t, http.MethodDelete, nil, http.StatusNoContent)
	defer server.Close()

	err := (NewClient(server.URL)).DeleteTask(context.Background(), 1)
	assert.NoError(t, err)
}
