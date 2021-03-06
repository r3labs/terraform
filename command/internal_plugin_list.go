// +build !core

//
// This file is automatically generated by scripts/generate-plugins.go -- Do not edit!
//
package command

import (
	azureprovider "github.com/r3labs/terraform/builtin/providers/azure"
	azurermprovider "github.com/r3labs/terraform/builtin/providers/azurerm"

	"github.com/r3labs/terraform/plugin"
	"github.com/r3labs/terraform/terraform"

	//New Provider Builds
	opcprovider "github.com/hashicorp/terraform-provider-opc/opc"

	// Legacy, will remove once it conforms with new structure
	chefprovisioner "github.com/r3labs/terraform/builtin/provisioners/chef"
)

var InternalProviders = map[string]plugin.ProviderFunc{
	"azure":   azureprovider.Provider,
	"azurerm": azurermprovider.Provider,
}

var InternalProvisioners = map[string]plugin.ProvisionerFunc{}

func init() {
	// Legacy provisioners that don't match our heuristics for auto-finding
	// built-in provisioners.
	InternalProvisioners["chef"] = func() terraform.ResourceProvisioner { return new(chefprovisioner.ResourceProvisioner) }

	// New Provider Layouts
	InternalProviders["opc"] = func() terraform.ResourceProvider { return opcprovider.Provider() }
}
