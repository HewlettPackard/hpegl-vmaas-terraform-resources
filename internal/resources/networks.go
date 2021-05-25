package resources

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const networkReadTimeout = 30 * time.Second

func NetworkData() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				Description: `Name of the network. This needs to be exact name or
				else will return error not found`,
			},
		},
		ReadContext: networkReadContext,
		Description: "Get the Network details",
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(networkReadTimeout),
		},
		SchemaVersion:  0,
		StateUpgraders: nil,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func networkReadContext(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	return nil
}
