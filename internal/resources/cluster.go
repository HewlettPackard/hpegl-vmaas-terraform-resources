// (C) Copyright 2020 Hewlett Packard Enterprise Development LP

package resources

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/hpe-hcss/hpecli-generated-caas-client/pkg/mcaasapi"

	"github.com/hpe-hcss/poc-caas-terraform-resources/pkg/client"
)

const (
	stateError   = "error"
	stateReady   = "ready"
	stateDeleted = "deleted"

	clusterAvailableTimeout = 60 * time.Minute
	clusterDeleteTimeout    = 60 * time.Minute
	pollingInterval         = 3 * time.Second

	errClusterPollForStateFmt = "Error in poll for cluster state %s"
)

func Cluster() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"blueprint_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"appliance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"space_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
		SchemaVersion:  0,
		StateUpgraders: nil,
		CreateContext:  clusterCreateContext,
		ReadContext:    clusterReadContext,
		// TODO figure out if a cluster can be updated
		// Update:             clusterUpdate,
		DeleteContext: clusterDeleteContext,
		CustomizeDiff: nil,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		DeprecationMessage: "",
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(clusterAvailableTimeout),
			// Update: schema.DefaultTimeout(clusterAvailableTimeout),
			Delete: schema.DefaultTimeout(clusterDeleteTimeout),
		},
		Description: "",
	}
}

func clusterCreateContext(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := client.GetClientFromMetaMap(meta)
	token := c.IAMToken
	clientCtx := context.WithValue(ctx, mcaasapi.ContextAccessToken, token)

	var diags diag.Diagnostics

	spaceID := d.Get("space_id").(string)

	createCluster := mcaasapi.CreateCluster{
		Name:               d.Get("name").(string),
		ClusterBlueprintId: d.Get("blueprint_id").(string),
		ApplianceID:        d.Get("appliance_id").(string),
		SpaceID:            spaceID,
	}

	cluster, _, err := c.CaasClient.ClusterAdminApi.ClustersPost(clientCtx, createCluster)
	if err != nil {
		diags = append(diags, diag.Errorf("Error in ClustersPost: %s", err)...)

		return diags
	}

	// TODO Should we be passing clientCtx here?
	if ds := pollForClusterState(ctx, cluster.Id, spaceID, stateReady, meta); ds.HasError() {
		diags = append(diags, ds...)
		diags = append(diags, diag.Errorf(errClusterPollForStateFmt, stateReady)...)

		return diags
	}

	// Only set id to non-empty string if resource has been successfully created
	d.SetId(cluster.Id)

	// TODO Should we be passing clientCtx here?
	return clusterReadContext(ctx, d, meta)
}

func clusterReadContext(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := client.GetClientFromMetaMap(meta)
	token := c.IAMToken
	clientCtx := context.WithValue(ctx, mcaasapi.ContextAccessToken, token)

	var diags diag.Diagnostics
	id := d.Id()
	spaceID := d.Get("space_id").(string)

	cluster, _, err := c.CaasClient.ClusterAdminApi.ClustersIdGet(clientCtx, id, spaceID)
	if err != nil {
		return diag.FromErr(err)
	}

	if err = d.Set("state", cluster.State); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func clusterDeleteContext(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := client.GetClientFromMetaMap(meta)
	token := c.IAMToken
	clientCtx := context.WithValue(ctx, mcaasapi.ContextAccessToken, token)

	var diags diag.Diagnostics
	id := d.Id()
	spaceID := d.Get("space_id").(string)

	_, err := c.CaasClient.ClusterAdminApi.ClustersIdDelete(clientCtx, id)
	if err != nil {
		return diag.FromErr(err)
	}

	// TODO Should we be passing clientCtx here?
	// TODO This needs to be fixed.  Unfortunately we get an "unauthorized" when trying to GET a deleted cluster,
	// there is no "deleted" state.
	if ds := pollForClusterState(ctx, id, spaceID, stateDeleted, meta); ds.HasError() {
		diags = append(diags, ds...)
		diags = append(diags, diag.Errorf(errClusterPollForStateFmt, stateDeleted)...)

		return diags
	}

	// Only set id to "" if delete has been successful, this means that terraform will delete the resource entry
	// This also means that the destroy can be reattempted by terraform if there was an error
	d.SetId("")

	return diags
}

func pollForClusterState(
	ctx context.Context,
	id, spaceID, expectedState string,
	meta interface{},
) diag.Diagnostics {
	c := client.GetClientFromMetaMap(meta)
	token := c.IAMToken
	clientCtx := context.WithValue(ctx, mcaasapi.ContextAccessToken, token)

	// Set up ticker for the loop below
	ticker := time.NewTicker(pollingInterval)

	for {
		_, ok := <-ticker.C
		if !ok {
			ticker.Stop()

			return diag.Errorf("error in getting ticker")
		}

		cluster, _, err := c.CaasClient.ClusterAdminApi.ClustersIdGet(clientCtx, id, spaceID)
		if err != nil {
			ticker.Stop()

			// TODO we might not want to expose the cluster ID
			return diag.Errorf("Error in poll of cluster id %s : %s", id, err)
		}

		switch cluster.State {
		case expectedState:
			ticker.Stop()

			return nil
		case stateError:
			ticker.Stop()

			return diag.Errorf("cluster in error state")
		}
	}
}
