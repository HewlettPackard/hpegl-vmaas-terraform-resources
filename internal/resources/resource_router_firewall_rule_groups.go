// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package resources

import (
	"context"

	"github.com/HewlettPackard/hpegl-vmaas-terraform-resources/internal/resources/validations"
	"github.com/HewlettPackard/hpegl-vmaas-terraform-resources/internal/utils"
	"github.com/HewlettPackard/hpegl-vmaas-terraform-resources/pkg/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func RouterFirewallRuleGroup() *schema.Resource {
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
				Description: "Name of the Firewall rule Group.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description for the Firewall rule Group.",
			},
			"priority": {
				Type:             schema.TypeInt,
				Optional:         true,
				Description:      "Firewall rule group priority",
				ValidateDiagFunc: validations.IntAtLeast(1),
			},
			"external_type": {
				Type:         schema.TypeString,
				Required:     true,
				InputDefault: "GatewayPolicy",
				Description:  "Platform/vendor specific type. Pass `GatewayPolicy`.",
			},
			"group_layer": {
				Type:     schema.TypeString,
				Required: true,
				ValidateDiagFunc: validations.StringInSlice([]string{
					"Emergency",
					"SharedPreRules",
					"LocalGatewayRules",
				}, false),
				Description: "Platform/vendor specific category",
			},
			"is_deprecated": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "If parent router not found, then is_deprecated will be true",
			},
		},
		ReadContext:   routerFirewallRuleGroupReadContext,
		CreateContext: routerFirewallRuleGroupCreateContext,
		UpdateContext: routerFirewallRuleGroupUpdateContext,
		DeleteContext: routerFirewallRuleGroupDeleteContext,
	}
}

func routerFirewallRuleGroupReadContext(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c, err := client.GetClientFromMetaMap(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	data := utils.NewData(rd)
	if err := c.CmpClient.RouterFirewallRuleGroup.Read(ctx, data, meta); err != nil {
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

func routerFirewallRuleGroupCreateContext(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c, err := client.GetClientFromMetaMap(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	data := utils.NewData(rd)
	if err := c.CmpClient.RouterFirewallRuleGroup.Create(ctx, data, meta); err != nil {
		return diag.FromErr(err)
	}

	return routerFirewallRuleGroupReadContext(ctx, rd, meta)
}

func routerFirewallRuleGroupUpdateContext(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c, err := client.GetClientFromMetaMap(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	data := utils.NewData(rd)
	if err := c.CmpClient.RouterFirewallRuleGroup.Update(ctx, data, meta); err != nil {
		return diag.FromErr(err)
	}

	return routerFirewallRuleGroupReadContext(ctx, rd, meta)
}

func routerFirewallRuleGroupDeleteContext(ctx context.Context, rd *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c, err := client.GetClientFromMetaMap(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	data := utils.NewData(rd)
	if err := c.CmpClient.RouterFirewallRuleGroup.Delete(ctx, data, meta); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
