//  (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package diffvalidation

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	dhcpNetwork   = "dhcp_network"
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
	if isEnabled == false {
		err := l.validateNetworkConfigs(dhcpNetwork)
		if err != nil {
			return err
		}
	}

	return nil
}

func (l *Network) validateNetworkConfigs(networkTypes string) error {
	value := l.diff.Get(networkTypes)
	if len((value).([]interface{})) != 0 {
		return fmt.Errorf("Do not provide the DHCP configurations when the" + isDhcpEnabled + "is disabled")
	}

	return nil
}
