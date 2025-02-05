// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package main

import (
	"context"
	"flag"
	"log"

	"terraform-provider-tasklite/internal/provider"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

var (
	// these will be set by the goreleaser configuration
	// to appropriate values for the compiled binary.
	version string = "dev"
)

func main() {
	var debug bool

	flag.BoolVar(&debug, "debug", true, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := providerserver.ServeOpts{
		Address: "registry.terraform.io/providers/tasklite",
		Debug:   debug,
	}

	err := providerserver.Serve(context.Background(), provider.New(version), opts)

	if err != nil {
		log.Fatal(err.Error())
	}
}
