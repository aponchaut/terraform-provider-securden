// Copyright (c) HashiCorp, Inc.

package main

import (
	"context"
	"flag"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"terraform-provider-securden/internal/provider"
)

var (
	version string = "dev"
)

func main() {
	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := providerserver.ServeOpts{
		// TODO: Update this string with the published name of your provider.
		//	Address: "hashicorp.com/provider/securden",
		Address: "terraform.io/local/securden",
		Debug:   debug,
		//ProtocolVersion: 6,  // v6 implicit by default
	}

	err := providerserver.Serve(context.Background(), provider.Provider(version), opts)
	if err != nil {
		log.Fatal(err.Error())
	}
}

// Acceptance tests -->  See later
//resource.Test(t, resource.TestCase{
//    ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error) {
//        // newProvider is an example function that returns a provider.Provider
//        "examplecloud": providerserver.NewProtocol6WithError(newProvider()),
//    },
//    Steps: []resource.TestStep{/* ... */},
//})
