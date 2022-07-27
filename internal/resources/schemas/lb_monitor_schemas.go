package schemas

import (
	"github.com/HewlettPackard/hpegl-vmaas-terraform-resources/internal/resources/validations"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func HttpMonitorSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		Description: "HTTP Monitor configuration",
		MaxItems:    1,
		ExactlyOneOf: []string{"http_monitor", "https_monitor", "icmp_monitor",
			"passive_monitor", "tcp_monitor", "udp_monitor"},
		ConflictsWith: []string{"https_monitor", "icmp_monitor",
			"passive_monitor", "tcp_monitor", "udp_monitor"},
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"fall_count": {
					Type:        schema.TypeInt,
					Optional:    true,
					Default:     3,
					Description: "Number of consecutive checks that must fail before marking it down.",
				},
				"interval": {
					Type:        schema.TypeInt,
					Default:     5,
					Description: "The frequency at which the system issues the monitor check (in seconds).",
					Optional:    true,
				},
				"monitor_port": {
					Type:        schema.TypeInt,
					Description: "Interval time for Network loadbalancer Monitor",
					Optional:    true,
				},
				"rise_count": {
					Type:        schema.TypeInt,
					Optional:    true,
					Default:     3,
					Description: "Number of consecutive checks that must pass before marking it up",
				},
				"timeout": {
					Type:        schema.TypeInt,
					Optional:    true,
					Default:     15,
					Description: "Timeout for Network loadbalancer Monitor",
				},
				"request_body": {
					Type:        schema.TypeString,
					Description: "String to send as HTTP health check request body. Valid only for certain HTTP methods like POST",
					Optional:    true,
				},
				"request_method": {
					Type:     schema.TypeString,
					Optional: true,
					ValidateDiagFunc: validations.StringInSlice([]string{"GET", "POST", "OPTIONS",
						"HEAD", "PUT"}, false),
					Default:     "GET",
					Description: "Health check method for HTTP monitor type. Valid values are GET, HEAD, PUT, POST and OPTIONS",
				},
				"request_url": {
					Type:        schema.TypeString,
					Description: "URL used for HTTP monitor",
					Optional:    true,
				},
				"request_version": {
					Type: schema.TypeString,
					ValidateDiagFunc: validations.StringInSlice([]string{"HTTP_VERSION_1_0",
						"HTTP_VERSION_1_1"}, false),
					Description: "HTTP request version. Valid values are HTTP_VERSION_1_0 and HTTP_VERSION_1_1",
					Optional:    true,
				},
				"response_data": {
					Type:        schema.TypeString,
					Description: "response data to get the monitor data",
					Optional:    true,
				},
				"response_status_codes": {
					Type:        schema.TypeString,
					Description: "HTTP response status code should be a valid HTTP status code",
					Optional:    true,
				},
			},
		},
	}
}

func HttpsMonitorSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		Description: "Https Monitor configuration",
		MaxItems:    1,
		ExactlyOneOf: []string{"http_monitor", "https_monitor", "icmp_monitor",
			"passive_monitor", "tcp_monitor", "udp_monitor"},
		ConflictsWith: []string{"http_monitor", "icmp_monitor",
			"passive_monitor", "tcp_monitor", "udp_monitor"},
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"fall_count": {
					Type:        schema.TypeInt,
					Optional:    true,
					Default:     3,
					Description: "Number of consecutive checks that must fail before marking it down.",
				},
				"interval": {
					Type:        schema.TypeInt,
					Default:     5,
					Description: "The frequency at which the system issues the monitor check (in seconds).",
					Optional:    true,
				},
				"monitor_port": {
					Type:        schema.TypeInt,
					Description: "If the monitor port is specified, it would override pool member port setting for healthcheck. A port range is not supported",
					Optional:    true,
				},
				"rise_count": {
					Type:        schema.TypeInt,
					Optional:    true,
					Default:     3,
					Description: "Number of consecutive checks that must pass before marking it up",
				},
				"timeout": {
					Type:        schema.TypeInt,
					Optional:    true,
					Default:     15,
					Description: "Number of seconds the target has to respond to the monitor request",
				},
				"request_body": {
					Type:        schema.TypeString,
					Description: "String to send as HTTPs health check request body. Valid only for certain HTTPs methods like POST",
					Optional:    true,
				},
				"request_method": {
					Type:     schema.TypeString,
					Optional: true,
					ValidateDiagFunc: validations.StringInSlice([]string{"GET", "POST", "OPTIONS",
						"HEAD", "PUT"}, false),
					Default:     "GET",
					Description: "Health check method for HTTPs monitor type. Valid values are GET, HEAD, PUT, POST and OPTIONS",
				},
				"request_url": {
					Type:        schema.TypeString,
					Description: "URL used for HTTP monitor",
					Optional:    true,
				},
				"request_version": {
					Type: schema.TypeString,
					ValidateDiagFunc: validations.StringInSlice([]string{"HTTP_VERSION_1_0",
						"HTTP_VERSION_1_1"}, false),
					Description: "HTTP request version. Valid values are HTTP_VERSION_1_0 and HTTP_VERSION_1_1",
					Optional:    true,
				},
				"response_data": {
					Type:        schema.TypeString,
					Description: "response data to get the monitor data",
					Optional:    true,
				},
				"response_status_codes": {
					Type:        schema.TypeString,
					Description: "HTTPs response status code should be a valid HTTPs status code",
					Optional:    true,
				},
			},
		},
	}
}

func IcmpMonitorSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		Description: "Icmp Monitor configuration",
		MaxItems:    1,
		ExactlyOneOf: []string{"http_monitor", "https_monitor", "icmp_monitor",
			"passive_monitor", "tcp_monitor", "udp_monitor"},
		ConflictsWith: []string{"http_monitor", "https_monitor",
			"passive_monitor", "tcp_monitor", "udp_monitor"},
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"fall_count": {
					Type:        schema.TypeInt,
					Optional:    true,
					Default:     3,
					Description: "Number of consecutive checks that must fail before marking it down.",
				},
				"interval": {
					Type:        schema.TypeInt,
					Default:     5,
					Description: "The frequency at which the system issues the monitor check (in seconds).",
					Optional:    true,
				},
				"monitor_port": {
					Type:        schema.TypeInt,
					Description: "If the monitor port is specified, it would override pool member port setting for healthcheck. A port range is not supported",
					Optional:    true,
				},
				"rise_count": {
					Type:        schema.TypeInt,
					Optional:    true,
					Default:     3,
					Description: "Number of consecutive checks that must pass before marking it up",
				},
				"timeout": {
					Type:        schema.TypeInt,
					Optional:    true,
					Default:     15,
					Description: "Number of seconds the target has to respond to the monitor request",
				},
				"data_length": {
					Type:        schema.TypeInt,
					Default:     56,
					Description: "data length is for the ICMP monitor type",
					Optional:    true,
				},
			},
		},
	}
}

func PassiveMonitorSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		Description: "Passive Monitor configuration",
		MaxItems:    1,
		ExactlyOneOf: []string{"http_monitor", "https_monitor", "icmp_monitor",
			"passive_monitor", "tcp_monitor", "udp_monitor"},
		ConflictsWith: []string{"http_monitor", "https_monitor",
			"icmp_monitor", "tcp_monitor", "udp_monitor"},
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"timeout": {
					Type:        schema.TypeInt,
					Optional:    true,
					Default:     15,
					Description: "Number of seconds the target has to respond to the monitor request",
				},
				"max_fail": {
					Type:        schema.TypeInt,
					Default:     5,
					Description: "maximum failure for the ICMP monitor type",
					Optional:    true,
				},
			},
		},
	}
}

func TcpMonitorSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		Description: "Tcp Monitor configuration",
		MaxItems:    1,
		ExactlyOneOf: []string{"http_monitor", "https_monitor", "icmp_monitor",
			"passive_monitor", "tcp_monitor", "udp_monitor"},
		ConflictsWith: []string{"http_monitor", "https_monitor",
			"icmp_monitor", "passive_monitor", "udp_monitor"},
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"fall_count": {
					Type:        schema.TypeInt,
					Optional:    true,
					Default:     3,
					Description: "Number of consecutive checks that must fail before marking it down.",
				},
				"interval": {
					Type:        schema.TypeInt,
					Default:     5,
					Description: "The frequency at which the system issues the monitor check (in seconds).",
					Optional:    true,
				},
				"monitor_port": {
					Type:        schema.TypeInt,
					Description: "If the monitor port is specified, it would override pool member port setting for healthcheck. A port range is not supported",
					Optional:    true,
				},
				"rise_count": {
					Type:        schema.TypeInt,
					Optional:    true,
					Default:     3,
					Description: "Number of consecutive checks that must pass before marking it up",
				},
				"timeout": {
					Type:        schema.TypeInt,
					Optional:    true,
					Default:     15,
					Description: "Number of seconds the target has to respond to the monitor request",
				},
				"request_body": {
					Type:        schema.TypeString,
					Description: "String to send as TCP health check request body. Valid only for certain TCP methods like POST",
					Optional:    true,
				},
				"response_data": {
					Type:        schema.TypeString,
					Description: "response data to get the monitor data",
					Optional:    true,
				},
			},
		},
	}
}

func UdpMonitorSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		Description: "Udp Monitor configuration",
		MaxItems:    1,
		ExactlyOneOf: []string{"http_monitor", "https_monitor", "icmp_monitor",
			"passive_monitor", "tcp_monitor", "udp_monitor"},
		ConflictsWith: []string{"http_monitor", "https_monitor",
			"icmp_monitor", "passive_monitor", "tcp_monitor"},
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"fall_count": {
					Type:        schema.TypeInt,
					Optional:    true,
					Default:     3,
					Description: "Number of consecutive checks that must fail before marking it down.",
				},
				"interval": {
					Type:        schema.TypeInt,
					Default:     5,
					Description: "The frequency at which the system issues the monitor check (in seconds).",
					Optional:    true,
				},
				"monitor_port": {
					Type:        schema.TypeInt,
					Description: "If the monitor port is specified, it would override pool member port setting for healthcheck. A port range is not supported",
					Optional:    true,
				},
				"rise_count": {
					Type:        schema.TypeInt,
					Optional:    true,
					Default:     3,
					Description: "Number of consecutive checks that must pass before marking it up",
				},
				"timeout": {
					Type:        schema.TypeInt,
					Optional:    true,
					Default:     15,
					Description: "Number of seconds the target has to respond to the monitor request",
				},
				"request_body": {
					Type:        schema.TypeString,
					Description: "String to send as UDP health check request body. Valid only for certain UDP methods like POST",
					Optional:    true,
				},
				"response_data": {
					Type:        schema.TypeString,
					Description: "response data to get the monitor data",
					Optional:    true,
				},
			},
		},
	}
}
