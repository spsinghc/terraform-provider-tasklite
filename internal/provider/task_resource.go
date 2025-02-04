package provider

import (
	"context"
	"fmt"

	"terraform-provider-tasklite/internal/task"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
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
		},
	}
}

func (r *taskResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
	var data taskModel

	// Read Terraform data data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	res, err := r.client.CreateTask(ctx, data.Title.ValueString())

	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to create task",
			err.Error(),
		)
		return
	}

	data.ID = types.Int32Value(res.ID)

	tflog.Trace(ctx, "created a resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, "An error return while saving state")
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *taskResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// not supported by the api
}

func (r *taskResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// not supported by the api
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *taskResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data taskModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Log the data being read
	tflog.Debug(ctx, "Read resource data", map[string]interface{}{
		"ID": data.ID.ValueInt32(),
	})

	err := r.client.DeleteTask(ctx, data.ID.ValueInt32())

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete task, got error: %s", err))
		return
	}
}
