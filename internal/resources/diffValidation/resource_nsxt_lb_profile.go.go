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

	LBHttpProfile                = "LBHttpProfile"
	LBFastTcpProfile             = "LBFastTcpProfile"
	LBFastUdpProfile             = "LBFastUdpProfile"
	LBCookiePersistenceProfile   = "LBCookiePersistenceProfile"
	LBSourceIpPersistenceProfile = "LBSourceIpPersistenceProfile"
	LBGenericPersistenceProfile  = "LBGenericPersistenceProfile"
	LBClientSslProfile           = "LBClientSslProfile"
	LBServerSslProfile           = "LBServerSslProfile"
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
	err := l.validateProfile()
	if err != nil {
		return err
	}
	return nil
}

func (l *LoadBalancerProfile) validateProfile() error {
	if l.diff.HasChange(serviceTypes) {
		service := l.diff.Get(serviceTypes)
		profileType := l.diff.Get(httpProfile)
		if service == LBHttpProfile {
			if len((profileType).([]interface{})) == 0 {
				return fmt.Errorf("please provide http_profile configurations for serviceType LBHttpProfile")
			}
		}
	}

	if l.diff.HasChange(serviceTypes) {
		service := l.diff.Get(serviceTypes)
		profileType := l.diff.Get(tcpProfile)
		if service == LBFastTcpProfile {
			if len((profileType).([]interface{})) == 0 {
				return fmt.Errorf("please provide tcp_profile configurations for serviceType LBFastTcpProfile")
			}
		}
	}

	if l.diff.HasChange(serviceTypes) {
		service := l.diff.Get(serviceTypes)
		profileType := l.diff.Get(udpProfile)
		if service == LBFastUdpProfile {
			if len((profileType).([]interface{})) == 0 {
				return fmt.Errorf("please provide udp_profile configurations for serviceType LBFastUdpProfile")
			}
		}
	}

	if l.diff.HasChange(serviceTypes) {
		service := l.diff.Get(serviceTypes)
		profileType := l.diff.Get(cookieProfile)
		if service == LBCookiePersistenceProfile {
			if len((profileType).([]interface{})) == 0 {
				return fmt.Errorf("please provide cookie_profile configurations for serviceType LBCookiePersistenceProfile")
			}
		}
	}

	if l.diff.HasChange(serviceTypes) {
		service := l.diff.Get(serviceTypes)
		profileType := l.diff.Get(sourceIpProfile)
		if service == LBSourceIpPersistenceProfile {
			if len((profileType).([]interface{})) == 0 {
				return fmt.Errorf("please provide sourceip_profile configurations for serviceType LBSourceIpPersistenceProfile")
			}
		}
	}

	if l.diff.HasChange(serviceTypes) {
		service := l.diff.Get(serviceTypes)
		profileType := l.diff.Get(genericProfile)
		if service == LBGenericPersistenceProfile {
			if len((profileType).([]interface{})) == 0 {
				return fmt.Errorf("please provide generic_profile configurations for serviceType LBGenericPersistenceProfile")
			}
		}
	}

	if l.diff.HasChange(serviceTypes) {
		service := l.diff.Get(serviceTypes)
		profileType := l.diff.Get(clientProfile)
		if service == LBClientSslProfile {
			if len((profileType).([]interface{})) == 0 {
				return fmt.Errorf("please provide client_profile configurations for serviceType LBClientSslProfile")
			}
		}
	}

	if l.diff.HasChange(serviceTypes) {
		service := l.diff.Get(serviceTypes)
		profileType := l.diff.Get(serverProfile)
		if service == LBServerSslProfile {
			if len((profileType).([]interface{})) == 0 {
				return fmt.Errorf("please provide server_profile configurations for serviceType LBServerSslProfile")
			}
		}
	}
	return nil
}
