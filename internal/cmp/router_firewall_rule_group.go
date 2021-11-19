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
	rClient *client.RouterAPIService
}

func newRouterFirewallRuleGroup(routerFirewallRuleGroupClient *client.RouterAPIService) *routerFirewallRuleGroup {
	return &routerFirewallRuleGroup{
		rClient: routerFirewallRuleGroupClient,
	}
}

func (r *routerFirewallRuleGroup) Read(ctx context.Context, d *utils.Data, meta interface{}) error {
	var tfModel models.CreateRouterFirewallRuleGroup
	if err := tftags.Get(d, &tfModel); err != nil {
		return err
	}
	// Get the router, if the router not exists, return warning
	if check, err := checkRouterDeprecated(
		ctx, r.rClient, d, tfModel.RouterID, &tfModel.IsDeprecated, &tfModel,
	); err != nil || check {
		return err
	}

	_, err := r.rClient.GetSpecificRouterFirewallRuleGroup(ctx, tfModel.RouterID,
		tfModel.ID)
	if err != nil {
		return err
	}
	tfModel.IsDeprecated = false

	return tftags.Set(d, tfModel)
}

func (r *routerFirewallRuleGroup) Create(ctx context.Context, d *utils.Data, meta interface{}) error {
	var tfModel models.CreateRouterFirewallRuleGroup
	err := tftags.Get(d, &tfModel)
	if err != nil {
		return err
	}
	// Setting to Default value "GatewayPolicy"
	tfModel.ExternalType = routerFirewallExternalPolicy
	firewallGroupRes, err := r.rClient.CreateRouterFirewallRuleGroup(ctx, tfModel.RouterID,
		models.CreateRouterFirewallRuleGroupRequest{CreateRouterFirewallRuleGroup: tfModel},
	)
	if err != nil {
		return err
	}

	if !firewallGroupRes.Success {
		return fmt.Errorf(successErr, "creating firewall rule group for the router")
	}
	tfModel.ID = firewallGroupRes.ID

	return tftags.Set(d, tfModel)
}

func (r *routerFirewallRuleGroup) Update(ctx context.Context, d *utils.Data, meta interface{}) error {
	// to be implemented
	return nil
}

func (r *routerFirewallRuleGroup) Delete(ctx context.Context, d *utils.Data, meta interface{}) error {
	var tfModel models.CreateRouterFirewallRuleGroup
	if err := tftags.Get(d, &tfModel); err != nil {
		return err
	}

	// if parent router got deleted, NAT is already deleted
	if tfModel.IsDeprecated {
		log.Printf("[WARNING] Firewall rule group already deleted since router is deleted")

		return nil
	}

	resp, err := r.rClient.DeleteRouterFirewallRuleGroup(ctx, tfModel.RouterID,
		tfModel.ID)
	if err != nil {
		return err
	}
	if !resp.Success {
		return fmt.Errorf("got success = 'false' while deleting firewall rule group rule")
	}

	return nil
}
