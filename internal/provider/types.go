package provider

import "github.com/hashicorp/terraform-plugin-framework/types"

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

type taskModel struct {
	ID    types.Int32  `tfsdk:"id"`
	Title types.String `tfsdk:"title"`
}

type tasksModel struct {
	Tasks []taskModel `tfsdk:"tasks"`
}
