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
	return nil
}

func (l *LoadBalancerMonitor) validateHttp() error {
	serviceType := "http_monitor.0.request_version"
	if l.diff.HasChange((serviceType)) {
		service := l.diff.Get(serviceType)
		switch service {
		case "HTTP_VERSION_1_0":
			if l.diff.Get("type") != "LBHttpMonitorProfile" || l.diff.Get("type") == "" {
				return fmt.Errorf("please provide type as LBHttpMonitorProfile for http_monitor")
			}
		case "HTTP_VERSION_1_1":
			if l.diff.Get("type") != "LBHttpMonitorProfile" || l.diff.Get("type") == "" {
				return fmt.Errorf("please provide type as LBHttpMonitorProfile for http_monitor")
			}
		}
	}
	return nil
}

func (l *LoadBalancerMonitor) validateHttps() error {
	serviceType := "https_monitor.0.request_version"
	if l.diff.HasChange((serviceType)) {
		service := l.diff.Get(serviceType)
		switch service {
		case "HTTP_VERSION_1_0":
			if l.diff.Get("type") != "LBHttpsMonitorProfile" || l.diff.Get("type") == "" {
				return fmt.Errorf("please provide type as LBHttpMonitorProfile for http_monitor")
			}
		case "HTTP_VERSION_1_1":
			if l.diff.Get("type") != "LBHttpsMonitorProfile" || l.diff.Get("type") == "" {
				return fmt.Errorf("please provide type as LBHttpMonitorProfile for http_monitor")
			}
		}
	}
	return nil
}
