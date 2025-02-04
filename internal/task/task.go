package task

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Task struct {
	ID    int32  `json:"id,omitempty"`
	Title string `json:"title"`
}

type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

func NewClient(baseURL string) *Client {
	return &Client{
		BaseURL:    baseURL,
		HTTPClient: &http.Client{},
	}
}

const TASK_URI = "/api/task/"

func apiPath(baseURL string) string {
	return fmt.Sprintf("%s%s", baseURL, TASK_URI)
}

func (c *Client) CreateTask(ctx context.Context, title string) (*Task, error) {
	taskRequest := &Task{Title: title}
	resp, err := c.doRequest(ctx, http.MethodPost, apiPath(c.BaseURL), taskRequest)
	if err != nil {
		return nil, err
	}

	var createdTask Task
	if err := c.parseResponse(resp, &createdTask); err != nil {
		return nil, err
	}

	return &createdTask, nil
}

func (c *Client) DeleteTask(ctx context.Context, id int32) error {

	url := fmt.Sprintf("%s%d/", apiPath(c.BaseURL), id)

	resp, err := c.doRequest(ctx, http.MethodDelete, url, nil)

	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to delete task: %s", resp.Status)
	}

	return nil
}

func (c *Client) doRequest(ctx context.Context, method, url string, body interface{}) (*http.Response, error) {
	var reqBody []byte
	var err error
	if body != nil {
		reqBody, err = json.Marshal(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *Client) parseResponse(resp *http.Response, out interface{}) error {
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, body)
	}

	return json.NewDecoder(resp.Body).Decode(out)
}
