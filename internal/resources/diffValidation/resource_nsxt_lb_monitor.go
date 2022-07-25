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
	err := l.validateHttp()
	if err != nil {
		return err
	}
	err = l.validateHttps()
	if err != nil {
		return err
	}
	err = l.validateIcmp()
	if err != nil {
		return err
	}
	err = l.validatePassive()
	if err != nil {
		return err
	}
	err = l.validateTcp()
	if err != nil {
		return err
	}
	err = l.validateUdp()
	if err != nil {
		return err
	}

	return nil
}

func (l *LoadBalancerMonitor) validateHttp() error {
	serviceType := "type"
	httpMonitor := "http_monitor"
	if l.diff.HasChange(serviceType) {
		service := l.diff.Get(serviceType)
		monitorType := l.diff.Get(httpMonitor)
		if service == "LBHttpMonitorProfile" {
			if len((monitorType).([]interface{})) == 0 {
				return fmt.Errorf("please provide http_monitor configurations for serviceType LBHttpMonitorProfile")
			}
		}
	}
	return nil
}

func (l *LoadBalancerMonitor) validateHttps() error {
	serviceType := "type"
	httpMonitor := "https_monitor"
	if l.diff.HasChange(serviceType) {
		service := l.diff.Get(serviceType)
		monitorType := l.diff.Get(httpMonitor)
		if service == "LBHttpsMonitorProfile" {
			if len((monitorType).([]interface{})) == 0 {
				return fmt.Errorf("please provide https_monitor configurations for Type LBHttpsMonitorProfile")
			}
		}
	}
	return nil
}

func (l *LoadBalancerMonitor) validateIcmp() error {
	serviceType := "type"
	httpMonitor := "icmp_monitor"
	if l.diff.HasChange(serviceType) {
		service := l.diff.Get(serviceType)
		monitorType := l.diff.Get(httpMonitor)
		if service == "LBIcmpMonitorProfile" {
			if len((monitorType).([]interface{})) == 0 {
				return fmt.Errorf("please provide icmp_monitor configurations for Type LBIcmpMonitorProfile")
			}
		}
	}
	return nil
}

func (l *LoadBalancerMonitor) validatePassive() error {
	serviceType := "type"
	httpMonitor := "passive_monitor"
	if l.diff.HasChange(serviceType) {
		service := l.diff.Get(serviceType)
		monitorType := l.diff.Get(httpMonitor)
		if service == "LBPassiveMonitorProfile" {
			if len((monitorType).([]interface{})) == 0 {
				return fmt.Errorf("please provide passive_monitor configurations for Type LBPassiveMonitorProfile")
			}
		}
	}
	return nil
}

func (l *LoadBalancerMonitor) validateTcp() error {
	serviceType := "type"
	httpMonitor := "tcp_monitor"
	if l.diff.HasChange(serviceType) {
		service := l.diff.Get(serviceType)
		monitorType := l.diff.Get(httpMonitor)
		if service == "LBTcpMonitorProfile" {
			if len((monitorType).([]interface{})) == 0 {
				return fmt.Errorf("please provide tcp_monitor configurations for Type LBTcpMonitorProfile")
			}
		}
	}
	return nil
}

func (l *LoadBalancerMonitor) validateUdp() error {
	serviceType := "type"
	httpMonitor := "udp_monitor"
	if l.diff.HasChange(serviceType) {
		service := l.diff.Get(serviceType)
		monitorType := l.diff.Get(httpMonitor)
		if service == "LBUdpMonitorProfile" {
			if len((monitorType).([]interface{})) == 0 {
				return fmt.Errorf("please provide udp_monitor configurations for Type LBUdpMonitorProfile")
			}
		}
	}
	return nil
}
