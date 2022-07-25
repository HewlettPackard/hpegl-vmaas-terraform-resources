//  (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package diffvalidation

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	serviceTypes   = "type"
	httpMonitor    = "http_monitor"
	httpsMonitor   = "https_monitor"
	icmpMonitor    = "icmp_monitor"
	passiveMonitor = "passive_monitor"
	tcpMonitor     = "tcp_monitor"
	udpMonitor     = "udp_monitor"

	LBHttpMonitorProfile    = "LBHttpMonitorProfile"
	LBHttpsMonitorProfile   = "LBHttpsMonitorProfile"
	LBIcmpMonitorProfile    = "LBIcmpMonitorProfile"
	LBPassiveMonitorProfile = "LBPassiveMonitorProfile"
	LBTcpMonitorProfile     = "LBTcpMonitorProfile"
	LBUdpMonitorProfile     = "LBUdpMonitorProfile"
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
	err := l.validateMonitor()
	if err != nil {
		return err
	}
	return nil
}

func (l *LoadBalancerMonitor) validateMonitor() error {
	if l.diff.HasChange(serviceTypes) {
		service := l.diff.Get(serviceTypes)
		profileType := l.diff.Get(http_monitor)
		if service == LBHttpMonitorProfile {
			if len((monitorType).([]interface{})) == 0 {
				return fmt.Errorf("please provide http_monitor configurations for serviceType LBHttpMonitorProfile")
			}
		}
	}

	if l.diff.HasChange(serviceType) {
		service := l.diff.Get(serviceType)
		monitorType := l.diff.Get(https_monitor)
		if service == LBHttpsMonitorProfile {
			if len((monitorType).([]interface{})) == 0 {
				return fmt.Errorf("please provide https_monitor configurations for Type LBHttpsMonitorProfile")
			}
		}
	}

	if l.diff.HasChange(serviceType) {
		service := l.diff.Get(serviceType)
		monitorType := l.diff.Get(icmp_monitor)
		if service == LBIcmpMonitorProfile {
			if len((monitorType).([]interface{})) == 0 {
				return fmt.Errorf("please provide icmp_monitor configurations for Type LBIcmpMonitorProfile")
			}
		}
	}

	if l.diff.HasChange(serviceType) {
		service := l.diff.Get(serviceType)
		monitorType := l.diff.Get(passive_monitor)
		if service == LBPassiveMonitorProfile {
			if len((monitorType).([]interface{})) == 0 {
				return fmt.Errorf("please provide passive_monitor configurations for Type LBPassiveMonitorProfile")
			}
		}
	}

	if l.diff.HasChange(serviceType) {
		service := l.diff.Get(serviceType)
		monitorType := l.diff.Get(tcp_monitor)
		if service == LBTcpMonitorProfile {
			if len((monitorType).([]interface{})) == 0 {
				return fmt.Errorf("please provide tcp_monitor configurations for Type LBTcpMonitorProfile")
			}
		}
	}

	if l.diff.HasChange(serviceType) {
		service := l.diff.Get(serviceType)
		monitorType := l.diff.Get(udp_monitor)
		if service == LBUdpMonitorProfile {
			if len((monitorType).([]interface{})) == 0 {
				return fmt.Errorf("please provide udp_monitor configurations for Type LBUdpMonitorProfile")
			}
		}
	}
	return nil
}
