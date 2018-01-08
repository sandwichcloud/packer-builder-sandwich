package main

import (
	"github.com/hashicorp/packer/packer/plugin"
	"github.com/sandwichcloud/packer-builder-sandwich/sandwich"
)

func main() {
	server, err := plugin.Server()
	if err != nil {
		panic(err)
	}
	server.RegisterBuilder(new(sandwich.Builder))
	server.Serve()
}