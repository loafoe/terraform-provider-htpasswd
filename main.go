package main

import (
	"context"
	"flag"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/loafoe/terraform-provider-htpasswd/htpasswd"
)

var version = "dev"

func main() {
	var debugMode bool

	flag.BoolVar(&debugMode, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := providerserver.ServeOpts{
		Address: "registry.terraform.io/loafoe/htpasswd",
		Debug:   debugMode,
	}

	err := providerserver.Serve(context.Background(), htpasswd.New(version), opts)
	if err != nil {
		log.Fatal(err.Error())
	}
}
