//  (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package diffvalidation

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	applicationTypes = "type"
	HTTP             = "http"
	UDP              = "udp"
	TCP              = "tcp"

	TCPProfile  = "tcp_application_profile"
	UDPProfile  = "udp_application_profile"
	HTTPProfile = "http_application_profile"

	persistenceTypes = "persistence"
	SOURCE_IP        = "SOURCE_IP"
	COOKIE           = "COOKIE"

	CookieProfile = "cookie_persistence_profile"
	SourceProfile = "sourceip_persistence_profile"
)

type loadBalancerVirtualServer struct {
	diff *schema.ResourceDiff
}

func NewLoadBalancerVirtualServerValidate(diff *schema.ResourceDiff) *loadBalancerVirtualServer {
	return &loadBalancerVirtualServer{
		diff: diff,
	}
}

func (l *loadBalancerVirtualServer) validateProfile(type1 string, value string) error {
	profileType := l.diff.Get(value)
	if len((profileType).([]interface{})) == 0 {
		return fmt.Errorf("please provide " + value + " " + "configurations for Type" + " " + type1)
	}
	return nil
}

func (l *loadBalancerVirtualServer) DiffValidate() error {
	types := l.diff.Get(applicationTypes)
	switch types {
	case TCP:
		err := l.validateProfile(TCP, TCPProfile)
		if err != nil {
			return err
		}
	case UDP:
		err := l.validateProfile(UDP, UDPProfile)
		if err != nil {
			return err
		}

	case HTTP:
		err := l.validateProfile(HTTP, HTTPProfile)
		if err != nil {
			return err
		}
	}
	persTypes := l.diff.Get(persistenceTypes)
	switch persTypes {
	case COOKIE:
		err := l.validateProfile(COOKIE, CookieProfile)
		if err != nil {
			return err
		}
	case SOURCE_IP:
		err := l.validateProfile(SOURCE_IP, SourceProfile)
		if err != nil {
			return err
		}
	}
	return nil
}
