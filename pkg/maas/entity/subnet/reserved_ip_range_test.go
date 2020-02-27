package subnet_test

import (
	"net"
	"testing"

	"github.com/google/go-cmp/cmp"
	. "github.com/roblox/terraform-provider-maas/pkg/maas/entity/subnet"
	"github.com/roblox/terraform-provider-maas/test/helper"
)

var sampleReservedIPRanges []ReservedIPRange = []ReservedIPRange{
	ReservedIPRange{
		IPRange: IPRange{
			Start:        net.ParseIP("172.16.2.1"),
			End:          net.ParseIP("172.16.2.1"),
			NumAddresses: 1,
		},
		Purpose: []string{"gateway-ip"},
	},
	ReservedIPRange{
		IPRange: IPRange{
			Start:        net.ParseIP("172.16.2.3"),
			End:          net.ParseIP("172.16.2.4"),
			NumAddresses: 2,
		},
		Purpose: []string{"assigned-ip"},
	},
	ReservedIPRange{
		IPRange: IPRange{
			Start:        net.ParseIP("172.16.2.11"),
			End:          net.ParseIP("172.16.2.11"),
			NumAddresses: 1,
		},
		Purpose: []string{"assigned-ip"},
	},
	ReservedIPRange{
		IPRange: IPRange{
			Start:        net.ParseIP("172.16.2.26"),
			End:          net.ParseIP("172.16.2.26"),
			NumAddresses: 1,
		},
		Purpose: []string{"assigned-ip"},
	},
	ReservedIPRange{
		IPRange: IPRange{
			Start:        net.ParseIP("172.16.2.62"),
			End:          net.ParseIP("172.16.2.63"),
			NumAddresses: 2,
		},
		Purpose: []string{"assigned-ip"},
	},
	ReservedIPRange{
		IPRange: IPRange{
			Start:        net.ParseIP("172.16.2.101"),
			End:          net.ParseIP("172.16.2.101"),
			NumAddresses: 1,
		},
		Purpose: []string{"gateway-ip"},
	},
	ReservedIPRange{
		IPRange: IPRange{
			Start:        net.ParseIP("172.16.2.109"),
			End:          net.ParseIP("172.16.2.109"),
			NumAddresses: 1,
		},
		Purpose: []string{"assigned-ip"},
	},
	ReservedIPRange{
		IPRange: IPRange{
			Start:        net.ParseIP("172.16.2.111"),
			End:          net.ParseIP("172.16.2.111"),
			NumAddresses: 1,
		},
		Purpose: []string{"assigned-ip"},
	},
	ReservedIPRange{
		IPRange: IPRange{
			Start:        net.ParseIP("172.16.2.116"),
			End:          net.ParseIP("172.16.2.116"),
			NumAddresses: 1,
		},
		Purpose: []string{"assigned-ip"},
	},
	ReservedIPRange{
		IPRange: IPRange{
			Start:        net.ParseIP("172.16.2.134"),
			End:          net.ParseIP("172.16.2.134"),
			NumAddresses: 1,
		},
		Purpose: []string{"assigned-ip"},
	},
	ReservedIPRange{
		IPRange: IPRange{
			Start:        net.ParseIP("172.16.2.174"),
			End:          net.ParseIP("172.16.2.174"),
			NumAddresses: 1,
		},
		Purpose: []string{"gateway-ip"},
	},
	ReservedIPRange{
		IPRange: IPRange{
			Start:        net.ParseIP("172.16.2.206"),
			End:          net.ParseIP("172.16.2.206"),
			NumAddresses: 1,
		},
		Purpose: []string{"assigned-ip"},
	},
	ReservedIPRange{
		IPRange: IPRange{
			Start:        net.ParseIP("172.16.2.235"),
			End:          net.ParseIP("172.16.2.235"),
			NumAddresses: 1,
		},
		Purpose: []string{"assigned-ip"},
	},
	ReservedIPRange{
		IPRange: IPRange{
			Start:        net.ParseIP("172.16.2.237"),
			End:          net.ParseIP("172.16.2.237"),
			NumAddresses: 1,
		},
		Purpose: []string{"gateway-ip"},
	},
	ReservedIPRange{
		IPRange: IPRange{
			Start:        net.ParseIP("172.16.2.252"),
			End:          net.ParseIP("172.16.2.252"),
			NumAddresses: 1,
		},
		Purpose: []string{"assigned-ip"},
	},
}

func TestReservedIPRange(t *testing.T) {
	ranges := new([]ReservedIPRange)

	// Unmarshal sample data into the types
	if err := helper.TestdataFromJSON("maas/subnets/reservedipranges.json", ranges); err != nil {
		t.Fatal(err)
	}

	// Verify the values are correct
	if diff := cmp.Diff(&sampleReservedIPRanges, ranges); diff != "" {
		t.Fatalf("json.Decode([]ReservedIPRange) mismatch (-want +got):\n%s", diff)
	}
}
