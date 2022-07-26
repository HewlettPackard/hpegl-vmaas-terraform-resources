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
	serviceType := l.diff.Get(serviceTypes)
	switch serviceType {
	case LBHttpProfile:
		err := l.validateProfile(httpProfile, LBHttpProfile)
		if err != nil {
			return err
		}

		err = l.validateApplication(httpProfile, LBHttpProfile)
		if err != nil {
			return err
		}

	case LBFastTcpProfile:
		err := l.validateProfile(tcpProfile, LBFastTcpProfile)
		if err != nil {
			return err
		}

		err = l.validateApplication(tcpProfile, LBFastTcpProfile)
		if err != nil {
			return err
		}
	case LBFastUdpProfile:
		err := l.validateProfile(udpProfile, LBFastUdpProfile)
		if err != nil {
			return err
		}

		err = l.validateApplication(udpProfile, LBFastUdpProfile)
		if err != nil {
			return err
		}
	case LBCookiePersistenceProfile:
		err := l.validateProfile(cookieProfile, LBCookiePersistenceProfile)
		if err != nil {
			return err
		}

		err = l.validatePersistence(cookieProfile, LBCookiePersistenceProfile)
		if err != nil {
			return err
		}
	case LBSourceIpPersistenceProfile:
		err := l.validateProfile(sourceIpProfile, LBSourceIpPersistenceProfile)
		if err != nil {
			return err
		}

		err = l.validatePersistence(sourceIpProfile, LBSourceIpPersistenceProfile)
		if err != nil {
			return err
		}
	case LBGenericPersistenceProfile:
		err := l.validateProfile(genericProfile, LBGenericPersistenceProfile)
		if err != nil {
			return err
		}

		err = l.validatePersistence(genericProfile, LBGenericPersistenceProfile)
		if err != nil {
			return err
		}
	case LBClientSslProfile:
		err := l.validateProfile(clientProfile, LBClientSslProfile)
		if err != nil {
			return err
		}

		err = l.validateSsl(clientProfile, LBClientSslProfile)
		if err != nil {
			return err
		}
	case LBServerSslProfile:
		err := l.validateProfile(serverProfile, LBServerSslProfile)
		if err != nil {
			return err
		}

		err = l.validateSsl(serverProfile, LBServerSslProfile)
		if err != nil {
			return err
		}
	}
	return nil
}

func (l *LoadBalancerProfile) validateApplication(profileTypes string, service_type string) error {
	p := l.diff.Get(profileTypes)
	for _, profile := range p.([]interface{}) {
		profile_type := profile.(map[string]interface{})["profile_type"].(string)
		if profile_type != applicationProfile {
			return fmt.Errorf("please provide profile_type as " + applicationProfile + " " + "for serviceType" + " " + service_type)
		}
	}
	return nil
}

func (l *LoadBalancerProfile) validatePersistence(profileTypes string, service_type string) error {
	p := l.diff.Get(profileTypes)
	for _, profile := range p.([]interface{}) {
		profile_type := profile.(map[string]interface{})["profile_type"].(string)
		if profile_type != persistenceProfile {
			return fmt.Errorf("please provide profile_type as " + persistenceProfile + " " + " for serviceType" + " " + service_type)
		}
	}
	return nil
}

func (l *LoadBalancerProfile) validateSsl(profileTypes string, service_type string) error {
	p := l.diff.Get(profileTypes)
	for _, profile := range p.([]interface{}) {
		profile_type := profile.(map[string]interface{})["profile_type"].(string)
		if profile_type != sslProfile {
			return fmt.Errorf("please provide profile_type as " + sslProfile + " " + "for serviceType" + " " + service_type)
		}
	}
	return nil
}

func (l *LoadBalancerProfile) validateProfile(profile_type string, service_type string) error {
	profileType := l.diff.Get(profile_type)
	if len((profileType).([]interface{})) == 0 {
		return fmt.Errorf("please provide " + profile_type + " " + "configurations for serviceType" + " " + service_type)
	}
	return nil
}
