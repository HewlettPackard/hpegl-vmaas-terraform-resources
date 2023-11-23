# Running and Writing Acceptance Tests for VMaaS
- [Running an Acceptance Test](#running-an-acceptance-test)
- [Writing an Acceptance Test](#writing-an-acceptance-test)
    - [Writing Acceptance Test for New Resource/Datasource](#writing-acceptance-test-for-new-resource-datasource)
- [Writing Acceptance Test Suite](#writing-acceptance-test-suite)
    - [Acceptance Test Suite folder Structure](#acceptance-test-suite-folder-structure)

## Running an Acceptance Test

Acceptance tests can be run using the `acceptance` target in the Terraform
`Makefile`. Prior to running acceptance test you may need to set provider
IAM configuration as environment.

```bash
export HPEGL_IAM_SERVICE_URL=<iam_service_url>
export HPEGL_TENANT_ID=<tenant_id>
export HPEGL_USER_SECRET=<user_secret>
export HPEGL_USER_ID=<user_id>
export HPEGL_VMAAS_LOCATION=<vmaas_location>
export HPEGL_VMAAS_SPACE_NAME=<vmaas_space_name>
```

By default the Terraform Provider will use the VMaaS production endpoint. To
use a non-production endpoint the `HPEGL_VMAAS_API_URL` must be set. For
example:
```bash
export HPEGL_VMAAS_API_URL="https://iac-vmaas.intg.hpedevops.net"
```

By Default acceptance test will run test cases from folder `acc-testcases`.
If you need to specify different test suite or folder, please set environment
`TF_ACC_TEST_PATH`

For example:
```bash
export TF_ACC_TEST_PATH=$PWD/mytestcases
make acceptance
```
above script will run acceptance test from specified folder and it is expected that
you will cover all the test cases there.


## Writing an Acceptance Test

Acceptance test should be written using vmaas acceptance test framework or `atf`. This
framework supports writing test suites for both resources and data sources.

### Writing Acceptance Test for New Resource/Datasource

Create new test file wil the format `<resource/data_source>_<name of the resource/data source without hpegl_vmaas_>_test.go`.
Exmaple:
- `resource_instance_test.go` (for resource `hpegl_vmaas_instance`).
- `data_source_network_test.go` (for data source `hpegl_vmaas_network`)
You only need to create a new file if there is no test file exists. Otherwise you can use the existing file
and add your test suite there.

If you need to create entire new test suite then create new Test function as well.
Please follow the naming convention for function name as `TestAcc<DataSource/Resource><Plan/Create>`.

General example to create acceptance test
```go
// This example is for acceptance test for plan data source
func TestAccResourceMyResourceCreate(t *testing.T) {
	acc := &atf.Acc{
		PreCheck:     testAccPreCheck,
		Providers:    testAccProviders,
		ResourceName: "hpegl_vmaas_my_resource",
        Version:      "case_2",     // Use version if you need to run multiple step cases, if not then skip this field
                                    // Version field should also be unique
		GetAPI: func(attr map[string]string) (interface{}, error) {
			cl, cfg := getAPIClient()
			iClient := api_client.PlansAPIService{
				Client: cl,
				Cfg:    cfg,
			}
			id := toInt(attr["id"])

			return iClient.GetSpecificMyResource(getAccContext(), id)
		},
	}

	acc.RunDataSourceTests(t)
}
```

`GetAPI` function should contains Get Specific API for the a resource/Data source. This is
used for validation of the specific resource/data source.

## Writing Acceptance Test Suite

Acceptance test suite should be written over default location `acc-testcases` if default acceptance
test suite is used, or else add test case in `TF_ACC_TEST_PATH` folder.

### Acceptance Test Suite folder Structure

TF_ACC_TEST_PATH folder structure as follows

    .
    ├── resource
    │    ├─ my_resource.yaml
    │    ├─ my_resource-case_2.yaml
    │    ├─ ...
    │    └─ .
    └── data-sources
         ├─ my_datasource.yaml
         ├─ ...
         └─ .

The filename of the test should be the resource/data_source name without "hpegl_vmaas_" string.
for example, resource `hpegl_vmaas_instance` filename will be `instance.yaml`.

Also you may needto append `Version` (if you specified one in the test case on corresponding resource test go file)
prepended with `-`.
for example, resource `hpegl_vmaas_instance` with `Version: "case2"` will be `instance-case2.yaml`.

### Creating new test suite

To create entire new test suite, where you need to create a resource (or run new Plan for data source/Resource),
optionally apply updates on resource, create new file with appended `Version`(if a test suite already presented)
along with filename (In this case you also need to create new test function in corresponding go test file with `Version`),

Structure of the test suite yaml file will be
```yml
acc: # this should be parent field for every test suite
    - name: <name of the test suite> # you can specify name of the test suite here, Apart from the readability
                                     # there are no real world use for the 'name' field
      config: |-                    # Config holds the terraform actual configuration.
        <
            First configuration in test suite will create the resource and rest will update the resource.
            If `expect_error` is specified in first config, then second configuration will be used to
            create a new resource.
            terraform configuration. Here you don't want to
            specify resource name or local name, but only the fields under it.
            for example:
            ```
            resource "hpegl_vmaas_my_resource" "tf_resource" {
                name = "resource_name"
                tag {
                    tag_name = "test"
                }
            }
            ```
            above configuration in test case should represent as follows

                config:
                    name = "resource_name"
                    tag {
                        tag_name = "test"
                    }

            Acceptance test framework will prepend the resource name and local name along
            with this >

      validation:               # validation can contains n number of child validations
        json.instance.status: "running"     # Currently json or tf validations are supported
                                            # in json validation, atf framework will call get API for instance with ID of the created instance
                                            # and get instance.status field from response json and check equality to "running"
        tf.network.0.is_primary: true    # in tf validation corresponding field in state field is compared to RHS
      expect_error: <regex expression if expecting error from above configuration>

    - name: <name of the test suite>
      config: |-
        < this configuration will result into update operation.
          Here you need to provide entire configuration with updated fields
          For example, to change name of the resource from "resource_name" to "resource_name_2"
          the config looks as follows

            config:
                name = "resource_name_2"
                tag {
                    tag_name = "test"
                }
      validation: # here you can validate updated fields as well as existing fields
        ...
    ...

ignore: <if true then entire test suite in this file will be skipped>
```
