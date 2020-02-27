package entity_test

import (
	"testing"

	. "github.com/roblox/terraform-provider-maas/pkg/maas/entity"
	"github.com/roblox/terraform-provider-maas/test/helper"
)

func TestZonet(t *testing.T) {
	zone := new(Zone)
	if err := helper.TestdataFromJSON("maas/zone.json", zone); err != nil {
		t.Fatal(err)
	}
}
