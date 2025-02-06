package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccTaskResource(t *testing.T) {
	resourceName := "tasklite_task.test"
	title := "Task created by terraform"
	updatedTitle := "Updated Task by terraform"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Step 1 - create.
			{
				Config: fmt.Sprintf(`%s
resource "tasklite_task" "test" {
   title = "%s"
}
`, providerConfig, title),
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
`, providerConfig, updatedTitle),
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
`, providerConfig, updatedTitle),
				PlanOnly: true,
			},
		},
	})
}
