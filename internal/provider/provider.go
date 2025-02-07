package provider

import (
	"context"
	"os"

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
				Description: "URL for TaskLite API. May also be provided via TASKLITE_HOST environment variable.",
				Optional:    true,
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

	if config.Host.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Unknown API Host is unknown or empty",
			"The provider cannot create the TaskLite API client as there is an unknown configuration value for the TaskLite API host. "+
				"Either target apply the source of the value first, set the value statically in the configuration, "+
				"or use the TASKLITE_HOST environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.

	host := os.Getenv("TASKLITE_HOST")
	if !config.Host.IsNull() {
		host = config.Host.ValueString()
	}
	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.
	if host == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Missing TaskLite API Host",
			"The provider cannot create the TaskLite API client as there is a missing or empty value for the TaskLite API host. "+
				"Set the host value in the configuration or use the TASKLITE_HOST environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Create a new task client using the configuration values
	client := task.NewClient(host)

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
