// (C) Copyright 2022 Hewlett Packard Enterprise Development LP

package resources

import (
	"context"

	diffvalidation "github.com/HewlettPackard/hpegl-vmaas-terraform-resources/internal/resources/diffValidation"
	"github.com/HewlettPackard/hpegl-vmaas-terraform-resources/internal/resources/schemas"
	"github.com/HewlettPackard/hpegl-vmaas-terraform-resources/internal/resources/validations"
	"github.com/HewlettPackard/hpegl-vmaas-terraform-resources/internal/utils"
	"github.com/HewlettPackard/hpegl-vmaas-terraform-resources/pkg/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func LoadBalancerProfiles() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"lb_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Parent lb ID, lb_id can be obtained by using LB datasource/resource.",
				ForceNew:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Network loadbalancer Profile Name",
			},
			"description": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Creating the Network Load balancer Profile",
			},
			"service_type": {
				Type:             schema.TypeString,
				ValidateDiagFunc: validations.StringInSlice([]string{"LBHttpProfile", "LBFastTcpProfile", "LBFastUdpProfile", "LBClientSslProfile", "LBServerSslProfile", "LBCookiePersistenceProfile", "LBGenericPersistenceProfile", "LBSourceIpPersistenceProfile"}, false),
				Required:         true,
				Description:      "Network Loadbalancer Supported values are `LBHttpProfile`,`LBFastTcpProfile`, `LBFastUdpProfile`, `LBClientSslProfile`,`LBServerSslProfile`, `LBCookiePersistenceProfile`,`LBGenericPersistenceProfile`,`LBSourceIpPersistenceProfile`"},
			"config": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "profile Configuration",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"http_profile":     schemas.HttpProfileSchema(),
						"tcp_profile":      schemas.TcpProfileSchema(),
						"udp_profile":      schemas.UdpProfileSchema(),
						"cookie_profile":   schemas.CookieProfileSchema(),
						"sourceip_profile": schemas.SourceIPProfileSchema(),
						"generic_profile":  schemas.GenericProfileSchema(),
						"client_profile":   schemas.ClientProfileSchema(),
						"server_profile":   schemas.ServerProfileSchema(),
						"tags": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "tags Configuration",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"tag": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "tag for Network Load balancer Profile",
									},
									"scope": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "scope for Network Load balancer Profile",
									},
								},
							},
						},
					},
				},
			},
		},
		ReadContext:   loadbalancerProfileReadContext,
		UpdateContext: loadbalancerProfileUpdateContext,
		CreateContext: loadbalancerProfileCreateContext,
		DeleteContext: loadbalancerProfileDeleteContext,
		CustomizeDiff: profileCustomDiff,
		Description: `loadbalancer Profile resource facilitates creating,
		and deleting NSX-T  Network Load Balancers.`,
	}
}

func loadbalancerProfileUpdateContext(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c, err := client.GetClientFromMetaMap(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	data := utils.NewData(rd)
	if err := c.CmpClient.LoadBalancerProfile.Update(ctx, data, meta); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func loadbalancerProfileReadContext(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c, err := client.GetClientFromMetaMap(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	data := utils.NewData(rd)
	if err := c.CmpClient.LoadBalancerProfile.Read(ctx, data, meta); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func loadbalancerProfileCreateContext(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c, err := client.GetClientFromMetaMap(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	data := utils.NewData(rd)
	if err := c.CmpClient.LoadBalancerProfile.Create(ctx, data, meta); err != nil {
		return diag.FromErr(err)
	}

	return loadbalancerProfileReadContext(ctx, rd, meta)
}

func loadbalancerProfileDeleteContext(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c, err := client.GetClientFromMetaMap(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	data := utils.NewData(rd)
	if err := c.CmpClient.LoadBalancerProfile.Delete(ctx, data, meta); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func profileCustomDiff(ctx context.Context, diff *schema.ResourceDiff, meta interface{}) error {
	return diffvalidation.NewLoadBalancerProfileValidate(diff).DiffValidate()
}
