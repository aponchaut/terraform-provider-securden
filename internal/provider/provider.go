package provider

import (
	"context"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	//"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	//"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var _ provider.Provider = &securdenProvider{}

// Provider is a helper function to simplify provider server and testing implementation.
// function is called in main.go
func Provider(version string) func() provider.Provider {
	return func() provider.Provider {
		return &securdenProvider{
			version: version,
		}
	}
}

// Inputs variables --> TO CHECK
//var SecurdenAuthToken string
//var SecurdenServerURL string
//var server_url string
//var SecurdenTimeout int64

// securdenProviderModel maps provider schema data to a Go type.
type securdenProviderModel struct {
	AuthToken types.String `tfsdk:"authtoken"`
	ServerURL types.String `tfsdk:"server_url"`
	//ServerTimeout types.Int64 `tfsdk:"server_timeout"`
	ServerTimeout types.String `tfsdk:"server_timeout"`
}

// securdenProvider is the provider implementation.
type securdenProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// Securden Client to API
type securdenClient struct {
	ServerUrl string
	AuthToken string
	Timeout   time.Duration
	Client    *http.Client
}

// Metadata returns the provider type name.
func (p *securdenProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "securden"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *securdenProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Interact with Securden.",
		Attributes: map[string]schema.Attribute{
			"authtoken": schema.StringAttribute{
				Required:    true,
				Sensitive:   true,
				Description: "Securden API Authentication Token",
			},
			"server_url": schema.StringAttribute{
				Required:    true,
				Description: "Securden Server URL. Example: https://example.securden.com",
			},
			//"server_timeout": schema.Int64Attribute{
			"server_timeout": schema.StringAttribute{
				Optional:    true,
				Description: "Securden Server Timeout in seconds (default=30)",
			},
		},
	}
}

// Configure prepares a securden client for data sources and resources.
func (p *securdenProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring Securden client")
	// Retrieve provider data from configuration
	var config securdenProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If practitioner provided a configuration value for any of the
	// attributes, it must be a known value.
	if config.ServerURL.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("server_url"),
			"Securden URL undefined",
			"Securden URL must be defined",
		)
	}

	if config.AuthToken.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("authtoken"),
			"Auth token undefined",
			"Authentication Token must be defined",
		)
	}

	// Exit if error
	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.
	authtoken := os.Getenv("SECURDEN_AUTHTOKEN")
	server_url := os.Getenv("SECURDEN_URL")
	server_timeout := os.Getenv("SECURDEN_TIMEOUT")
	//timeout := int64(30)

	// Check inputs
	if !config.ServerURL.IsNull() {
		server_url = config.ServerURL.ValueString()
	}
	if !config.AuthToken.IsNull() {
		authtoken = config.AuthToken.ValueString()
	}
	if !config.ServerTimeout.IsNull() {
		//server_timeout = config.ServerTimeout.ValueInt64()
		server_timeout = config.ServerTimeout.ValueString()
	}

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.
	if server_url == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("server_url"),
			"Missing Securden URL",
			"The provider cannot create the Securden client as there is a missing or empty value for the Securden URL. "+
				"Set the host value in the configuration or use the SECURDEN_URL environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if authtoken == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("authtoken"),
			"Missing Authentication Token to Securden",
			"The provider cannot create the Securden client as there is a missing or empty value for the Securden Authentication Token. "+
				"Set the host value in the configuration or use the SECURDEN_AUTHTOKEN environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if server_timeout == "" {
		//force default value to 30 seconds
		server_timeout = "30"
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Timeout calculation
	seconds, _ := strconv.ParseInt(server_timeout, 0, 64)
	timeoutDuration := time.Duration(seconds) * time.Second

	// Create a new Securden client using the configuration values
	httpClient := &http.Client{
		Timeout: timeoutDuration,
	}

	// Create and return the client
	client := &securdenClient{
		ServerUrl: server_url,
		AuthToken: authtoken,
		Timeout:   timeoutDuration,
		Client:    httpClient,
	}

	// Make the Securden client available during DataSource and Resource
	// type Configure methods.
	resp.DataSourceData = client
	resp.ResourceData = client
}

// Resources defines the resources implemented in the provider.
func (p *securdenProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{}
}

// DataSources defines the data sources implemented in the provider.
func (p *securdenProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		account_get_password_data_source,
	}
}
