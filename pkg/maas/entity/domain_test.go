package entity_test

import (
	"testing"

	. "github.com/roblox/terraform-provider-maas/pkg/maas/entity"
	"github.com/roblox/terraform-provider-maas/test/helper"
)

func TestDomaint(t *testing.T) {
	domains := new([]Domain)
	if err := helper.TestdataFromJSON("maas/domains.json", domains); err != nil {
		t.Fatal(err)
	}
}
