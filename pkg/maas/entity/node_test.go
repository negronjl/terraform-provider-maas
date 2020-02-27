package entity_test

import (
	"testing"

	. "github.com/roblox/terraform-provider-maas/pkg/maas/entity"
	"github.com/roblox/terraform-provider-maas/test/helper"
)

func TestNodet(t *testing.T) {
	node := new(Node)
	nodes := new([]Node)

	// Unmarshal sample data into the types
	if err := helper.TestdataFromJSON("maas/node.json", node); err != nil {
		t.Fatal(err)
	}
	if err := helper.TestdataFromJSON("maas/nodes.json", nodes); err != nil {
		t.Fatal(err)
	}
}
