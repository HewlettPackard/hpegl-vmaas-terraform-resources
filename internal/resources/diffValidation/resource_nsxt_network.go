//  (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package diffvalidation

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	dhcpNetwork   = "dhcp_network"
	staticNetwork = "static_network"
	dhcpConfig    = "dhcp_config"
	staticConfig  = "static_config"
)

type Network struct {
	diff *schema.ResourceDiff
}

func NewNetworkValidate(diff *schema.ResourceDiff) *Network {
	return &Network{
		diff: diff,
	}
}

func (l *Network) DiffValidate() error {
	err := l.validateNetworks()
	if err != nil {
		return err
	}

	return nil
}

func (l *Network) validateNetworks() error {
	err := l.validateNetworkConfigs(dhcpNetwork, staticNetwork)
	if err != nil {
		return err
	}

	return nil
}

func (l *Network) validateNetworkConfigs(dhcpTypes interface{}, staticTypes interface{}) error {
	if len((dhcpTypes).([]interface{})) != 0 {
		for _, networks := range dhcpTypes.([]interface{}) {
			if len(networks.([]interface{})) == 0 {
				return fmt.Errorf("please provide " + dhcpConfig + " for the Configuration")
			}
		}
	}
	if len((staticTypes).([]interface{})) != 0 {
		for _, networks := range staticTypes.([]interface{}) {
			if len(networks.([]interface{})) == 0 {
				return fmt.Errorf("please provide " + staticConfig + " for the Configuration")
			}
		}
	}

	return nil
}
