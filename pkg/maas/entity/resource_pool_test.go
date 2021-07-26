package entity_test

import (
	"testing"

	. "github.com/roblox/terraform-provider-maas/pkg/maas/entity"
	"github.com/roblox/terraform-provider-maas/test/helper"
)

func TestResourcePoolt(t *testing.T) {
	pool := new(ResourcePool)
	pools := new([]ResourcePool)

	// Unmarshal sample data into the types
	if err := helper.TestdataFromJSON("maas/resource_pool.json", pool); err != nil {
		t.Fatal(err)
	}
	if err := helper.TestdataFromJSON("maas/resource_pools.json", pools); err != nil {
		t.Fatal(err)
	}
}
