// (C) Copyright 2022 Hewlett Packard Enterprise Development LP

package resources

import (
	"context"

	"github.com/HewlettPackard/hpegl-vmaas-terraform-resources/internal/resources/validations"
	"github.com/HewlettPackard/hpegl-vmaas-terraform-resources/internal/utils"
	"github.com/HewlettPackard/hpegl-vmaas-terraform-resources/pkg/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func LBProfileData() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: f(generalNamedesc, "ResLoadBalancerProfile", "ResLoadBalancerProfile"),
			},
			"service_type": {
				Type:     schema.TypeString,
				Required: true,
				ValidateDiagFunc: validations.StringInSlice([]string{
					"LBHttpProfile", "LBFastTcpProfile", "LBFastUdpProfile",
					"LBClientSslProfile", "LBServerSslProfile", "LBCookiePersistenceProfile", "LBGenericPersistenceProfile"}, false),
				Description: "Network Loadbalancer Supported values are `LBHttpProfile`,`LBFastTcpProfile`, `LBFastUdpProfile`, `LBClientSslProfile`, `LBServerSslProfile`, `LBCookiePersistenceProfile`,`LBGenericPersistenceProfile`",
			},
			"config": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "profile Configuration",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"profile_type": {
							Type: schema.TypeString,
							ValidateDiagFunc: validations.StringInSlice([]string{
								"application-profile", "ssl-profile", "persistence-profile",
							}, false),
							Required:    true,
							Description: "Network Loadbalancer Supported values are `application-profile`,`ssl-profile`, `persistence-profile`"},
						"ssl_suite": {
							Type: schema.TypeString,
							ValidateDiagFunc: validations.StringInSlice([]string{
								"BALANCED", "HIGH_SECURITY", "HIGH_COMPATIBILITY", "CUSTOM"}, false),
							Optional:    true,
							Description: "Network Loadbalancer Supported values are `BALANCED`,`HIGH_SECURITY`, `HIGH_COMPATIBILITY`,`CUSTOM`"},
						"cookie_mode": {
							Type:             schema.TypeString,
							ValidateDiagFunc: validations.StringInSlice([]string{"INSERT", "PREFIX", "REWRITE"}, false),
							Required:         true,
							Description:      "Network Loadbalancer Supported values are `INSERT`,`PREFIX`, `REWRITE`"},
						"cookie_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "cookie_name for Network Load balancer Profile",
						},
						"cookie_type": {
							Type:             schema.TypeString,
							ValidateDiagFunc: validations.StringInSlice([]string{"LBPersistenceCookieTime", "LBSessionCookieTime"}, false),
							Required:         true,
							Description:      "Network Loadbalancer Supported values are `LBPersistenceCookieTime`,`LBSessionCookieTime`"},
					},
				},
			},
		},
		ReadContext: LBProfileReadContext,
		Description: `The ` + DSLBProfile + ` data source can be used to discover the ID of a hpegl vmaas router.
		This can then be used with resources or data sources that require a ` + DSLBProfile + `,
		such as the ` + ResLoadBalancerProfiles + ` resource.`,
		SchemaVersion:  0,
		StateUpgraders: nil,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func LBProfileReadContext(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c, err := client.GetClientFromMetaMap(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	data := utils.NewData(d)
	err = c.CmpClient.LoadBalancerProfile.Read(ctx, data, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
