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
	types := l.diff.Get(serviceTypes)
	switch types {
	case LBHttpMonitorProfile:
		err := l.validateMonitor(httpMonitor, LBHttpMonitorProfile)
		if err != nil {
			return err
		}
	case LBHttpsMonitorProfile:
		err := l.validateMonitor(httpsMonitor, LBHttpsMonitorProfile)
		if err != nil {
			return err
		}
	case LBIcmpMonitorProfile:
		err := l.validateMonitor(icmpMonitor, LBIcmpMonitorProfile)
		if err != nil {
			return err
		}
	case LBPassiveMonitorProfile:
		err := l.validateMonitor(passiveMonitor, LBPassiveMonitorProfile)
		if err != nil {
			return err
		}
	case LBTcpMonitorProfile:
		err := l.validateMonitor(tcpMonitor, LBTcpMonitorProfile)
		if err != nil {
			return err
		}
	case LBUdpMonitorProfile:
		err := l.validateMonitor(udpMonitor, LBUdpMonitorProfile)
		if err != nil {
			return err
		}
	}
	return nil
}

func (l *LoadBalancerMonitor) validateMonitor(monitor_type string, service_type string) error {
	monitorType := l.diff.Get(monitor_type)
	if len((monitorType).([]interface{})) == 0 {
		return fmt.Errorf("please provide " + monitor_type + " " + "configurations for Type" + " " + service_type)
	}
	return nil
}
