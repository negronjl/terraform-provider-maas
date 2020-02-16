package subnet_test

import (
	"net"
	"testing"

	"github.com/google/go-cmp/cmp"
	. "github.com/roblox/terraform-provider-maas/pkg/maas/entity/subnet"
	"github.com/roblox/terraform-provider-maas/test/helper"
)

var sampleStatistics Statistics = Statistics{
	NumAvailable:     232,
	LargestAvailable: 41,
	NumUnavailable:   22,
	TotalAddresses:   254,
	Usage:            0.08661417322834646,
	UsageString:      "9%",
	AvailableString:  "91%",
	FirstAddress:     net.ParseIP("172.16.1.1"),
	LastAddress:      net.ParseIP("172.16.1.254"),
	IPVersion:        4,
}

func TestStatistics(t *testing.T) {
	stats := new(Statistics)

	// Unmarshal sample data into the types
	if err := helper.TestdataFromJSON("maas/subnets/statistics.json", stats); err != nil {
		t.Fatal(err)
	}

	// Verify the values are correct
	if diff := cmp.Diff(&sampleStatistics, stats); diff != "" {
		t.Fatalf("json.Decode(Statistics) mismatch (-want +got):\n%s", diff)
	}
}
