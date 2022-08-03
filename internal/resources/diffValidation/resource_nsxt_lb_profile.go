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
		err := l.validateProfileTypes(httpType, tcpType, udpType, applicationProfile)
		if err != nil {
			return err
		}
	case persistenceProfile:
		cookieType := l.diff.Get(cookieProfile)
		sourceipType := l.diff.Get(sourceIpProfile)
		genericType := l.diff.Get(genericProfile)
		err := l.validateProfileTypes(cookieType, sourceipType, genericType, persistenceProfile)
		if err != nil {
			return err
		}

	case sslProfile:
		serverType := l.diff.Get(serverProfile)
		clientType := l.diff.Get(clientProfile)
		err := l.validateProfileTypes(serverType, clientType, nil, sslProfile)
		if err != nil {
			return err
		}
	}
	return nil
}

func (l *LoadBalancerProfile) validateProfileTypes(serviceType1 interface{}, serviceType2 interface{}, serviceType3 interface{}, profileType string) error {

	if profileType == applicationProfile {
		if len((serviceType1).([]interface{})) != 0 && len((serviceType2).([]interface{})) == 0 && len((serviceType3).([]interface{})) == 0 {
			for _, profile := range serviceType1.([]interface{}) {
				profile_type := profile.(map[string]interface{})["service_type"].(string)
				if profile_type != LBHttpProfile {
					return fmt.Errorf("please provide service_type as " + LBHttpProfile + " for HTTP Profile Configuration")
				}
			}
		} else if len((serviceType2).([]interface{})) != 0 && len((serviceType1).([]interface{})) == 0 && len((serviceType3).([]interface{})) == 0 {
			for _, profile := range serviceType2.([]interface{}) {
				profile_type := profile.(map[string]interface{})["service_type"].(string)
				if profile_type != LBFastTcpProfile {
					return fmt.Errorf("please provide service_type as " + LBFastTcpProfile + " for TCP Profile Configuration")
				}
			}
		} else if len((serviceType3).([]interface{})) != 0 && len((serviceType2).([]interface{})) == 0 && len((serviceType1).([]interface{})) == 0 {
			for _, profile := range serviceType3.([]interface{}) {
				profile_type := profile.(map[string]interface{})["service_type"].(string)
				if profile_type != LBFastUdpProfile {
					return fmt.Errorf("please provide service_type as " + LBFastUdpProfile + " for UDP Profile Configuration")
				}
			}
		}
	} else if profileType == persistenceProfile {
		if len((serviceType1).([]interface{})) != 0 && len((serviceType2).([]interface{})) == 0 && len((serviceType3).([]interface{})) == 0 {
			for _, profile := range serviceType1.([]interface{}) {
				profile_type := profile.(map[string]interface{})["service_type"].(string)
				if profile_type != LBCookiePersistenceProfile {
					return fmt.Errorf("please provide service_type as " + LBCookiePersistenceProfile + " for COOKIE Profile Configuration")
				}
			}
		} else if len((serviceType2).([]interface{})) != 0 && len((serviceType1).([]interface{})) == 0 && len((serviceType3).([]interface{})) == 0 {
			for _, profile := range serviceType2.([]interface{}) {
				profile_type := profile.(map[string]interface{})["service_type"].(string)
				if profile_type != LBSourceIpPersistenceProfile {
					return fmt.Errorf("please provide service_type as " + LBSourceIpPersistenceProfile + " for SOURCEIP Profile Configuration")
				}
			}
		} else if len((serviceType3).([]interface{})) != 0 && len((serviceType2).([]interface{})) == 0 && len((serviceType1).([]interface{})) == 0 {
			for _, profile := range serviceType3.([]interface{}) {
				profile_type := profile.(map[string]interface{})["service_type"].(string)
				if profile_type != LBGenericPersistenceProfile {
					return fmt.Errorf("please provide service_type as " + LBGenericPersistenceProfile + " for GENERIC Profile Configuration")
				}
			}
		}
	} else if profileType == sslProfile {
		if len((serviceType1).([]interface{})) != 0 && len((serviceType2).([]interface{})) == 0 {
			for _, profile := range serviceType1.([]interface{}) {
				profile_type := profile.(map[string]interface{})["service_type"].(string)
				if profile_type != LBServerSslProfile {
					return fmt.Errorf("please provide service_type as " + LBServerSslProfile + " for SSL-SERVER Profile Configuration")
				}
			}
		} else if len((serviceType2).([]interface{})) != 0 && len((serviceType1).([]interface{})) == 0 {
			for _, profile := range serviceType2.([]interface{}) {
				profile_type := profile.(map[string]interface{})["service_type"].(string)
				if profile_type != LBClientSslProfile {
					return fmt.Errorf("please provide service_type as " + LBClientSslProfile + " for SSL-CLIENT Profile Configuration")
				}
			}
		}
	}
	return nil
}
