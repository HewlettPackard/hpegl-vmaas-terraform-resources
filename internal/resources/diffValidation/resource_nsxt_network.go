//  (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package diffvalidation

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	dhcpNetwork   = "dhcp_network"
	staticNetwork = "static_network"
	isDhcpEnabled = "dhcp_enabled"
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
	isEnabled := l.diff.Get(isDhcpEnabled)
	if isEnabled == true {
		err := l.validateNetworkConfigs(dhcpNetwork)
		if err != nil {
			return err
		}
	} else if isEnabled == false {
		err := l.validateNetworkConfigs(staticNetwork)
		if err != nil {
			return err
		}
	}

	return nil
}

func (l *Network) validateNetworkConfigs(networkTypes interface{}) error {
	if len((networkTypes).([]interface{})) != 0 {
		return fmt.Errorf("please provide " + networkTypes.(string) + " for the Configuration")
	}
	return nil
}
