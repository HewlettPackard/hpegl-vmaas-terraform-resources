//  (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package diffvalidation

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type LoadBalancerProfile struct {
	diff *schema.ResourceDiff
}

func NewLoadBalancerProfileValidate(diff *schema.ResourceDiff) *LoadBalancerProfile {
	return &LoadBalancerProfile{
		diff: diff,
	}
}

func (l *LoadBalancerProfile) DiffValidate() error {
	err := l.validateProfileServiceTypes()
	if err != nil {
		return err
	}

	return nil
}

func (l *LoadBalancerProfile) validateProfileServiceTypes() error {
	serviceType := "config.0.http_profile.0.profile_type"
	if l.diff.HasChange(serviceType) {
		service := l.diff.Get(serviceType)
		switch service {
		case "application-profile":
			if l.diff.Get("service_type") != "LBHttpProfile" {
				return fmt.Errorf("please provide http_monitor configurations for serviceType LBHttpMonitorProfile")
			}

			// actionPath := "config.0.action"
			// if r.diff.HasChange(actionPath) {
			// 	action := r.diff.Get(actionPath)
			// 	switch action {
			// 	case "DNAT":
			// 		if r.diff.Get("destination_network") == "" {
			// 			return fmt.Errorf("destination_network should be set for DNAT")
			// 		}
			// 	case "SNAT":
			// 		if r.diff.Get("source_network") == "" {
			// 			return fmt.Errorf("source_network should be set for SNAT")
			// 		}
			// 	}
			// }
			// case "LBFastTcpProfile":
			// 	if l.diff.Get("tcp_profile") != nil {
			// 		return fmt.Errorf("please provide https_monitor configurations for serviceType LBHttpsMonitorProfile")
			// 	}
			// case "LBFastUdpProfile":
			// 	if l.diff.Get("udp_profile") != nil {
			// 		return fmt.Errorf("please provide icmp_monitor configurations for serviceType LBIcmpMonitorProfile")
			// 	}
			// case "LBCookiePersistenceProfile":
			// 	if l.diff.Get("cookie_profile") != nil {
			// 		return fmt.Errorf("please provide passive_monitor configurations for serviceType LBPassiveMonitorProfile")
			// 	}
			// case "LBGenericPersistenceProfile":
			// 	if l.diff.Get("generic_profile") != nil {
			// 		return fmt.Errorf("please provide tcp_monitor configurations for serviceType LBTcpMonitorProfile")
			// 	}
			// case "LBSourceIpPersistenceProfile":
			// 	if l.diff.Get("sourceip_profile") != nil {
			// 		return fmt.Errorf("please provide udp_monitor configurations for serviceType LBUdpMonitorProfile")
			// 	}
			// case "LBClientSslProfile":
			// 	if l.diff.Get("client_profile") != nil {
			// 		return fmt.Errorf("please provide udp_monitor configurations for serviceType LBUdpMonitorProfile")
			// 	}
			// case "LBServerSslProfile":
			// 	if l.diff.Get("server_profile") != nil {
			// 		return fmt.Errorf("please provide udp_monitor configurations for serviceType LBUdpMonitorProfile")
			// 	}
		}
	}
	return nil
}
