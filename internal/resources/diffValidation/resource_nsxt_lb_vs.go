//  (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package diffvalidation

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	applicationTypes = "type"
	http             = "http"
	udp              = "udp"
	tcp              = "tcp"

	tcpProfile  = "tcp_application_profile"
	udpProfile  = "udp_application_profile"
	httpProfile = "http_application_profile"

	persistenceTypes = "persistence"
	sourceIP         = "SOURCE_IP"
	cookie           = "COOKIE"

	CookieProfile = "cookie_persistence_profile"
	SourceProfile = "sourceip_persistence_profile"
)

type LoadBalancerVirtualServers struct {
	diff *schema.ResourceDiff
}

func NewLoadBalancerVirtualServerValidate(diff *schema.ResourceDiff) *LoadBalancerVirtualServers {
	return &LoadBalancerVirtualServers{
		diff: diff,
	}
}

func (l *LoadBalancerVirtualServers) validateProfile(profile_type string, value string) error {
	profileType := l.diff.Get(value)
	if len((profileType).([]interface{})) == 0 {
		return fmt.Errorf("please provide " + value + " " + "configurations for Type" + " " + profile_type)
	}

	return nil
}

func (l *LoadBalancerVirtualServers) DiffValidate() error {
	types := l.diff.Get(applicationTypes)
	switch types {
	case tcp:
		err := l.validateProfile(tcp, tcpProfile)
		if err != nil {
			return err
		}
	case udp:
		err := l.validateProfile(udp, udpProfile)
		if err != nil {
			return err
		}

	case http:
		err := l.validateProfile(http, httpProfile)
		if err != nil {
			return err
		}
	}
	persTypes := l.diff.Get(persistenceTypes)
	switch persTypes {
	case cookie:
		err := l.validateProfile(cookie, CookieProfile)
		if err != nil {
			return err
		}
	case sourceIP:
		err := l.validateProfile(sourceIP, SourceProfile)
		if err != nil {
			return err
		}
	}

	return nil
}
