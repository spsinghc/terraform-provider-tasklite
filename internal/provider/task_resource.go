package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32default"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"terraform-provider-tasklite/internal/task"
)

var (
	_ resource.Resource              = &taskResource{}
	_ resource.ResourceWithConfigure = &taskResource{}
)

func NewTaskResource() resource.Resource {
	return &taskResource{}
}

type taskResource struct {
	client *task.Client
}

// Metadata returns the resource type name.
func (r *taskResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_task"
}

// Schema defines the schema for the resource.
func (r *taskResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"title": schema.StringAttribute{
				Required: true,
			},
			"id": schema.Int32Attribute{
				Computed: true,
			},
			"priority": schema.Int32Attribute{
				Optional: true,
				Computed: true,
				Default:  int32default.StaticInt32(0),
			},
			"complete": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(false),
			},
		},
	}
}

func (r *taskResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*task.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *tasklite.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client

}

// Create creates the resource and sets the initial Terraform state.
func (r *taskResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan taskModel
	// Read Terraform data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Creating task", map[string]any{"task": plan})
	t, err := r.client.CreateTask(ctx, mapTaskModelToTask(plan))

	if err != nil {
		logErrorAndAddDiagnostic(ctx, req, resp, err)
		return
	}

	tflog.Debug(ctx, "Task created", map[string]any{"task": t})
	s := mapTaskToTaskModel(t)
	resp.Diagnostics.Append(resp.State.Set(ctx, &s)...)

	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, "An error return while saving state")
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *taskResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state taskModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Refreshing task with the server data", map[string]interface{}{
		"ID": state.ID.ValueInt32(),
	})

	// get refreshed task from the api
	t, err := r.client.ReadTask(ctx, state.ID.ValueInt32())

	if err != nil {
		logErrorAndAddDiagnostic(ctx, req, resp, err)
		return
	}

	state = mapTaskToTaskModel(t)

	// set refreshed state
	resp.Diagnostics.Append(req.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *taskResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state taskModel
	// read Terraform state
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var plan taskModel
	// read Terraform plan into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if state.ID.ValueInt32() == 0 {
		logErrorAndAddDiagnostic(ctx, req, resp, fmt.Errorf("unexpected task ID %d found in the terraform state", state.ID.ValueInt32()))
		return
	}
	plan.ID = state.ID

	tflog.Debug(ctx, "Updating task", map[string]any{"task": plan})

	t, err := r.client.UpdateTask(ctx, mapTaskModelToTask(plan))

	if err != nil {
		logErrorAndAddDiagnostic(ctx, req, resp, err)
		return
	}

	plan = mapTaskToTaskModel(t)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, "An error return while saving state")
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *taskResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state taskModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Deleting task", map[string]interface{}{
		"ID": state.ID.ValueInt32(),
	})

	err := r.client.DeleteTask(ctx, state.ID.ValueInt32())

	if err != nil {
		logErrorAndAddDiagnostic(ctx, req, resp, err)
		return
	}
}

func logErrorAndAddDiagnostic(ctx context.Context, req any, resp any, err error) {
	operation := ""
	switch req.(type) {
	case resource.CreateRequest:
		operation = "Create"
		if respPtr, ok := resp.(*resource.CreateResponse); ok {
			respPtr.Diagnostics.AddError(fmt.Sprintf("%s Operation Error", operation), fmt.Sprintf("Failed to %s the task, got error: %s", operation, err))
		}
	case resource.ReadRequest:
		operation = "Read"
		if respPtr, ok := resp.(*resource.ReadResponse); ok {
			respPtr.Diagnostics.AddError(fmt.Sprintf("%s Operation Error", operation), fmt.Sprintf("Failed to %s the task, got error: %s", operation, err))
		}
	case resource.UpdateRequest:
		operation = "Update"
		if respPtr, ok := resp.(*resource.UpdateResponse); ok {
			respPtr.Diagnostics.AddError(fmt.Sprintf("%s Operation Error", operation), fmt.Sprintf("Failed to %s the task, got error: %s", operation, err))
		}
	case resource.DeleteRequest:
		operation = "Delete"
		if respPtr, ok := resp.(*resource.DeleteResponse); ok {
			respPtr.Diagnostics.AddError(fmt.Sprintf("%s Operation Error", operation), fmt.Sprintf("Failed to %s the task, got error: %s", operation, err))
		}
	}
	tflog.Error(ctx, fmt.Sprintf("Failed to %s the task", operation), map[string]interface{}{"error": err})
}
