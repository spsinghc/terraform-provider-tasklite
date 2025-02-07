package provider

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const (
	providerConfigTemplate = `
provider "tasklite" {
  host = "%s"
}`
)

func newResourceServer(t *testing.T) *httptest.Server {
	var data atomic.Value
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodPost:
			body, _ := io.ReadAll(r.Body)
			parsedBody := make(map[string]interface{})
			_ = json.Unmarshal(body, &parsedBody)
			parsedBody["id"] = 1
			body, _ = json.Marshal(parsedBody)
			data.Store(body)
			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write(body)
		case r.Method == http.MethodGet:
			if storedData, ok := data.Load().([]byte); ok {
				_, _ = w.Write(storedData)
			} else {
				t.Fatal("Failed to assert type []byte for data.Load()")
			}
		case r.Method == http.MethodPut:
			body, _ := io.ReadAll(r.Body)
			parsedBody := make(map[string]interface{})
			_ = json.Unmarshal(body, &parsedBody)
			parsedBody["id"] = 1
			body, _ = json.Marshal(parsedBody)
			data.Store(body)
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(body)
		case r.Method == http.MethodDelete:
			w.WriteHeader(http.StatusNoContent)
			data.Store([]byte(nil))
		default:
			t.Fatal("Unexpected " + r.Method + " " + r.RequestURI)
		}
	}))
}

func TestAccTaskResource(t *testing.T) {
	server := newResourceServer(t)
	defer server.Close()
	resourceName := "tasklite_task.test"
	title := "Task created by terraform"
	updatedTitle := "Updated Task by terraform"
	c := fmt.Sprintf(providerConfigTemplate, server.URL)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Step 1 - create.
			{
				Config: fmt.Sprintf(`%s
resource "tasklite_task" "test" {
   title = "%s"
}
`, c, title),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "title", title),
					resource.TestCheckResourceAttr(resourceName, "priority", "0"),
					resource.TestCheckResourceAttr(resourceName, "complete", "false"),
				),
			},
			// Step 2 - update.
			{
				Config: fmt.Sprintf(`%s
resource "tasklite_task" "test" {
   title = "%s"
   priority = 1
   complete = true
}
`, c, updatedTitle),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "title", updatedTitle),
					resource.TestCheckResourceAttr(resourceName, "priority", "1"),
					resource.TestCheckResourceAttr(resourceName, "complete", "true"),
				),
			},
			// Step 3 - make no changes, check plan is empty.
			{
				Config: fmt.Sprintf(`%s
resource "tasklite_task" "test" {
   title = "%s"
   priority = 1
   complete = true
}
`, c, updatedTitle),
				PlanOnly: true,
			},
		},
	})
}
