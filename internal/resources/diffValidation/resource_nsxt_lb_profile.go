//  (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package diffvalidation

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	profiles        = "profile_type"
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

	applicationProfile = "application-profile"
	persistenceProfile = "persistence-profile"
	sslProfile         = "ssl-profile"
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
	pTypes := l.diff.Get(profiles)
	switch pTypes {
	case applicationProfile:
		p := l.diff.Get(httpProfile)
		if len((p).([]interface{})) != 0 {
			for _, profile := range p.([]interface{}) {
				profile_type := profile.(map[string]interface{})["service_type"].(string)
				if profile_type != LBHttpProfile {
					return fmt.Errorf("please provide service_type as " + LBHttpProfile + " for " + httpProfile + " configuration")
				}
			}
		}

		p = l.diff.Get(tcpProfile)
		if len((p).([]interface{})) != 0 {
			for _, profile := range p.([]interface{}) {
				profile_type := profile.(map[string]interface{})["service_type"].(string)
				if profile_type != LBFastTcpProfile {
					return fmt.Errorf("please provide service_type as " + LBFastTcpProfile + " for " + tcpProfile + " configuration")
				}
			}
		}

		p = l.diff.Get(udpProfile)
		if len((p).([]interface{})) != 0 {
			for _, profile := range p.([]interface{}) {
				profile_type := profile.(map[string]interface{})["service_type"].(string)
				if profile_type != LBFastUdpProfile {
					return fmt.Errorf("please provide service_type as " + LBFastUdpProfile + " for " + udpProfile + " configuration")
				}
			}
		}
	case persistenceProfile:
		p := l.diff.Get(cookieProfile)
		if len((p).([]interface{})) != 0 {
			for _, profile := range p.([]interface{}) {
				profile_type := profile.(map[string]interface{})["service_type"].(string)
				if profile_type != LBCookiePersistenceProfile {
					return fmt.Errorf("please provide service_type as " + LBCookiePersistenceProfile + " for " + cookieProfile + " configuration")
				}
			}
		}

		p = l.diff.Get(genericProfile)
		if len((p).([]interface{})) != 0 {
			for _, profile := range p.([]interface{}) {
				profile_type := profile.(map[string]interface{})["service_type"].(string)
				if profile_type != LBGenericPersistenceProfile {
					return fmt.Errorf("please provide service_type as " + LBGenericPersistenceProfile + " for " + genericProfile + " configuration")
				}
			}
		}

		p = l.diff.Get(sourceIpProfile)
		if len((p).([]interface{})) != 0 {
			for _, profile := range p.([]interface{}) {
				profile_type := profile.(map[string]interface{})["service_type"].(string)
				if profile_type != LBSourceIpPersistenceProfile {
					return fmt.Errorf("please provide service_type as " + LBSourceIpPersistenceProfile + " for " + sourceIpProfile + " configuration")
				}
			}
		}

	case sslProfile:
		p := l.diff.Get(clientProfile)
		if len((p).([]interface{})) != 0 {
			for _, profile := range p.([]interface{}) {
				profile_type := profile.(map[string]interface{})["service_type"].(string)
				if profile_type != LBClientSslProfile {
					return fmt.Errorf("please provide service_type as " + LBClientSslProfile + " for " + clientProfile + " configuration")
				}
			}
		}

		p = l.diff.Get(serverProfile)
		if len((p).([]interface{})) != 0 {
			for _, profile := range p.([]interface{}) {
				profile_type := profile.(map[string]interface{})["service_type"].(string)
				if profile_type != LBServerSslProfile {
					return fmt.Errorf("please provide service_type as " + LBServerSslProfile + " for " + serverProfile + " configuration")
				}
			}
		}
	}
	return nil
}
