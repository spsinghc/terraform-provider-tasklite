package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"terraform-provider-tasklite/internal/task"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ provider.Provider = &taskLiteProvider{}
)

// New is a helper function to simplify provider server and testing implementation.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &taskLiteProvider{
			version: version,
		}
	}
}

// Metadata returns the provider type name.
func (p *taskLiteProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "tasklite"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *taskLiteProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				Description: "URI for TaskLite API.",
				Required:    true,
			},
		},
	}
}

// Configure prepares a taskLite API client for data sources and resources.
func (p *taskLiteProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Debug(ctx, "Configuring Tasklite client")
	// Retrieve provider data from configuration
	var config hashicupsProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.Host.IsUnknown() || config.Host.IsNull() || config.Host.ValueString() == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"TaskLite API Host is unknown or empty",
			"The provider cannot create the TaskLite API client as there is an invalid configuration value for the TaskLite API host. "+
				"Either target apply the source of the value first, set the value statically in the configuration",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Create a new task client using the configuration values
	client := task.NewClient(config.Host.ValueString())

	resp.ResourceData = client

	ctx = tflog.SetField(ctx, "Tasklite host", config.Host)
	tflog.Debug(ctx, "Configured Tasklite client", map[string]any{"success": true})
}

// DataSources defines the data sources implemented in the provider.
func (p *taskLiteProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return nil
}

// Resources defines the resources implemented in the provider.
func (p *taskLiteProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewTaskResource,
	}
}
