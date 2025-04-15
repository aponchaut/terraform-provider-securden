// Copyright (c) HashiCorp, Inc.

package provider

import (
	"context"
	"encoding/json"
	"fmt"
	//"log"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

//example
//curl -k -H "authtoken: xxxxxx" -X GET "https://bel.securden-vault.com/api/get_password?account_id=1000000027810&reason=Test"

// Ensure provider defined types fully satisfy framework interfaces.
// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &accountDataSource{}
	_ datasource.DataSourceWithConfigure = &accountDataSource{}
)

// account_data_souorce is a helper function to simplify the provider implementation.
func accountGetPasswordDataSource() datasource.DataSource {
	return &accountDataSource{}
}

// accountDataSource defines the data source implementation.
type accountDataSource struct {
	client *securdenClient
}

// Account Data Model
type accountDataModel struct {
	//	ID          types.String `tfsdk:"id"`
	AccountID types.String `tfsdk:"account_id"`
	Password  types.String `tfsdk:"password"`
	Label     types.String `tfsdk:"label"`
}

// API response attended
type GetpasswordResponse struct {
	Password string `json:"password"`
	Label    string `json:"label"`
}

// Account implementation
func (d *accountDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_account_get_password"
}

// Schema defines the schema for the data source.
func (d *accountDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Get password from Securden account ID",
		Attributes: map[string]schema.Attribute{
			//"id": schema.StringAttribute{
			//	Computed:    true,
			//	Description: "Unique identifier for this resource",
			//},
			"account_id": schema.StringAttribute{
				Required:    true,
				Description: "account ID",
			},
			"password": schema.StringAttribute{
				Computed:    true,
				Sensitive:   true,
				Description: "Account Password",
			},
			"label": schema.StringAttribute{
				Computed:    true,
				Description: "Label",
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *accountDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*securdenClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *securden.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.client = client
}

// Read refreshes the Terraform state with the latest data.
func (d *accountDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data accountDataModel

	// configuration reading ...
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Defined apiURL
	apiURL := "api/get_password"

	// builds params
	//params := make(map[string]string)
	params := map[string]string{
		"account_id": data.AccountID.ValueString(),
	}

	// Call getRequest function get_request
	responseBody, err := d.client.getRequest(ctx, apiURL, params)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error on Get Request",
			"Impossible to execute the API request: "+err.Error(),
		)
		return
	}

	// Decode JSON response (GO version)
	var apiResponse GetpasswordResponse
	err = json.Unmarshal(responseBody, &apiResponse)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error on parsing response",
			"Impossible to decode JSON response: "+err.Error(),
		)
		return
	}

	// Set values in the data object
	//data.ID = types.StringValue(data.AccountID.ValueString())
	data.Password = types.StringValue(apiResponse.Password)
	data.Label = types.StringValue(apiResponse.Label)

	// Set the state
	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
