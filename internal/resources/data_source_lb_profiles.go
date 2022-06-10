// (C) Copyright 2022 Hewlett Packard Enterprise Development LP

package resources

import (
	"context"

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
			"serviceType": {
				Type: schema.TypeString,
				ValidateDiagFunc: validations.StringInSlice([]string{
					"LBHttpProfile", "LBFastTcpProfile", "LBFastUdpProfile",
					 "LBClientSslProfile","LBServerSslProfile",
					 "LBCookiePersistenceProfile","LBGenericPersistenceProfile"
				}, false),
				Computed:    true,
				Description: "Network Loadbalancer Supported values are `LBHttpProfile`, 
				`LBFastTcpProfile`, `LBFastUdpProfile`, `LBClientSslProfile`,
				 `LBServerSslProfile`, `LBCookiePersistenceProfile`,`LBGenericPersistenceProfile`",
			},
			"config": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "profile Configuration",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"profileType": {
							Type: schema.TypeString,
				            ValidateDiagFunc: validations.StringInSlice([]string{
					            "application-profile", "ssl-profile", "persistence-profile",
				            }, false),
				            Computed:    true,
				            Description: "Network Loadbalancer Supported values are `application-profile`, 
				            `ssl-profile`, `persistence-profile`",
						},
						"sslSuite": {
							Type: schema.TypeString,
				            ValidateDiagFunc: validations.StringInSlice([]string{
				        	    "BALANCED", "HIGH_SECURITY", "HIGH_COMPATIBILITY", "CUSTOM"
				             }, false),
				            Computed:    true,
				            Description: "Network Loadbalancer Supported values are `BALANCED`, 
				                `HIGH_SECURITY`, `HIGH_COMPATIBILITY`,`CUSTOM`",
						},
						"cookieMode": {
							Type: schema.TypeString,
				            ValidateDiagFunc: validations.StringInSlice([]string{
					            "INSERT", "PREFIX", "REWRITE"
				            }, false),
				            Computed:    true,
				            Description: "Network Loadbalancer Supported values are `INSERT`, 
				                `PREFIX`, `REWRITE`",
						},
						"cookieName": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "cookieName for Network Load balancer Profile",
						},
						"cookieType": {
							Type: schema.TypeString,
				            ValidateDiagFunc: validations.StringInSlice([]string{
					            "LBPersistenceCookieTime", "LBSessionCookieTime"
				            }, false),
				            Computed:    true,
				            Description: "Network Loadbalancer Supported values are `LBPersistenceCookieTime`, 
				            `LBSessionCookieTime`",
						},
					},
				},
			},
		},
		ReadContext:   LBProfileReadContext,
		Description: `The ` + DSLoadBalancer + ` profile data source can be used to discover the ID of a hpegl vmaas router.
		This can then be used with resources or data sources that require a ` + DSLoadBalancer + `,
		profile such as the ` + ResLoadBalancer + ` profile resource.`,
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
	err = c.CmpClient.DSLBProfile.Read(ctx, data, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
