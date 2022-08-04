//  (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package diffvalidation

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	profiles        = "profile_type"
	http            = "http_profile"
	tcp             = "tcp_profile"
	udp             = "udp_profile"
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
	types := l.diff.Get(profiles)
	switch types {
	case applicationProfile:
		httpType := l.diff.Get(http)
		tcpType := l.diff.Get(tcp)
		udpType := l.diff.Get(udp)
		if len((httpType).([]interface{})) != 0 {
			err := l.validateProfilesTypes(httpType, LBHttpProfile)
			if err != nil {
				return err
			}
		} else if len((tcpType).([]interface{})) != 0 {
			err := l.validateProfilesTypes(tcpType, LBFastTcpProfile)
			if err != nil {
				return err
			}
		} else if len((udpType).([]interface{})) != 0 {
			err := l.validateProfilesTypes(udpType, LBFastUdpProfile)
			if err != nil {
				return err
			}
		}

	case persistenceProfile:
		cookieType := l.diff.Get(cookieProfile)
		sourceipType := l.diff.Get(sourceIpProfile)
		genericType := l.diff.Get(genericProfile)
		if len((cookieType).([]interface{})) != 0 {
			err := l.validateProfilesTypes(cookieType, LBCookiePersistenceProfile)
			if err != nil {
				return err
			}
		} else if len((sourceipType).([]interface{})) != 0 {
			err := l.validateProfilesTypes(sourceipType, LBSourceIpPersistenceProfile)
			if err != nil {
				return err
			}
		} else if len((genericType).([]interface{})) != 0 {
			err := l.validateProfilesTypes(genericType, LBGenericPersistenceProfile)
			if err != nil {
				return err
			}
		}

	case sslProfile:
		serverType := l.diff.Get(serverProfile)
		clientType := l.diff.Get(clientProfile)
		if len((serverType).([]interface{})) != 0 {
			err := l.validateProfilesTypes(serverType, LBServerSslProfile)
			if err != nil {
				return err
			}
		} else if len((clientType).([]interface{})) != 0 {
			err := l.validateProfilesTypes(clientType, LBClientSslProfile)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (l *LoadBalancerProfile) validateProfilesTypes(profileType interface{}, serviceType string) error {
	for _, profile := range profileType.([]interface{}) {
		service_type := profile.(map[string]interface{})["service_type"].(string)
		if service_type != serviceType {
			return fmt.Errorf("please provide service_type as " + serviceType + " for the Configuration")
		}
	}

	return nil
}