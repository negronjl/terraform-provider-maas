package entity_test

import (
	"testing"

	. "github.com/roblox/terraform-provider-maas/pkg/maas/entity"
	"github.com/roblox/terraform-provider-maas/test/helper"
)

func TestMachinet(t *testing.T) {
	machine := new(Machine)
	machines := new([]Machine)

	// Unmarshal sample data into the types
	if err := helper.TestdataFromJSON("maas/machine.json", machine); err != nil {
		t.Fatal(err)
	}
	if err := helper.TestdataFromJSON("maas/machines.json", machines); err != nil {
		t.Fatal(err)
	}
}
