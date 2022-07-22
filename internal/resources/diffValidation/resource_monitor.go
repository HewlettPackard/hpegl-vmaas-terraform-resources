//  (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package diffvalidation

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type LoadBalancerMonitor struct {
	diff *schema.ResourceDiff
}

func NewLoadBalancerMonitorValidate(diff *schema.ResourceDiff) *LoadBalancerMonitor {
	return &LoadBalancerMonitor{
		diff: diff,
	}
}

func (l *LoadBalancerMonitor) DiffValidate() error {
	err := l.validateMonitorServiceTypes()
	if err != nil {
		return err
	}

	return nil
}

func (l *LoadBalancerMonitor) validateMonitorServiceTypes() error {
	serviceType := "type"
	if l.diff.HasChange(serviceType) {
		service := l.diff.Get(serviceType)
		switch service {
		case "LBHttpMonitorProfile":
			if l.diff.Get("http_monitor") != nil {
				return fmt.Errorf("please provide http_monitor configurations for serviceType LBHttpMonitorProfile")
			}
		case "LBHttpsMonitorProfile":
			if l.diff.Get("https_monitor") != nil {
				return fmt.Errorf("please provide https_monitor configurations for serviceType LBHttpsMonitorProfile")
			}
		case "LBIcmpMonitorProfile":
			if l.diff.Get("icmp_monitor") != nil {
				return fmt.Errorf("please provide icmp_monitor configurations for serviceType LBIcmpMonitorProfile")
			}
		case "LBPassiveMonitorProfile":
			if l.diff.Get("passive_monitor") != nil {
				return fmt.Errorf("please provide passive_monitor configurations for serviceType LBPassiveMonitorProfile")
			}
		case "LBTcpMonitorProfile":
			if l.diff.Get("tcp_monitor") != nil {
				return fmt.Errorf("please provide tcp_monitor configurations for serviceType LBTcpMonitorProfile")
			}
		case "LBUdpMonitorProfile":
			if l.diff.Get("udp_monitor") != nil {
				return fmt.Errorf("please provide udp_monitor configurations for serviceType LBUdpMonitorProfile")
			}
		}
	}
	return nil
}
