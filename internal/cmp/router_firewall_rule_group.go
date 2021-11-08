// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package cmp

import (
	"context"
	"fmt"
	"log"

	"github.com/HewlettPackard/hpegl-vmaas-cmp-go-sdk/pkg/client"
	"github.com/HewlettPackard/hpegl-vmaas-cmp-go-sdk/pkg/models"
	"github.com/HewlettPackard/hpegl-vmaas-terraform-resources/internal/utils"
	"github.com/tshihad/tftags"
)

type routerFirewallRuleGroup struct {
	routerFirewallRuleGroupClient *client.RouterAPIService
}

func newRouterFirewallRuleGroup(routerFirewallRuleGroupClient *client.RouterAPIService) *routerFirewallRuleGroup {
	return &routerFirewallRuleGroup{
		routerFirewallRuleGroupClient: routerFirewallRuleGroupClient,
	}
}

func (r *routerFirewallRuleGroup) Read(ctx context.Context, d *utils.Data, meta interface{}) error {
	var tfFirewallRuleGroup models.CreateRouterFirewallRuleGroup
	if err := tftags.Get(d, &tfFirewallRuleGroup); err != nil {
		return err
	}
	// Get the router, if the router not exists, return warning
	router, err := r.routerFirewallRuleGroupClient.GetSpecificRouter(ctx, tfFirewallRuleGroup.RouterID)
	if err != nil {
		return err
	}
	// if router not found set is_deprecated flag=true
	if router.ID == 0 {
		log.Printf("[ERROR] Router with %d id is not found on Firewall Rule Group plan", tfFirewallRuleGroup.RouterID)
		tfFirewallRuleGroup.IsDeprecated = true

		return tftags.Set(d, tfFirewallRuleGroup)
	}

	_, err = r.routerFirewallRuleGroupClient.GetSpecificRouterFirewallRuleGroup(ctx, tfFirewallRuleGroup.RouterID, tfFirewallRuleGroup.ID)
	if err != nil {
		return err
	}
	tfFirewallRuleGroup.IsDeprecated = false

	return tftags.Set(d, tfFirewallRuleGroup)
}

func (r *routerFirewallRuleGroup) Create(ctx context.Context, d *utils.Data, meta interface{}) error {
	var tfFirewallRuleGroup models.CreateRouterFirewallRuleGroup
	err := tftags.Get(d, &tfFirewallRuleGroup)
	if err != nil {
		return err
	}
	firewallGroupRes, err := r.routerFirewallRuleGroupClient.CreateRouterFirewallRuleGroup(ctx, tfFirewallRuleGroup.RouterID,
		models.CreateRouterFirewallRuleGroupRequest{CreateRouterFirewallRuleGroup: tfFirewallRuleGroup},
	)
	if err != nil {
		return err
	}

	if !firewallGroupRes.Success {
		return fmt.Errorf(successErr, "creating firewall rule group for the router")
	}
	tfFirewallRuleGroup.ID = firewallGroupRes.ID

	return tftags.Set(d, tfFirewallRuleGroup)
}

func (r *routerFirewallRuleGroup) Update(ctx context.Context, d *utils.Data, meta interface{}) error {
	// to be implemented
	return nil
}

func (r *routerFirewallRuleGroup) Delete(ctx context.Context, d *utils.Data, meta interface{}) error {
	var tfFirewallRuleGroup models.CreateRouterFirewallRuleGroup
	if err := tftags.Get(d, &tfFirewallRuleGroup); err != nil {
		return err
	}

	// if parent router got deleted, NAT is already deleted
	if tfFirewallRuleGroup.IsDeprecated {
		log.Printf("[WARNING] Firewall rule group already deleted since router is deleted")

		return nil
	}

	resp, err := r.routerFirewallRuleGroupClient.DeleteRouterFirewallRuleGroup(ctx, tfFirewallRuleGroup.RouterID, tfFirewallRuleGroup.ID)
	if err != nil {
		return err
	}
	if !resp.Success {
		return fmt.Errorf("got success = 'false' while deleting firewall rule group rule")
	}

	return nil
}
