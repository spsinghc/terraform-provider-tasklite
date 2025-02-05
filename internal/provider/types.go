package provider

import (
	"terraform-provider-tasklite/internal/task"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

// taskLiteProvider is the provider implementation.
type taskLiteProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// hashicupsProviderModel maps provider schema data to a Go type.
type hashicupsProviderModel struct {
	Host types.String `tfsdk:"host"`
}

// TODO update model to remaining properties; completed, priority
type taskModel struct {
	ID        types.Int32  `tfsdk:"id"`
	Title     types.String `tfsdk:"title"`
	Priority  types.Int32  `tfsdk:"priority"`
	Completed types.Bool   `tfsdk:"completed"`
}

type tasksModel struct {
	Tasks []taskModel `tfsdk:"tasks"`
}

// mapTaskToTaskModel maps api client task struct to provider task type
func mapTaskToTaskModel(t *task.Task) taskModel {
	return taskModel{
		ID:        types.Int32Value(t.ID),
		Title:     types.StringValue(t.Title),
		Priority:  types.Int32Value(t.Priority),
		Completed: types.BoolValue(t.Completed),
	}
}

// mapTaskModelToTask maps provider task type to api client task struct
func mapTaskModelToTask(t taskModel) task.Task {
	return task.Task{
		ID:        t.ID.ValueInt32(),
		Title:     t.Title.ValueString(),
		Priority:  t.Priority.ValueInt32(),
		Completed: t.Completed.ValueBool(),
	}
}
