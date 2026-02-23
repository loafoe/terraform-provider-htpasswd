package htpasswd

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"htpasswd": providerserver.NewProtocol6WithError(New("test")()),
}

func TestProvider(t *testing.T) {
	New("test")()
}
