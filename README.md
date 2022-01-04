![Acceptance workflow](https://github.com/HewlettPackard/hpegl-vmaas-terraform-resources/actions/workflows/acc.yml/badge.svg)

# vmaas-terraform-resources

- [vmaas-terraform-resources](#vmaas-terraform-resources)
    * [Introduction](#introduction)
    * [Terraform versions](#terraform-versions)
    * [Terraform provider & v2.0 SDK](#terraform-provider--v20-sdk)
    * [Basic structure](#basic-structure)
    * [pkg directory](#pkg-directory)
        + [pkg/client](#pkgclient)
        + [pkg/resources](#pkgresources)
    * [internal](#internal)
        + [acceptance_test](#acceptance_test)
        + [resources](#resources)
        + [test-utils](#test-utils)
    * [Building and using the "dummy" service-specific provider](#building-and-using-the-dummy-service-specific-provider)
        + [Building the service-specific provider](#building-the-service-specific-provider)
        + [Using the service-specific provider](#using-the-service-specific-provider)

## Introduction

This repo contains VMaaS terraform provider code that we've used while developing
the hpegl provider [here](https://github.com/hpe-hcss/terraform-provider-hpegl).  It is
the exemplar service provider repo, and will form the basis of a Github template.

## Terraform versions

Terraform versions >= v0.13 should be used while developing the service provider code.

## Terraform provider & v2.0 SDK

We are mandating the use of the v2.0 Terraform SDK.  This version of the SDK is documented
[here](https://www.terraform.io/docs/extend/guides/v2-upgrade-guide.html#new-diagnostics-enabled-map-validators).  Perhaps
the two most pertinent enhancements are the passing of a context down to the CRUD implementation functions, and the
ability to bundle errors and warnings into a Diagnostics list ([]Diagnostic) that is returned from the CRUD functions.
These Diagnostics can be warnings or errors, and are processed by Terraform for presentation on the console.
The Terraform provider writing tutorial [here](https://learn.hashicorp.com/collections/terraform/providers) has been updated
to use the v2.0 SDK.

## Basic structure

The repo follows golang standards:

```bash
.
├── build
│   └── terraform-provider-hpegl
├── cmd
│   └── terraform-provider-hpegl
│       └── main.go
├── examples
│   └── main.tf
├── golangci-lint-config.yaml
├── go.mod
├── go.sum
├── internal
│   ├── acceptance_test
│   │   └── provider_test.go
│   ├── resources
│   │   ├── cluster_blueprint.go
│   │   └── cluster.go
│   └── test-utils
│       └── test-provider.go
├── Makefile
├── pkg
│   ├── client
│   │   ├── client.go
│   │   └── config.go
│   └── resources
│       └── registration.go
└── README.md
```

## pkg directory

The pkg directory contains implementations of the following interfaces from
[hpegl-provider-lib](https://github.com/hewlettpackard/hpegl-provider-lib):
* pkg/client: <br>
    [client.Initialisation](https://github.com/hewlettpackard/hpegl-provider-lib/blob/bcedebfa36825c4e78b4df73ecaa0e2ed77b92e2/pkg/client/client.go#L9)
    which is <br>
    ```go
    package client

    import "github.com/hewlettpackard/hpegl-provider-lib/pkg/provider"

    // Initialisation interface, service Client creation code will have to satisfy this interface
    // The hpegl provider will iterate over a slice of these to initialise service clients
    type Initialisation interface {
    	// NewClient is run by hpegl to initialise the service client
    	NewClient(config provider.ConfigData) (interface{}, error)

    	// ServiceName is used by hpegl, it returns the key to be used for the client returned by NewClient
    	// in the map[string]interface{} passed-down to provider code by terraform
    	ServiceName() string
    }
    ```

* pkg/resources: <br>
    [registration.ServiceRegistration](https://github.com/hewlettpackard/hpegl-provider-lib/blob/bcedebfa36825c4e78b4df73ecaa0e2ed77b92e2/pkg/registration/service_registration.go#L9)
    which is: <br>
    ```go
    package registration

    import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

    type ServiceRegistration interface {
    	// Name is the name of this Service
    	Name() string

    	// SupportedDataSources returns the supported Data Sources supported by this Service
    	SupportedDataSources() map[string]*schema.Resource

    	// SupportedResources returns the supported Resources supported by this Service
    	SupportedResources() map[string]*schema.Resource
    }
    ```

### pkg/client

This contains the implementation of the client.Initialisation interface from hpegl-provider-lib.
The struct that implements the interface is named InitialiseClient and this name should be retained
for all service provider code.  The hpegl provider will iterate over a slice of InitialiseClient structs
from each supported service to populate the map[string]interface{} struct of service Clients that is
passed-down to provider code by terraform and from which service provider code extracts the service Client.
This is the slice that is created in the hpegl provider:

```go
ppackage client

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Initialisation interface, service Client creation code will have to satisfy this interface
// The hpegl provider will iterate over a slice of these to initialise service clients
type Initialisation interface {
	// NewClient is run by hpegl to initialise the service client
	NewClient(r *schema.ResourceData) (interface{}, error)

	// ServiceName is used by hpegl, it returns the key to be used for the client returned by NewClient
	// in the map[string]interface{} passed-down to provider code by terraform
	ServiceName() string
}
```

The hpegl code that iterates over the slice is:

```go
package client

import (
	"fmt"
	"os"

	api_client "github.com/HewlettPackard/hpegl-vmaas-cmp-go-sdk/pkg/client"
	"github.com/HewlettPackard/hpegl-vmaas-cmp-go-sdk/pkg/models"
	cmp_client "github.com/HewlettPackard/hpegl-vmaas-terraform-resources/internal/cmp"
	"github.com/HewlettPackard/hpegl-vmaas-terraform-resources/pkg/constants"
	"github.com/HewlettPackard/hpegl-vmaas-terraform-resources/pkg/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hewlettpackard/hpegl-provider-lib/pkg/client"
	"github.com/tshihad/tftags"
)

// keyForGLClientMap is the key in the map[string]interface{} that is passed down by hpegl used to store *Client
// This must be unique, hpegl will error-out if it isn't
const keyForGLClientMap = "vmaasClient"

// Assert that InitialiseClient satisfies the client.Initialisation interface
var _ client.Initialisation = (*InitialiseClient)(nil)

// Client is the client struct that is used by the provider code
type Client struct {
	CmpClient *cmp_client.Client
}

// InitialiseClient is imported by hpegl from each service repo
type InitialiseClient struct{}

// NewClient takes an argument of all of the provider.ConfigData, and returns an interface{} and error
// If there is no error interface{} will contain *Client.
// The hpegl provider will put *Client at the value of keyForGLClientMap (returned by ServiceName) in
// the map of clients that it creates and passes down to provider code.  hpegl executes NewClient for each service.
func (i InitialiseClient) NewClient(r *schema.ResourceData) (interface{}, error) {
	var tfprovider models.TFProvider
	if err := tftags.Get(r, &tfprovider); err != nil {
		return nil, err
	}
	// Create VMaas Client
	client := new(Client)

	cfg := api_client.Configuration{
		Host:          serviceURL,
		DefaultHeader: getHeaders(),
		DefaultQueryParams: map[string]string{
			constants.SpaceKey:    tfprovider.Vmaas.SpaceName,
			constants.LocationKey: tfprovider.Vmaas.Location,
		},
	}
	apiClient := api_client.NewAPIClient(&cfg)
	utils.SetMeta(apiClient, r)
	client.CmpClient = cmp_client.NewClient(apiClient, cfg)

	return client, nil
}

// ServiceName is used to return the value of keyForGLClientMap, for use by hpegl
func (i InitialiseClient) ServiceName() string {
	return keyForGLClientMap
}
// GetClientFromMetaMap is a convenience function used by provider code to extract *Client from the
// meta argument passed-in by terraform
func GetClientFromMetaMap(meta interface{}) (*Client, error) {
	cli := meta.(map[string]interface{})[keyForGLClientMap]
	if cli == nil {
		return nil, fmt.Errorf("client is not initialised, make sure that vmaas block is defined in hpegl stanza")
	}

	return cli.(*Client), nil
}

// Get env configurations for VmaaS services
func getHeaders() map[string]string {
	token := os.Getenv("HPEGL_IAM_TOKEN")
	header := make(map[string]string)
	serviceURL = utils.GetServiceEndpoint()
	if utils.GetEnvBool(constants.MockIAMKey) {
		header["subject"] = os.Getenv(constants.CmpSubjectKey)
		header["Authorization"] = token
	}

	return header
}

```

Note the following features of pkg/client:
* We assert that InitaliseClient{} implements the client.Initialisation interface so that the
    compiler will error out if it doesn't

* We use a constant keyForGLClientMap to define the key for the Client struct in the map[string]interface{}
    that is returned by NewClientMap above.  This constant is returned by ServiceName().  Note that hpegl
    will check that the keys returned by service-teams are unique when it starts-up.  It will error out
    if it detects a repeated key.  We suggest that the service mnemonic (e.g. vmaas etc)
    forms part of the key.

* We've added a helper function GetClientFromMetaMap() to extract Client from map[string]interface{}.  This
    function is used in provider code

* The Client struct itself is indicative of what we expect will be used by service code, assuming that
    the service client is generated from a Swagger definition.  Note that the IAM Token is
    contained in the struct.  This Token is generated by the hpegl provider code and passed
    down as part of the provider.ConfigData struct

* The provider.ConfigData struct is defined in
    [hpegl-provider-lib](https://github.com/hewlettpackard/hpegl-provider-lib/blob/bcedebfa36825c4e78b4df73ecaa0e2ed77b92e2/pkg/provider/provider.go#L15).
    This struct is passed into NewClient() and contains all of the hpegl provider configuration
    information, which can be used as necessary in NewClient()

* We have added support for a .gltform file, the relevant code is in pkg/client/config.go.
  At present we are using a .gltform file to provide the IAM token and other information to the GreenLake provider
  when it is run in the context of the POC [Genesis service](https://github.com/hpe-hcss/genesis).  The use of a file
  to deliver information to the terraform provider is a pattern that we've adopted from the stand-alone Metal terraform
  provider.  It is TBD if we will persist with this method as the GreenLake provider is developed further.

* We hope that service teams can use the basic structure of pkg/client, and change the following
    to implement their client:
    * Change the value of keyForGLClientMap - remember that it should contain the service mnemonic
    * Change the contents of Client as appropriate, the actual service client at least will need to be changed
    * Add whatever code is necessary to NewClient()

### pkg/resources

This contains the implementation of the registration.ServiceRegistration interface from hpegl-provider-lib.
This interface is used to associate hcl names with specific service object CRUD code.  The struct
that implements this interface is called Registration and this name should be retained for all service
provider code.  The hpegl provider will iterate over a slice of ServiceRegistration structs from
each supported service to populate the maps that the provider uses to associate hcl object names with
the relevant CRUD code for resources and for data-sources.

This is the slice that is created in the hpegl provider:

```go
package resources

import (
	"github.com/hewlettpackard/hpegl-provider-lib/pkg/registration"

	resvmaas "github.com/HewlettPackard/hpegl-vmaas-terraform-resources/pkg/resources"
)

func SupportedServices() []registration.ServiceRegistration {
	return []registration.ServiceRegistration{
		resvmaas.Registration{},
	}
}
```

The code that iterates over the slice is contained in
[hpe-provider-lib NewProviderFunc](https://github.com/hewlettpackard/hpegl-provider-lib/blob/bcedebfa36825c4e78b4df73ecaa0e2ed77b92e2/pkg/provider/provider.go#L34):

```go
// NewProviderFunc is called from hpegl and service-repos to create a plugin.ProviderFunc which is used
// to define the provider that is exposed to Terraform.  The hpegl repo will use this to create a provider
// that spans all supported services.  A service repo will use this to create a "dummy" provider restricted
// to just the service that can be used for development purposes and for acceptance testing
func NewProviderFunc(reg []registration.ServiceRegistration, pf ConfigureFunc) plugin.ProviderFunc {
	return func() *schema.Provider {
		dataSources := make(map[string]*schema.Resource)
		resources := make(map[string]*schema.Resource)
		for _, service := range reg {
			for k, v := range service.SupportedDataSources() {
				dataSources[k] = v
			}
			for k, v := range service.SupportedResources() {
				resources[k] = v
			}
		}
...

```

We are mandating the following resource and data-source naming structure:

```bash
hpegl_<service mnemonic>_<service resource or data-source name>
```

To be specific at the moment we have the following names:
* VmaaS: <br>
  hpegl_vmaas_network


## internal

The internal directory contains:
* Terraform acceptance tests in acceptance_test
* Service resource terraform CRUD implementation code in resources
* Code in test-utils that returns a plugin.ProviderFunc object that can be used
    to create a "dummy" provider for development purposes (see later).  This object can also be used
    for acceptance tests.


### acceptance_test

There are two basic types of [terraform test](https://www.terraform.io/docs/extend/testing/index.html):
* [Unit tests](https://www.terraform.io/docs/extend/testing/unit-testing.html): <br>
  Terraform unit-tests are, according to the documentation, used mainly for "testing helper methods that
  expand or flatten API responses into data structures for storage into state by Terraform"


* [Acceptance tests](https://www.terraform.io/docs/extend/testing/acceptance-tests/index.html): <br>
  Acceptance tests require real service instances to run against, and exercise CRUD of real service objects.
  These tests are the primary method of ensuring that Terraform providers work as advertised.
  Acceptance tests need to be developed for each service that is added to the GreenLake provider.
  Acceptance tests can be run from this service provider repo using the plugin.ProviderFunc object
  returned from test-utils. Moreover the intention is that we will be able to copy acceptance tests
  verbatim (or with minimal changes to test set-up) from the service repos to hpegl.
  <br><br>
  Some information on writing acceptance tests can be found
  [here](https://www.terraform.io/docs/extend/testing/acceptance-tests/testcase.html)

An example of how to use test-utils.ProviderFunc() to populate a test provider map required by acceptance
tests is contained in provider_test.go

#### Running acceptance tests

To run acceptance tests:

```bash
$ make acceptance
```

### resources

This code is ultimately executed by the hpegl terraform provider.  Note that we are using
the v2.0 SDK.  The Terraform provider writing tutorial [here](https://learn.hashicorp.com/collections/terraform/providers) has been updated
to use the v2.0 SDK.  One of the features of the v2.0 SDK is that
a context is passed-down to the CRUD functions which allows terraform to time-out
operations by cancelling the context.  With regard to setting timeouts, note this map in the cluster resource definition:
```go
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(clusterAvailableTimeout),
			// Update: schema.DefaultTimeout(clusterAvailableTimeout),
			Delete: schema.DefaultTimeout(clusterDeleteTimeout),
		},
```

The timeout settings are used by Terraform itself.  When a timeout is exceeded the context passed-in is cancelled,
and OpenAPI-generated clients will exit.

A diag.Diagnostics slice is returned
by the CRUD functions which is processed by terraform.  Errors and warnings are presented on the console to
the user.  There are helper functions in the diag library to create errors and warnings, see the code
for examples of how to use them.

The resource code makes use of the client.GetClientFromMetaMap() function to extract the Client object from
the meta argument to each CRUD function.  The meta argument is passed-in by terraform, and in the case
of hpegl is a map[string]interface{} of service clients, where each service client defines the key for its
Client in the map.

The bulk of service-team development will occur in this directory.  Please ensure that as resource and data-source
terraform objects are created that they are added to pkg/resources/registration.go

### test-utils

The code is this directory should not need to be changed outside of imports.  It constructs a terraform plugin.ProviderFunc object
that is used to create a "dummy" service-specific terraform provider called hpegl that can be used for
development work.  See later for information on how to build and use this "dummy" provider.

The code uses the registration.ServiceRegistration and client.InitialiseClient interfaces created in pkg/reources
and pkg/client respectively to create the ProviderFunc.  This means that only the resources and data-sources in
addition to the Client defined in this repo are available to the service-specific terraform provider.

The code calls into the hpegl-provider-lib provider.NewProviderFunc like so:

```go
package testutils

import (
	"context"

	"github.com/HewlettPackard/hpegl-vmaas-terraform-resources/pkg/client"
	"github.com/HewlettPackard/hpegl-vmaas-terraform-resources/pkg/resources"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/hewlettpackard/hpegl-provider-lib/pkg/provider"
	"github.com/hewlettpackard/hpegl-provider-lib/pkg/token/common"
	"github.com/hewlettpackard/hpegl-provider-lib/pkg/token/retrieve"
	"github.com/hewlettpackard/hpegl-provider-lib/pkg/token/serviceclient"
)

func ProviderFunc() plugin.ProviderFunc {
	return provider.NewProviderFunc(provider.ServiceRegistrationSlice(resources.Registration{}), providerConfigure)
}

func providerConfigure(p *schema.Provider) schema.ConfigureContextFunc { // nolint staticcheck
	return func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		cli, err := client.InitialiseClient{}.NewClient(d)
		if err != nil {
			return nil, diag.Errorf("error in creating client: %s", err)
		}

		// Initialise token handler
		h, err := serviceclient.NewHandler(d)
		if err != nil {
			return nil, diag.FromErr(err)
		}

		// Returning a map[string]interface{} with the Client from pkg.client at the
		// key specified in that repo to ensure compatibility with the hpegl terraform
		// provider
		return map[string]interface{}{
			client.InitialiseClient{}.ServiceName(): cli,
			common.TokenRetrieveFunctionKey:         retrieve.NewTokenRetrieveFunc(h),
		}, nil
	}
}
```

Note the following:
* We use a hpegl-provider-lib convenience function to convert resources.Registration from pkg/resources into
    a slice

* We run provider.GetConfigData() to generate the provider.ConfigData structure needed by NewClient

* We create a map[string]interface{} with just one entry at the key specified by pkg/client that contains
    the pkg/client Client{} created for compatibility with the resource code

* The "dummy" service-specific provider that is created with this function has all of the provider configuration
    information that is required by the real hpegl provider, which calls-into the same NewProviderFunc()

## Building and using the "dummy" service-specific provider

The code that exposes the plugin.ProviderFunc object created in internal/test-utils to terraform as a provider
is contained in cmd/terraform-provider-hpegl and should not need to be altered outside of imports.

The service-specific provider is exposed to terraform with the name "hpegl", the same as that of the overall GreenLake
provider.  This means that the resource definitions in pkg/resources do not have to be changed in any way
when using the service-specific provider in hcl for provider development or in acceptance tests.

### Building the service-specific provider

To build the provider type "make install".  This will build the provider binary and also place it
in a .local directory that is compatible with terraform versions >= v0.13.  It is important to note
that the name of the service under development should be included in the .local directory path.  Note the
following section at the head of the Makefile:

```makefile
NAME=$(shell find cmd -name "main.go" -exec dirname {} \; | sort -u | sed -e 's|cmd/||')
VERSION=0.0.1
DUMMY_PROVIDER=vmaas
LOCAL_LOCATION=~/.local/share/terraform/plugins/terraform.example.com/$(DUMMY_PROVIDER)/hpegl/$(VERSION)/linux_amd64/
```


See below for how to use this service-specific provider in development.

### Using the service-specific provider

The service specific provider will be exposed to terraform under the name "hpegl".  We need to make sure
that terraform will use the provider stored in the DUMMY_PROVIDER location specified in the Makefile.
To do this note main.tf in examples, in particular the terraform stanza:

```hcl
# Set-up for terraform >= v0.13
terraform {
  required_providers {
    hpegl = {
      # We are specifying a location that is specific to the service under development
      # In this example it is vmaas (see "source" below).  The service-specific replacement
      # to vmaas must be specified in "source" below and also in the Makefile as the
      # value of DUMMY_PROVIDER.
      source  = "terraform.example.com/vmaas/hpegl"
      version = ">= 0.0.1"
    }
  }
}
```

Note the value of "source" under hpegl.  It needs to be of the form:
```bash
"terraform.example.com/$(DUMMY_PROVIDER)/hpegl"
```
where $(DUMMY_PROVIDER) is replaced by the value of DUMMY_PROVIDER from the Makefile.  With this
"source" terraform will use the service-specific provider that was generated by "make install".
The "hpegl" provider stanza can be provided and service resources and data-sources can be referred to by
the keys specified in pkg/resources in hcl (i.e. .tf) files for development and testing purposes.
