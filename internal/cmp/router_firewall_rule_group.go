// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package cmp

import (
	"context"
	"fmt"

	"github.com/HewlettPackard/hpegl-vmaas-cmp-go-sdk/pkg/client"
	"github.com/HewlettPackard/hpegl-vmaas-cmp-go-sdk/pkg/models"
	"github.com/HewlettPackard/hpegl-vmaas-terraform-resources/internal/utils"
	"github.com/tshihad/tftags"
)

type routerFirewallRuleGroup struct {
	rClient *client.RouterAPIService
}

func newRouterFirewallRuleGroup(rClient *client.RouterAPIService) *routerFirewallRuleGroup {
	return &routerFirewallRuleGroup{
		rClient: rClient,
	}
}

func (r *routerFirewallRuleGroup) Read(ctx context.Context, d *utils.Data, meta interface{}) error {
	setMeta(meta, r.rClient.Client)
	var tfModel models.CreateRouterFirewallRuleGroup
	if err := tftags.Get(d, &tfModel); err != nil {
		return err
	}

	_, err := r.rClient.GetSpecificRouterFirewallRuleGroup(ctx, tfModel.RouterID,
		tfModel.ID)
	if err != nil {
		return err
	}

	return tftags.Set(d, tfModel)
}

func (r *routerFirewallRuleGroup) Create(ctx context.Context, d *utils.Data, meta interface{}) error {
	setMeta(meta, r.rClient.Client)
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
	setMeta(meta, r.rClient.Client)
	var tfFirewallRuleGroup models.CreateRouterFirewallRuleGroup
	if err := tftags.Get(d, &tfFirewallRuleGroup); err != nil {
		return err
	}

	resp, err := r.rClient.DeleteRouterFirewallRuleGroup(ctx, tfFirewallRuleGroup.RouterID,
		tfFirewallRuleGroup.ID)
	if err != nil {
		return err
	}
	if !resp.Success {
		return fmt.Errorf("got success = 'false' while deleting firewall rule group rule")
	}

	return nil
}
