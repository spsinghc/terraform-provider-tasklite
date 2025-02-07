# Terraform TaskLite Provider

## Prerequisites
- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.22
- [TechChallengeApp](https://github.com/servian/TechChallengeApp)

## Building The Provider

* Clone the repository
* Enter the repository directory
* Build the provider using the Go `install` command:

```shell
go install
```

## Using the provider

To use the provider, follow the steps below:

1. Update `.terraformrc` file to use `dev_overrides`. If the file `.terraformrc` doesn't exist in the home directory `~`, create one, then add the following code. Change `<PATH>` to the value returned from the `go env GOBIN` command above.

```HCL
provider_installation {
  dev_overrides {
      "registry.terraform.io/providers/tasklite" = "<PATH>"
  }
  # For all other providers, install them directly from their origin provider
  # registries as normal. If you omit this, Terraform will _only_ use
  # the dev_overrides block, and so no other providers will be available.
  direct {}
}
```

2. Initialize the provider in your Terraform configuration:

```HCL
terraform {
  required_providers {
    tasklite = {
      source = "registry.terraform.io/providers/tasklite"
    }
  }
}

provider "tasklite" {
  host = "<HOST>" # replace it with TechChallengeApp api url
}
```

3. Define resources using the provider:

```HCL
resource "tasklite_task" "example" {
    title = "Task title"
}
```

4. Initialize Terraform and apply the configuration:

```HCL
    terraform plan # to check plan
    terraform apply # to apply changes
```

```shell
go install
```

Check the `docs` directory for more information on the provider.

## Developing the Provider
If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `go install`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

To generate or update documentation, run `make generate`.

### Tests
* Use `make test` to run api client unit tests.
* user `make testacc` to run API client unit tests as well as acceptance tests for the resource.

*Note:* Acceptance tests use the mock server, which can be replaced with real API server in the [task_resource_test.go](internal/provider/task_resource_test.go) file. Replace ```go u := server.URL``` with real API server URL.