package provider_test

import (
	"testing"

	. "github.com/roblox/terraform-provider-maas/internal/provider"
	"github.com/roblox/terraform-provider-maas/internal/tfschema"
)

func TestInstance(t *testing.T) {
	var i interface{} = &Instance{}
	if _, ok := i.(tfschema.Endpoint); !ok {
		t.Fail()
	}
}
