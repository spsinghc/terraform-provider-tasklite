package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"

	"terraform-provider-tasklite/internal/task"
)

// taskLiteProvider is the provider implementation.
type taskLiteProvider struct {
	version string
}

// hashicupsProviderModel maps provider schema data to a Go type.
type hashicupsProviderModel struct {
	Host types.String `tfsdk:"host"`
}

type taskModel struct {
	ID       types.Int32  `tfsdk:"id"`
	Title    types.String `tfsdk:"title"`
	Priority types.Int32  `tfsdk:"priority"`
	Complete types.Bool   `tfsdk:"complete"`
}

// mapTaskToTaskModel maps api client task struct to provider task type.
func mapTaskToTaskModel(t *task.Task) taskModel {
	return taskModel{
		ID:       types.Int32Value(t.ID),
		Title:    types.StringValue(t.Title),
		Priority: types.Int32Value(t.Priority),
		Complete: types.BoolValue(t.Complete),
	}
}

// mapTaskModelToTask maps provider task type to api client task struct.
func mapTaskModelToTask(t taskModel) task.Task {
	return task.Task{
		ID:       t.ID.ValueInt32(),
		Title:    t.Title.ValueString(),
		Priority: t.Priority.ValueInt32(),
		Complete: t.Complete.ValueBool(),
	}
}
