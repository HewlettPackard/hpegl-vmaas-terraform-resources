package schemas

import (
	"github.com/HewlettPackard/hpegl-vmaas-terraform-resources/internal/resources/validations"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func HTTPProfileSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		Description: "HTTP Profile configuration",
		MaxItems:    1,
		ExactlyOneOf: []string{
			"http_profile",
			"tcp_profile",
			"udp_profile",
			"cookie_profile",
			"sourceip_profile",
			"generic_profile",
			"client_profile",
			"server_profile",
		},
		ConflictsWith: []string{
			"tcp_profile",
			"udp_profile",
			"cookie_profile",
			"sourceip_profile",
			"generic_profile",
			"client_profile",
			"server_profile",
		},
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"service_type": {
					Type: schema.TypeString,
					ValidateDiagFunc: validations.StringInSlice([]string{
						"LBHttpProfile",
						"LBFastTcpProfile",
						"LBFastUdpProfile",
						"LBClientSslProfile",
						"LBServerSslProfile",
						"LBCookiePersistenceProfile",
						"LBGenericPersistenceProfile",
						"LBSourceIpPersistenceProfile",
					}, false),
					Required:     true,
					InputDefault: "LBHttpProfile",
					Description:  "Provide the Supported values for serviceTypes",
				},
				"http_idle_timeout": {
					Type:             schema.TypeInt,
					Optional:         true,
					Default:          15,
					ValidateDiagFunc: validations.IntBetween(1, 5400),
					Description:      "Timeout in seconds to specify how long an HTTP application can remain idle",
				},
				"request_header_size": {
					Type:             schema.TypeInt,
					Optional:         true,
					Default:          1024,
					ValidateDiagFunc: validations.IntBetween(1, 65536),
					Description:      "Specify the maximum buffer size in bytes used to store HTTP request headers",
				},
				"response_header_size": {
					Type:             schema.TypeInt,
					Optional:         true,
					ValidateDiagFunc: validations.IntBetween(1, 65536),
					Default:          4096,
					Description:      "Specify the maximum buffer size in bytes used to store HTTP response headers.",
				},
				"redirection": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "If a website is temporarily down or has moved, incoming requests for that virtual server can be temporarily redirected to a URL specified here.",
				},
				"x_forwarded_for": {
					Type: schema.TypeString,
					ValidateDiagFunc: validations.StringInSlice([]string{
						"INSERT",
						"REPLACE",
					}, false),
					Required:     true,
					InputDefault: "INSERT",
					Description:  "When this value is set, the x_forwarded_for header in the incoming request will be inserted or replaced.",
				},
				"request_body_size": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Enter value for the maximum size of the buffer used to store the HTTP request body",
				},
				"response_timeout": {
					Type:             schema.TypeInt,
					Optional:         true,
					Default:          60,
					ValidateDiagFunc: validations.IntBetween(1, 2147483647),
					Description:      "Number of seconds waiting for the server response before the connection is closed.",
				},
				"ntlm_authentication": {
					Type:        schema.TypeBool,
					Optional:    true,
					Description: "Toggle the button for the load balancer to turn off TCP multiplexing and enable HTTP keep-alive.",
				},
			},
		},
	}
}

func TCPProfileSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		Description: "TCP Profile configuration",
		MaxItems:    1,
		ExactlyOneOf: []string{
			"http_profile",
			"tcp_profile",
			"udp_profile",
			"cookie_profile",
			"sourceip_profile",
			"generic_profile",
			"client_profile",
			"server_profile",
		},
		ConflictsWith: []string{
			"http_profile",
			"udp_profile",
			"cookie_profile",
			"sourceip_profile",
			"generic_profile",
			"client_profile",
			"server_profile",
		},
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"service_type": {
					Type: schema.TypeString,
					ValidateDiagFunc: validations.StringInSlice([]string{
						"LBHttpProfile",
						"LBFastTcpProfile",
						"LBFastUdpProfile",
						"LBClientSslProfile",
						"LBServerSslProfile",
						"LBCookiePersistenceProfile",
						"LBGenericPersistenceProfile",
						"LBSourceIpPersistenceProfile",
					}, false),
					Required:     true,
					InputDefault: "LBFastTcpProfile",
					Description:  "Provide the  Supported values for serviceType",
				},
				"fast_tcp_idle_timeout": {
					Type:             schema.TypeInt,
					Optional:         true,
					Default:          1800,
					ValidateDiagFunc: validations.IntBetween(1, 2147483647),
					Description:      "Timeout in seconds to specify how long an idle TCP connection in ESTABLISHED",
				},
				"connection_close_timeout": {
					Type:             schema.TypeInt,
					Optional:         true,
					Default:          8,
					ValidateDiagFunc: validations.IntBetween(1, 60),
					Description:      "Timeout in seconds to specify how long a closed TCP connection",
				},
				"ha_flow_mirroring": {
					Type:        schema.TypeBool,
					Optional:    true,
					Default:     false,
					Description: "Toggle the button to make all the flows to the associated virtual server mirrored to the HA standby node",
				},
			},
		},
	}
}

func UDPProfileSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		Description: "UDP Profile configuration",
		MaxItems:    1,
		ExactlyOneOf: []string{
			"http_profile",
			"tcp_profile",
			"udp_profile",
			"cookie_profile",
			"sourceip_profile",
			"generic_profile",
			"client_profile",
			"server_profile",
		},
		ConflictsWith: []string{
			"http_profile",
			"tcp_profile",
			"cookie_profile",
			"sourceip_profile",
			"generic_profile",
			"client_profile",
			"server_profile",
		},
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"service_type": {
					Type: schema.TypeString,
					ValidateDiagFunc: validations.StringInSlice([]string{
						"LBHttpProfile",
						"LBFastTcpProfile",
						"LBFastUdpProfile",
						"LBClientSslProfile",
						"LBServerSslProfile",
						"LBCookiePersistenceProfile",
						"LBGenericPersistenceProfile",
						"LBSourceIpPersistenceProfile",
					}, false),
					Required:     true,
					InputDefault: "LBFastUdpProfile",
					Description:  "Provide the  Supported values for serviceType",
				},
				"fast_udp_idle_timeout": {
					Type:             schema.TypeInt,
					Optional:         true,
					Default:          300,
					ValidateDiagFunc: validations.IntBetween(1, 2147483647),
					Description:      "Timeout in seconds to specify how long an idle UDP connection in ESTABLISHED",
				},
				"ha_flow_mirroring": {
					Type:        schema.TypeBool,
					Optional:    true,
					Default:     false,
					Description: "Toggle the button to make all the flows to the associated virtual server mirrored to the HA standby node",
				},
			},
		},
	}
}

func CookieProfileSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		Description: "Cookie Profile configuration",
		MaxItems:    1,
		ExactlyOneOf: []string{
			"http_profile",
			"tcp_profile",
			"udp_profile",
			"cookie_profile",
			"sourceip_profile",
			"generic_profile",
			"client_profile",
			"server_profile",
		},
		ConflictsWith: []string{
			"http_profile",
			"tcp_profile",
			"udp_profile",
			"sourceip_profile",
			"generic_profile",
			"client_profile",
			"server_profile",
		},
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"service_type": {
					Type: schema.TypeString,
					ValidateDiagFunc: validations.StringInSlice([]string{
						"LBHttpProfile",
						"LBFastTcpProfile",
						"LBFastUdpProfile",
						"LBClientSslProfile",
						"LBServerSslProfile",
						"LBCookiePersistenceProfile",
						"LBGenericPersistenceProfile",
						"LBSourceIpPersistenceProfile",
					}, false),
					Required:     true,
					InputDefault: "LBCookiePersistenceProfile",
					Description:  "Provide the  Supported values for serviceType",
				},
				"cookie_name": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "cookie_name for Network Load balancer Profile",
				},
				"cookie_fallback": {
					Type:     schema.TypeBool,
					Optional: true,
					Default:  true,
					Description: "Cookie fallback enabled means,so that the client request is rejected" +
						"if cookie points to a server that is in a DISABLED or is in a DOWN state",
				},
				"cookie_garbling": {
					Type:     schema.TypeBool,
					Optional: true,
					Default:  true,
					Description: "When garbling is disabled, the cookie server IP address" +
						"and port information is in a plain text",
				},
				"cookie_mode": {
					Type: schema.TypeString,
					ValidateDiagFunc: validations.StringInSlice([]string{
						"INSERT",
						"PREFIX",
						"REWRITE",
					}, false),
					Required:     true,
					InputDefault: "INSERT",
					Description:  "The cookie persistence mode",
				},
				"cookie_type": {
					Type: schema.TypeString,
					ValidateDiagFunc: validations.StringInSlice([]string{
						"LBPersistenceCookieTime",
						"LBSessionCookieTime",
					}, false),
					Required:     true,
					InputDefault: "LBSessionCookieTime",
					Description:  "Provide the  Supported values for cookie_type",
				},
				"cookie_domain": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Enter the domain name. HTTP cookie domain can be configured only in the `INSERT` mode",
				},
				"cookie_path": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Enter the cookie URL path. HTTP cookie path can be set only in the `INSERT` mode",
				},
				"max_idle_time": {
					Type:        schema.TypeInt,
					Optional:    true,
					Description: "Enter the time in seconds that the cookie type can be idle before a cookie expires",
				},
				"share_persistence": {
					Type:        schema.TypeBool,
					Optional:    true,
					Description: "Toggle the button to share the persistence so that all virtual servers this profile is associated with can share the persistence table",
				},
			},
		},
	}
}

func SourceIPProfileSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		Description: "SourceIP Profile configuration",
		MaxItems:    1,
		ExactlyOneOf: []string{
			"http_profile",
			"tcp_profile",
			"udp_profile",
			"cookie_profile",
			"sourceip_profile",
			"generic_profile",
			"client_profile",
			"server_profile",
		},
		ConflictsWith: []string{
			"http_profile",
			"tcp_profile",
			"udp_profile",
			"cookie_profile",
			"generic_profile",
			"client_profile",
			"server_profile",
		},
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"service_type": {
					Type: schema.TypeString,
					ValidateDiagFunc: validations.StringInSlice([]string{
						"LBHttpProfile",
						"LBFastTcpProfile",
						"LBFastUdpProfile",
						"LBClientSslProfile",
						"LBServerSslProfile",
						"LBCookiePersistenceProfile",
						"LBGenericPersistenceProfile",
						"LBSourceIpPersistenceProfile",
					}, false),
					Required:     true,
					InputDefault: "LBSourceIpPersistenceProfile",
					Description:  "Provide the  Supported values for serviceType",
				},
				"share_persistence": {
					Type:     schema.TypeBool,
					Optional: true,
					Default:  false,
					Description: "Toggle the button to share the persistence so that all virtual servers" +
						"this profile is associated with can share the persistence table",
				},
				"ha_persistence_mirroring": {
					Type:     schema.TypeBool,
					Optional: true,
					Default:  false,
					Description: "Toggle the button to synchronize persistence entries to the HA peer. When HA persistence mirroring is enabled," +
						"the client IP persistence remains in the case of load balancer failover",
				},
				"persistence_entry_timeout": {
					Type:             schema.TypeInt,
					Optional:         true,
					Default:          300,
					ValidateDiagFunc: validations.IntBetween(1, 2147483647),
					Description:      "Persistence expiration time in seconds, counted from the time all the connections are completed",
				},
				"purge_entries_when_full": {
					Type:        schema.TypeBool,
					Optional:    true,
					Default:     true,
					Description: "When this option is enabled, the oldest entry is deleted to accept the newest entry in the persistence table",
				},
			},
		},
	}
}

func GenericProfileSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		Description: "Generic Profile configuration",
		MaxItems:    1,
		ExactlyOneOf: []string{
			"http_profile",
			"tcp_profile",
			"udp_profile",
			"cookie_profile",
			"sourceip_profile",
			"generic_profile",
			"client_profile",
			"server_profile",
		},
		ConflictsWith: []string{
			"http_profile",
			"tcp_profile",
			"udp_profile",
			"cookie_profile",
			"sourceip_profile",
			"client_profile",
			"server_profile",
		},
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"service_type": {
					Type: schema.TypeString,
					ValidateDiagFunc: validations.StringInSlice([]string{
						"LBHttpProfile",
						"LBFastTcpProfile",
						"LBFastUdpProfile",
						"LBClientSslProfile",
						"LBServerSslProfile",
						"LBCookiePersistenceProfile",
						"LBGenericPersistenceProfile",
						"LBSourceIpPersistenceProfile",
					}, false),
					Required:     true,
					InputDefault: "LBGenericPersistenceProfile",
					Description:  "Provide the  Supported values for serviceType",
				},
				"share_persistence": {
					Type:     schema.TypeBool,
					Optional: true,
					Default:  false,
					Description: "Toggle the button to share the persistence so" +
						"that all virtual servers this profile is associated with can share the persistence table",
				},
				"ha_persistence_mirroring": {
					Type:     schema.TypeBool,
					Optional: true,
					Default:  false,
					Description: "Toggle the button to synchronize persistence entries to the HA peer. When HA persistence mirroring is enabled," +
						"the client IP persistence remains in the case of load balancer failover.",
				},
				"persistence_entry_timeout": {
					Type:             schema.TypeInt,
					Optional:         true,
					Default:          300,
					ValidateDiagFunc: validations.IntBetween(1, 2147483647),
					Description: "Persistence expiration time in seconds," +
						"counted from the time all the connections are completed",
				},
			},
		},
	}
}

func ClientProfileSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		Description: "Client Profile configuration",
		MaxItems:    1,
		ExactlyOneOf: []string{
			"http_profile",
			"tcp_profile",
			"udp_profile",
			"cookie_profile",
			"sourceip_profile",
			"generic_profile",
			"client_profile",
			"server_profile",
		},
		ConflictsWith: []string{
			"http_profile",
			"tcp_profile",
			"udp_profile",
			"cookie_profile",
			"sourceip_profile",
			"generic_profile",
			"server_profile",
		},
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"service_type": {
					Type: schema.TypeString,
					ValidateDiagFunc: validations.StringInSlice([]string{
						"LBHttpProfile",
						"LBFastTcpProfile",
						"LBFastUdpProfile",
						"LBClientSslProfile",
						"LBServerSslProfile",
						"LBCookiePersistenceProfile",
						"LBGenericPersistenceProfile",
						"LBSourceIpPersistenceProfile",
					}, false),
					Required:     true,
					InputDefault: "LBClientSslProfile",
					Description:  "Provide the  Supported values for serviceType",
				},
				"ssl_suite": {
					Type: schema.TypeString,
					ValidateDiagFunc: validations.StringInSlice([]string{
						"BALANCED",
						"HIGH_SECURITY",
						"HIGH_COMPATIBILITY",
						"CUSTOM",
					}, false),
					Required:     true,
					InputDefault: "BALANCED",
					Description:  "Provide the  Supported values for ssl_suite",
				},
				"session_cache": {
					Type:     schema.TypeBool,
					Optional: true,
					Default:  true,
					Description: "To allow the SSL client and server to reuse previously negotiated security parameters avoiding" +
						"the expensive public key operation during an SSL handshake",
				},
				"session_cache_entry_timeout": {
					Type:     schema.TypeInt,
					Optional: true,
					Default:  300,
					Description: "Enter the cache timeout in seconds to specify how long the SSL session" +
						"parameters must be kept and can be reused",
				},
				"prefer_server_cipher": {
					Type:     schema.TypeBool,
					Optional: true,
					Default:  true,
					Description: "During SSL handshake as part of the SSL client sends an ordered list" +
						"of ciphers that it can support (or prefers) and typically server selects the first one from the top" +
						"of that list it can also support.For Perfect Forward Secrecy(PFS), server could override the client's preference",
				},
			},
		},
	}
}

func ServerProfileSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		Description: "Server Profile configuration",
		MaxItems:    1,
		ExactlyOneOf: []string{
			"http_profile",
			"tcp_profile",
			"udp_profile",
			"cookie_profile",
			"sourceip_profile",
			"generic_profile",
			"client_profile",
			"server_profile",
		},
		ConflictsWith: []string{
			"http_profile",
			"tcp_profile",
			"udp_profile",
			"cookie_profile",
			"sourceip_profile",
			"generic_profile",
			"client_profile",
		},
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"service_type": {
					Type: schema.TypeString,
					ValidateDiagFunc: validations.StringInSlice([]string{
						"LBHttpProfile",
						"LBFastTcpProfile",
						"LBFastUdpProfile",
						"LBClientSslProfile",
						"LBServerSslProfile",
						"LBCookiePersistenceProfile",
						"LBGenericPersistenceProfile",
						"LBSourceIpPersistenceProfile",
					}, false),
					Required:     true,
					InputDefault: "LBServerSslProfile",
					Description:  "Provide the  Supported values for serviceType",
				},
				"ssl_suite": {
					Type: schema.TypeString,
					ValidateDiagFunc: validations.StringInSlice([]string{
						"BALANCED",
						"HIGH_SECURITY",
						"HIGH_COMPATIBILITY",
						"CUSTOM",
					}, false),
					Required:     true,
					InputDefault: "BALANCED",
					Description:  "Provide the  Supported values for ssl_suite",
				},
				"session_cache": {
					Type:     schema.TypeBool,
					Required: true,
					Description: "To allow the SSL client and server to reuse previously negotiated security parameters avoiding" +
						"the expensive public key operation during an SSL handshake",
				},
			},
		},
	}
}
