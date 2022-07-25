//  (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package diffvalidation

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	serviceTypes    = "service_type"
	httpProfile     = "http_profile"
	tcpProfile      = "tcp_profile"
	udpProfile      = "udp_profile"
	cookieProfile   = "cookie_profile"
	sourceIpProfile = "sourceip_profile"
	genericProfile  = "generic_profile"
	clientProfile   = "client_profile"
	serverProfile   = "server_profile"
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
	err := l.validateHttpProfile()
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
	err = l.validateCookie()
	if err != nil {
		return err
	}
	err = l.validateGeneric()
	if err != nil {
		return err
	}
	err = l.validateSourceIP()
	if err != nil {
		return err
	}

	err = l.validateClient()
	if err != nil {
		return err
	}
	err = l.validateServer()
	if err != nil {
		return err
	}
	return nil
}

func (l *LoadBalancerProfile) validateHttpProfile() error {
	serviceType := serviceTypes
	http_profile := httpProfile
	if l.diff.HasChange(serviceType) {
		service := l.diff.Get(serviceType)
		profileType := l.diff.Get(http_profile)
		if service == "LBHttpProfile" {
			if len((profileType).([]interface{})) == 0 {
				return fmt.Errorf("please provide http_profile configurations for serviceType LBHttpProfile")
			}
		}
	}
	return nil
}

func (l *LoadBalancerProfile) validateTcp() error {
	serviceType := serviceTypes
	tcp_profile := tcpProfile
	if l.diff.HasChange(serviceType) {
		service := l.diff.Get(serviceType)
		profileType := l.diff.Get(tcp_profile)
		if service == "LBFastTcpProfile" {
			if len((profileType).([]interface{})) == 0 {
				return fmt.Errorf("please provide tcp_profile configurations for serviceType LBFastTcpProfile")
			}
		}
	}
	return nil
}

func (l *LoadBalancerProfile) validateUdp() error {
	serviceType := serviceTypes
	udp_profile := udpProfile
	if l.diff.HasChange(serviceType) {
		service := l.diff.Get(serviceType)
		profileType := l.diff.Get(udp_profile)
		if service == "LBFastUdpProfile" {
			if len((profileType).([]interface{})) == 0 {
				return fmt.Errorf("please provide udp_profile configurations for serviceType LBFastUdpProfile")
			}
		}
	}
	return nil
}

func (l *LoadBalancerProfile) validateCookie() error {
	serviceType := serviceTypes
	cookie_profile := cookieProfile
	if l.diff.HasChange(serviceType) {
		service := l.diff.Get(serviceType)
		profileType := l.diff.Get(cookie_profile)
		if service == "LBCookiePersistenceProfile" {
			if len((profileType).([]interface{})) == 0 {
				return fmt.Errorf("please provide cookie_profile configurations for serviceType LBCookiePersistenceProfile")
			}
		}
	}
	return nil
}

func (l *LoadBalancerProfile) validateSourceIP() error {
	serviceType := serviceTypes
	source_profile := sourceIpProfile
	if l.diff.HasChange(serviceType) {
		service := l.diff.Get(serviceType)
		profileType := l.diff.Get(source_profile)
		if service == "LBSourceIpPersistenceProfile" {
			if len((profileType).([]interface{})) == 0 {
				return fmt.Errorf("please provide sourceip_profile configurations for serviceType LBSourceIpPersistenceProfile")
			}
		}
	}
	return nil
}

func (l *LoadBalancerProfile) validateGeneric() error {
	serviceType := serviceTypes
	generic_profile := genericProfile
	if l.diff.HasChange(serviceType) {
		service := l.diff.Get(serviceType)
		profileType := l.diff.Get(generic_profile)
		if service == "LBGenericPersistenceProfile" {
			if len((profileType).([]interface{})) == 0 {
				return fmt.Errorf("please provide generic_profile configurations for serviceType LBGenericPersistenceProfile")
			}
		}
	}
	return nil
}

func (l *LoadBalancerProfile) validateClient() error {
	serviceType := serviceTypes
	client_profile := clientProfile
	if l.diff.HasChange(serviceType) {
		service := l.diff.Get(serviceType)
		profileType := l.diff.Get(client_profile)
		if service == "LBClientSslProfile" {
			if len((profileType).([]interface{})) == 0 {
				return fmt.Errorf("please provide client_profile configurations for serviceType LBClientSslProfile")
			}
		}
	}
	return nil
}

func (l *LoadBalancerProfile) validateServer() error {
	serviceType := serviceTypes
	server_profile := serverProfile
	if l.diff.HasChange(serviceType) {
		service := l.diff.Get(serviceType)
		profileType := l.diff.Get(server_profile)
		if service == "LBServerSslProfile" {
			if len((profileType).([]interface{})) == 0 {
				return fmt.Errorf("please provide server_profile configurations for serviceType LBServerSslProfile")
			}
		}
	}
	return nil
}
