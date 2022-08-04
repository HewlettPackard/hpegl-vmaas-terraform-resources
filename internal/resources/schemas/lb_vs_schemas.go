package schemas

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TCPAppProfileSchema() *schema.Schema {
	return &schema.Schema{
		Type:          schema.TypeList,
		Optional:      true,
		Description:   "TCP Profile configuration",
		MaxItems:      1,
		ExactlyOneOf:  []string{"tcp_application_profile", "udp_application_profile", "http_application_profile"},
		ConflictsWith: []string{"udp_application_profile", "http_application_profile"},
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"application_profile": {
					Type:        schema.TypeInt,
					Required:    true,
					Description: "Network Loadbalancer Supported values are `SOURCE_IP`,`COOKIE`",
				},
			},
		},
	}
}

func UDPAppProfileSchema() *schema.Schema {
	return &schema.Schema{
		Type:          schema.TypeList,
		Optional:      true,
		Description:   "UDP profile configuration",
		MaxItems:      1,
		ExactlyOneOf:  []string{"tcp_application_profile", "udp_application_profile", "http_application_profile"},
		ConflictsWith: []string{"tcp_application_profile", "http_application_profile"},
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"application_profile": {
					Type:        schema.TypeInt,
					Required:    true,
					Description: "Network Loadbalancer Supported values are `SOURCE_IP`,`COOKIE`",
				},
			},
		},
	}
}

func HTTPAppProfileSchema() *schema.Schema {
	return &schema.Schema{
		Type:          schema.TypeList,
		Optional:      true,
		Description:   "HTTP Profile configuration",
		MaxItems:      1,
		ExactlyOneOf:  []string{"tcp_application_profile", "udp_application_profile", "http_application_profile"},
		ConflictsWith: []string{"tcp_application_profile", "udp_application_profile"},
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"application_profile": {
					Type:        schema.TypeInt,
					Required:    true,
					Description: "Network Loadbalancer Supported values are `SOURCE_IP`,`COOKIE`",
				},
			},
		},
	}
}

func CookiePersProfileSchema() *schema.Schema {
	return &schema.Schema{
		Type:          schema.TypeList,
		Optional:      true,
		Description:   "Cookie profile configuration",
		MaxItems:      1,
		ExactlyOneOf:  []string{"cookie_persistence_profile", "sourceip_persistence_profile"},
		ConflictsWith: []string{"sourceip_persistence_profile"},
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"persistence_profile": {
					Type:        schema.TypeInt,
					Required:    true,
					Description: "Network Loadbalancer Supported values are `SOURCE_IP`,`COOKIE`",
				},
			},
		},
	}
}

func SourceipPersProfileSchema() *schema.Schema {
	return &schema.Schema{
		Type:          schema.TypeList,
		Optional:      true,
		Description:   "HTTP profile configuration",
		MaxItems:      1,
		ExactlyOneOf:  []string{"cookie_persistence_profile", "sourceip_persistence_profile"},
		ConflictsWith: []string{"cookie_persistence_profile"},
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"persistence_profile": {
					Type:        schema.TypeInt,
					Required:    true,
					Description: "Network Loadbalancer Supported values are `SOURCE_IP`,`COOKIE`",
				},
			},
		},
	}
}
