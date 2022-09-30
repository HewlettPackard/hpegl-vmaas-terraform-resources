package schemas

import (
	"github.com/HewlettPackard/hpegl-vmaas-terraform-resources/internal/resources/validations"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func HTTPMonitorSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		Description: "HTTP Monitor configuration",
		MaxItems:    1,
		ExactlyOneOf: []string{
			"http_monitor",
			"https_monitor",
			"icmp_monitor",
			"passive_monitor",
			"tcp_monitor",
			"udp_monitor",
		},
		ConflictsWith: []string{
			"https_monitor",
			"icmp_monitor",
			"passive_monitor",
			"tcp_monitor",
			"udp_monitor",
		},
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
					Description: "Set the number of times the server is tested before it is considered as DOWN",
					Optional:    true,
				},
				"monitor_port": {
					Type:        schema.TypeInt,
					Description: "Set the value of the monitoring port.",
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
					Description: "Set the number of times the server is tested before it is considered as DOWN",
				},
				"request_body": {
					Type:        schema.TypeString,
					Description: "Enter the request body. Valid for the POST and PUT methods",
					Optional:    true,
				},
				"request_method": {
					Type:     schema.TypeString,
					Optional: true,
					ValidateDiagFunc: validations.StringInSlice([]string{
						"GET",
						"POST",
						"OPTIONS",
						"HEAD",
						"PUT",
					}, false),
					Default:     "GET",
					Description: "Select the method to detect the server status",
				},
				"request_url": {
					Type:        schema.TypeString,
					Description: "Enter the request URI for the method",
					Optional:    true,
				},
				"request_version": {
					Type: schema.TypeString,
					ValidateDiagFunc: validations.StringInSlice([]string{
						"HTTP_VERSION_1_0",
						"HTTP_VERSION_1_1",
					}, false),
					Description: "HTTP request version. Valid values are HTTP_VERSION_1_0 and HTTP_VERSION_1_1",
					Optional:    true,
				},
				"response_data": {
					Type: schema.TypeString,
					Description: "If the HTTP response body string and the HTTP health check response body match," +
						"then the server is considered as healthy",
					Optional: true,
				},
				"response_status_codes": {
					Type: schema.TypeString,
					Description: "Enter the string that the monitor expects to match in the status line of HTTP response body." +
						"The response code is a comma-separated list",
					Optional: true,
				},
			},
		},
	}
}

func HTTPSMonitorSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		Description: "Https Monitor configuration",
		MaxItems:    1,
		ExactlyOneOf: []string{
			"http_monitor",
			"https_monitor",
			"icmp_monitor",
			"passive_monitor",
			"tcp_monitor",
			"udp_monitor",
		},
		ConflictsWith: []string{
			"http_monitor",
			"icmp_monitor",
			"passive_monitor",
			"tcp_monitor",
			"udp_monitor",
		},
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
					Description: "Set the number of times the server is tested before it is considered as DOWN",
					Optional:    true,
				},
				"monitor_port": {
					Type:        schema.TypeInt,
					Description: "Set the value of the monitoring port",
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
					Description: "Set the number of times the server is tested before it is considered as DOWN",
				},
				"request_body": {
					Type:        schema.TypeString,
					Description: "Enter the request body. Valid for the POST and PUT methods",
					Optional:    true,
				},
				"request_method": {
					Type:     schema.TypeString,
					Optional: true,
					ValidateDiagFunc: validations.StringInSlice([]string{
						"GET",
						"POST",
						"OPTIONS",
						"HEAD",
						"PUT",
					}, false),
					Default:     "GET",
					Description: "Select the method to detect the server status",
				},
				"request_url": {
					Type:        schema.TypeString,
					Description: "Enter the request URI for the method",
					Optional:    true,
				},
				"request_version": {
					Type: schema.TypeString,
					ValidateDiagFunc: validations.StringInSlice([]string{
						"HTTP_VERSION_1_0",
						"HTTP_VERSION_1_1",
					}, false),
					Description: "HTTP request version. Valid values are HTTP_VERSION_1_0 and HTTP_VERSION_1_1",
					Optional:    true,
				},
				"response_data": {
					Type: schema.TypeString,
					Description: "If the HTTP response body string and the HTTP health check response body match," +
						"then the server is considered as healthy",
					Optional: true,
				},
				"response_status_codes": {
					Type: schema.TypeString,
					Description: "Enter the string that the monitor expects to match in the status line of HTTP response body." +
						"The response code is a comma-separated list",
					Optional: true,
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
		ExactlyOneOf: []string{
			"http_monitor",
			"https_monitor",
			"icmp_monitor",
			"passive_monitor",
			"tcp_monitor",
			"udp_monitor",
		},
		ConflictsWith: []string{
			"http_monitor",
			"https_monitor",
			"passive_monitor",
			"tcp_monitor",
			"udp_monitor",
		},
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"fall_count": {
					Type:        schema.TypeInt,
					Optional:    true,
					Default:     3,
					Description: "Number of consecutive checks that must fail before marking it down",
				},
				"interval": {
					Type:        schema.TypeInt,
					Default:     5,
					Description: "Set the number of times the server is tested before it is considered as DOWN",
					Optional:    true,
				},
				"monitor_port": {
					Type:        schema.TypeInt,
					Description: "Set the value of the monitoring port.",
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
					Description: "Set the number of times the server is tested before it is considered as DOWN",
				},
				"data_length": {
					Type:        schema.TypeInt,
					Default:     56,
					Description: "Maximum size of the ICMP data packet",
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
		ExactlyOneOf: []string{
			"http_monitor",
			"https_monitor",
			"icmp_monitor",
			"passive_monitor",
			"tcp_monitor",
			"udp_monitor",
		},
		ConflictsWith: []string{
			"http_monitor",
			"https_monitor",
			"icmp_monitor",
			"tcp_monitor",
			"udp_monitor",
		},
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"timeout": {
					Type:        schema.TypeInt,
					Optional:    true,
					Default:     15,
					Description: "Set the number of times the server is tested before it is considered as DOWN",
				},
				"max_fail": {
					Type:    schema.TypeInt,
					Default: 5,
					Description: "Set a value when the consecutive failures reach this value," +
						"the server is considered temporarily unavailable",
					Optional: true,
				},
			},
		},
	}
}

func TCPMonitorSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		Description: "Tcp Monitor configuration",
		MaxItems:    1,
		ExactlyOneOf: []string{
			"http_monitor",
			"https_monitor",
			"icmp_monitor",
			"passive_monitor",
			"tcp_monitor",
			"udp_monitor",
		},
		ConflictsWith: []string{
			"http_monitor",
			"https_monitor",
			"icmp_monitor",
			"passive_monitor",
			"udp_monitor",
		},
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
					Description: "Set the number of times the server is tested before it is considered as DOWN",
					Optional:    true,
				},
				"monitor_port": {
					Type:        schema.TypeInt,
					Description: "Set the value of the monitoring port.",
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
					Description: "Set the number of times the server is tested before it is considered as DOWN",
				},
				"request_body": {
					Type:        schema.TypeString,
					Description: "Enter the request body. Valid for the POST and PUT methods",
					Optional:    true,
				},
				"response_data": {
					Type: schema.TypeString,
					Description: "If the HTTP response body string and the HTTP health check response body match" +
						"then the server is considered as healthy",
					Optional: true,
				},
			},
		},
	}
}

func UDPMonitorSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		Description: "Udp Monitor configuration",
		MaxItems:    1,
		ExactlyOneOf: []string{
			"http_monitor",
			"https_monitor",
			"icmp_monitor",
			"passive_monitor",
			"tcp_monitor",
			"udp_monitor",
		},
		ConflictsWith: []string{
			"http_monitor",
			"https_monitor",
			"icmp_monitor",
			"passive_monitor",
			"tcp_monitor",
		},
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
					Description: "Set the number of times the server is tested before it is considered as DOWN",
					Optional:    true,
				},
				"monitor_port": {
					Type:        schema.TypeInt,
					Description: "Set the value of the monitoring port.",
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
					Description: "Set the number of times the server is tested before it is considered as DOWN",
				},
				"request_body": {
					Type:        schema.TypeString,
					Description: "Enter the request body. Valid for the POST and PUT methods",
					Optional:    true,
				},
				"response_data": {
					Type: schema.TypeString,
					Description: "If the HTTP response body string and the HTTP health check response body match," +
						"then the server is considered as healthy",
					Optional: true,
				},
			},
		},
	}
}
