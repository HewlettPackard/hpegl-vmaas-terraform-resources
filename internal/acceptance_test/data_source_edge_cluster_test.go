// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package acceptancetest

import (
	"testing"

	api_client "github.com/HewlettPackard/hpegl-vmaas-cmp-go-sdk/pkg/client"
	"github.com/HewlettPackard/hpegl-vmaas-terraform-resources/internal/utils"
	"github.com/HewlettPackard/hpegl-vmaas-terraform-resources/pkg/atf"
)

func TestAccDataSourceNetworkEdgeCluster(t *testing.T) {
	acc := &atf.Acc{
		PreCheck:     testAccPreCheck,
		Providers:    testAccProviders,
		ResourceName: "hpegl_vmaas_edge_cluster",
		GetAPI: func(attr map[string]string) (interface{}, error) {
			cl, cfg := getAPIClient()
			iClient := api_client.RouterAPIService{
				Client: cl,
				Cfg:    cfg,
			}
			name := attr["name"]
			var serverID int
			ServerResp, err := iClient.GetNetworkServices(getAccContext(), nil)
			if err != nil {
				return nil, err
			}
			// Align request for Network Server
			if len(ServerResp.NetworkServices) == 0 {
				return nil, err
			}
			cmpVersion, err := utils.GetCmpVersion(iClient.Client)
			if err != nil {
				return nil, err
			}
			nsxVar := "NSX-T"
			if v, _ := utils.ParseVersion("6.2.4"); v <= cmpVersion {
				// from 6.2.4 onwards the display name of NSX-T has been change to NSX
				nsxVar = "NSX"
			}
			for i, n := range ServerResp.NetworkServices {
				if n.TypeName == nsxVar {
					serverID = ServerResp.NetworkServices[i].ID

					break
				}
			}

			return iClient.GetEdgeCluster(getAccContext(), serverID, name)
		},
	}

	acc.RunDataSourceTests(t)
}
