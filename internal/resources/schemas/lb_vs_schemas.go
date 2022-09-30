package schemas

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	DSLBProfile = "hpegl_vmaas_load_balancer_profile"
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
					Description: "TCP application_profile Id, Get the `id` from " + DSLBProfile + " datasource to obtain the TCP application_profile Id",
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
					Description: "UDP application_profile Id, Get the `id` from " + DSLBProfile + " datasource to obtain the UDP application_profile Id",
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
					Description: "HTTP application_profile Id, Get the `id` from " + DSLBProfile + " datasource to obtain the HTTP application_profile Id",
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
					Description: "COOKIE persistence_profile Id, Get the `id` from " + DSLBProfile + " datasource to obtain the COOKIE persistence_profile Id",
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
					Description: "SOURCEIP persistence_profile Id, Get the `id` from " + DSLBProfile + " datasource to obtain the SOURCEIP persistence_profile Id",
				},
			},
		},
	}
}
