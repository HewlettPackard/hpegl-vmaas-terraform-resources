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
		err := l.validateAppProfile(http, tcp, udp)
		if err != nil {
			return err
		}
	case persistenceProfile:
		err := l.validatePersProfile(cookieProfile, sourceIpProfile, genericProfile)
		if err != nil {
			return err
		}

	case sslProfile:
		err := l.validateSslProfile(serverProfile, clientProfile)
		if err != nil {
			return err
		}
	}
	return nil
}

func (l *LoadBalancerProfile) validateAppProfile(http string, tcp string, udp string) error {
	httpType := l.diff.Get(http)
	tcpType := l.diff.Get(tcp)
	udpType := l.diff.Get(udp)

	if len((httpType).([]interface{})) != 0 && len((tcpType).([]interface{})) == 0 && len((udpType).([]interface{})) == 0 {
		for _, profile := range httpType.([]interface{}) {
			profile_type := profile.(map[string]interface{})["service_type"].(string)
			if profile_type != LBHttpProfile {
				return fmt.Errorf("please provide service_type as " + LBHttpProfile + " for " + http + " configuration")
			}
		}
	}

	if len((httpType).([]interface{})) == 0 && len((tcpType).([]interface{})) != 0 && len((udpType).([]interface{})) == 0 {
		for _, profile := range tcpType.([]interface{}) {
			profile_type := profile.(map[string]interface{})["service_type"].(string)
			if profile_type != LBFastTcpProfile {
				return fmt.Errorf("please provide service_type as " + LBFastTcpProfile + " for " + tcp + " configuration")
			}
		}
	}

	if len((httpType).([]interface{})) == 0 && len((tcpType).([]interface{})) == 0 && len((udpType).([]interface{})) != 0 {
		for _, profile := range udpType.([]interface{}) {
			profile_type := profile.(map[string]interface{})["service_type"].(string)
			if profile_type != LBFastUdpProfile {
				return fmt.Errorf("please provide service_type as " + LBFastUdpProfile + " for " + udp + " configuration")
			}
		}
	}

	return nil
}
func (l *LoadBalancerProfile) validatePersProfile(cookie string, sourceIp string, generic string) error {
	cookieType := l.diff.Get(cookie)
	sourceIpType := l.diff.Get(sourceIp)
	genericType := l.diff.Get(generic)

	if len((cookieType).([]interface{})) != 0 && len((sourceIpType).([]interface{})) == 0 && len((genericType).([]interface{})) == 0 {
		for _, profile := range cookieType.([]interface{}) {
			profile_type := profile.(map[string]interface{})["service_type"].(string)
			if profile_type != LBCookiePersistenceProfile {
				return fmt.Errorf("please provide service_type as " + LBCookiePersistenceProfile + " for " + cookie + " configuration")
			}
		}
	}

	if len((cookieType).([]interface{})) == 0 && len((sourceIpType).([]interface{})) != 0 && len((genericType).([]interface{})) == 0 {
		for _, profile := range sourceIpType.([]interface{}) {
			profile_type := profile.(map[string]interface{})["service_type"].(string)
			if profile_type != LBSourceIpPersistenceProfile {
				return fmt.Errorf("please provide service_type as " + LBSourceIpPersistenceProfile + " for " + sourceIp + " configuration")
			}
		}
	}

	if len((cookieType).([]interface{})) == 0 && len((sourceIpType).([]interface{})) == 0 && len((genericType).([]interface{})) != 0 {
		for _, profile := range genericType.([]interface{}) {
			profile_type := profile.(map[string]interface{})["service_type"].(string)
			if profile_type != LBGenericPersistenceProfile {
				return fmt.Errorf("please provide service_type as " + LBGenericPersistenceProfile + " for " + generic + " configuration")
			}
		}
	}

	return nil
}
func (l *LoadBalancerProfile) validateSslProfile(server string, client string) error {
	serverType := l.diff.Get(server)
	clientType := l.diff.Get(client)

	if len((serverType).([]interface{})) != 0 && len((clientType).([]interface{})) == 0 {
		for _, profile := range serverType.([]interface{}) {
			profile_type := profile.(map[string]interface{})["service_type"].(string)
			if profile_type != LBServerSslProfile {
				return fmt.Errorf("please provide service_type as " + LBServerSslProfile + " for " + server + " configuration")
			}
		}
	}

	if len((serverType).([]interface{})) == 0 && len((clientType).([]interface{})) != 0 {
		for _, profile := range clientType.([]interface{}) {
			profile_type := profile.(map[string]interface{})["service_type"].(string)
			if profile_type != LBClientSslProfile {
				return fmt.Errorf("please provide service_type as " + LBClientSslProfile + " for " + client + " configuration")
			}
		}
	}

	return nil
}
