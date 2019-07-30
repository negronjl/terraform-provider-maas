package provider_test

import (
	"testing"

	. "github.com/roblox/terraform-provider-maas/internal/provider"
)

func TestInstance(t *testing.T) {
	var i interface{} = &Instance{}
	if _, ok := i.(Endpoint); !ok {
		t.Fail()
	}
}
