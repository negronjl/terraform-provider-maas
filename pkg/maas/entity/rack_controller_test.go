package entity_test

import (
	"testing"

	. "github.com/roblox/terraform-provider-maas/pkg/maas/entity"
	"github.com/roblox/terraform-provider-maas/test/helper"
)

func TestRackControllert(t *testing.T) {
	rackController := new(RackController)
	rackControllers := new([]RackController)

	// Unmarshal sample data into the types
	if err := helper.TestdataFromJSON("maas/rack_controller.json", rackController); err != nil {
		t.Fatal(err)
	}
	if err := helper.TestdataFromJSON("maas/rack_controllers.json", rackControllers); err != nil {
		t.Fatal(err)
	}
}
