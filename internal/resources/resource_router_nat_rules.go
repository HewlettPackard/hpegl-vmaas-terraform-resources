// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package resources

import (
	"context"

	diffvalidation "github.com/HewlettPackard/hpegl-vmaas-terraform-resources/internal/resources/diffValidation"
	"github.com/HewlettPackard/hpegl-vmaas-terraform-resources/internal/resources/validations"
	"github.com/HewlettPackard/hpegl-vmaas-terraform-resources/internal/utils"
	"github.com/HewlettPackard/hpegl-vmaas-terraform-resources/pkg/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func RouterNatRule() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"router_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Parent router ID, router_id can be obtained by using router datasource/resource.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the NAT rule.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description for the NAT rule.",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Default:     false,
				Optional:    true,
				Description: "If true then NAT rule will be active/enabled.",
			},
			"config": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Required:    true,
				Description: "NAT configurations",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"action": {
							Type: schema.TypeString,
							ValidateDiagFunc: validations.StringInSlice([]string{
								"DNAT", "SNAT",
							}, false),
							Required:    true,
							Description: "Supported values are DNAT and SNAT",
						},
						"service": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Type of the service",
						},
						"firewall": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "MATCH_INTERNAL_ADDRESS",
							ValidateDiagFunc: validations.StringInSlice([]string{
								"MATCH_EXTERNAL_ADDRESS", "MATCH_INTERNAL_ADDRESS", "BYPASS",
							}, false),
							// "MATCH_INTERNAL_ADDRESS",
						},
						// This field will added on later versions
						// "scope": {
						// 	Type:        schema.TypeString,
						// 	Optional:    true,
						// 	Description: "Scope to particular router interface",
						// },
						"logging": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},
			"source_network": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: validations.ValidateCidr,
				Description:      "Source Network CIDR Address",
			},
			"destination_network": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: validations.ValidateCidr,
				Description:      "Destination Network CIDR Address",
			},
			"translated_network": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validations.ValidateCidr,
				Description:      "Translated Network CIDR Address",
			},
			"translated_ports": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Translated Network Port",
			},
			"priority": {
				Type:             schema.TypeInt,
				Optional:         true,
				Default:          100,
				Description:      "Priority for the rule",
				ValidateDiagFunc: validations.IntAtLeast(1),
			},
			"is_deprecated": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "If parent router not found, then is_deprecated will be true",
			},
		},
		ReadContext:   routerNatRuleReadContext,
		CreateContext: routerNatRuleCreateContext,
		UpdateContext: routerNatRuleUpdateContext,
		DeleteContext: routerNatRuleDeleteContext,
		CustomizeDiff: routerNatCustomDiff,
	}
}

func routerNatRuleReadContext(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c, err := client.GetClientFromMetaMap(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	data := utils.NewData(rd)
	if err := c.CmpClient.RouterNat.Read(ctx, data, meta); err != nil {
		return diag.FromErr(err)
	}
	isDeprecated := data.GetBool("is_deprecated")
	if isDeprecated {
		return diag.Diagnostics{
			{
				Severity: diag.Warning,
				Summary:  "Parent router is deleted. This resource is deprecated!!!",
			},
		}
	}
	return nil
}

func routerNatRuleCreateContext(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c, err := client.GetClientFromMetaMap(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	data := utils.NewData(rd)
	if err := c.CmpClient.RouterNat.Create(ctx, data, meta); err != nil {
		return diag.FromErr(err)
	}

	return routerNatRuleReadContext(ctx, rd, meta)
}

func routerNatRuleUpdateContext(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c, err := client.GetClientFromMetaMap(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	data := utils.NewData(rd)
	if err := c.CmpClient.RouterNat.Update(ctx, data, meta); err != nil {
		return diag.FromErr(err)
	}

	return routerNatRuleReadContext(ctx, rd, meta)
}

func routerNatRuleDeleteContext(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c, err := client.GetClientFromMetaMap(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	data := utils.NewData(rd)
	if err := c.CmpClient.RouterNat.Delete(ctx, data, meta); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func routerNatCustomDiff(ctx context.Context, diff *schema.ResourceDiff, meta interface{}) error {
	return diffvalidation.NewRouterNatValidate(diff).DiffValidate()
}
