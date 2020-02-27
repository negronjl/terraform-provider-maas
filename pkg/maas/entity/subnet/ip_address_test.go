package subnet_test

import (
	"net"
	"testing"

	"github.com/google/go-cmp/cmp"
	. "github.com/roblox/terraform-provider-maas/pkg/maas/entity/subnet"
	"github.com/roblox/terraform-provider-maas/test/helper"
)

var sampleIPAddresses []IPAddress = []IPAddress{
	IPAddress{
		IP:        net.ParseIP("172.16.2.3"),
		AllocType: 1,
		Created:   "Tue, 27 Nov. 2018 18:06:19",
		Updated:   "Tue, 27 Nov. 2018 18:06:19",
		NodeSummary: NodeSummary{
			SystemID:    "76y7pg",
			NodeType:    2,
			FQDN:        "happy-rack.maas",
			Hostname:    "happy-rack",
			IsContainer: false,
			Via:         "eth2",
		},
	},
	IPAddress{
		IP:        net.ParseIP("172.16.2.4"),
		AllocType: 1,
		Created:   "Tue, 27 Nov. 2018 18:06:19",
		Updated:   "Tue, 27 Nov. 2018 18:06:19",
		NodeSummary: NodeSummary{
			SystemID:    "nfkend",
			NodeType:    3,
			FQDN:        "happy-region.maas",
			Hostname:    "happy-region",
			IsContainer: false,
			Via:         "eth2",
		},
	},
	IPAddress{
		IP:        net.ParseIP("172.16.2.11"),
		AllocType: 1,
		Created:   "Tue, 27 Nov. 2018 18:07:25",
		Updated:   "Tue, 27 Nov. 2018 18:07:25",
		NodeSummary: NodeSummary{
			SystemID:    "dq3sda",
			NodeType:    0,
			FQDN:        "kind-dory.maas",
			Hostname:    "kind-dory",
			IsContainer: false,
			Via:         "eth-xnY2lB",
		},
		User: "user2",
	},
	IPAddress{
		IP:        net.ParseIP("172.16.2.26"),
		AllocType: 1,
		Created:   "Tue, 27 Nov. 2018 18:07:09",
		Updated:   "Tue, 27 Nov. 2018 18:07:09",
		NodeSummary: NodeSummary{
			SystemID:    "seebkg",
			NodeType:    0,
			FQDN:        "grown-cougar.maas",
			Hostname:    "grown-cougar",
			IsContainer: false,
			Via:         "eth-GZF5ig",
		},
		User: "user2",
	},
	IPAddress{
		IP:        net.ParseIP("172.16.2.62"),
		AllocType: 1,
		Created:   "Tue, 27 Nov. 2018 18:07:24",
		Updated:   "Tue, 27 Nov. 2018 18:07:24",
		NodeSummary: NodeSummary{
			SystemID:    "r7enqt",
			NodeType:    0,
			FQDN:        "nice-gannet.maas",
			Hostname:    "nice-gannet",
			IsContainer: false,
			Via:         "eth-nlvMd2",
		},
		User: "user1",
	},
	IPAddress{
		IP:        net.ParseIP("172.16.2.63"),
		AllocType: 1,
		Created:   "Tue, 27 Nov. 2018 18:06:29",
		Updated:   "Tue, 27 Nov. 2018 18:06:29",
		NodeSummary: NodeSummary{
			SystemID:    "pme7wb",
			NodeType:    0,
			FQDN:        "solid-liger.sample",
			Hostname:    "solid-liger",
			IsContainer: false,
			Via:         "eth-2VCthm",
		},
	},
	IPAddress{
		IP:        net.ParseIP("172.16.2.109"),
		AllocType: 1,
		Created:   "Tue, 27 Nov. 2018 18:07:25",
		Updated:   "Tue, 27 Nov. 2018 18:07:25",
		NodeSummary: NodeSummary{
			SystemID:    "dq3sda",
			NodeType:    0,
			FQDN:        "kind-dory.maas",
			Hostname:    "kind-dory",
			IsContainer: false,
			Via:         "eth-fS18k5",
		},
		User: "user2",
	},
	IPAddress{
		IP:        net.ParseIP("172.16.2.111"),
		AllocType: 1,
		Created:   "Tue, 27 Nov. 2018 18:07:07",
		Updated:   "Tue, 27 Nov. 2018 18:07:07",
		NodeSummary: NodeSummary{
			SystemID:    "ydpcwh",
			NodeType:    0,
			FQDN:        "game-owl.ubnt",
			Hostname:    "game-owl",
			IsContainer: false,
			Via:         "eth-Cyk2jC",
		},
		User: "user2",
	},
	IPAddress{
		IP:        net.ParseIP("172.16.2.116"),
		AllocType: 1,
		Created:   "Tue, 27 Nov. 2018 18:07:10",
		Updated:   "Tue, 27 Nov. 2018 18:07:10",
		NodeSummary: NodeSummary{
			SystemID:    "seebkg",
			NodeType:    0,
			FQDN:        "grown-cougar.maas",
			Hostname:    "grown-cougar",
			IsContainer: false,
			Via:         "eth-0nEEnB",
		},
		User: "user2",
	},
	IPAddress{
		IP:        net.ParseIP("172.16.2.134"),
		AllocType: 1,
		Created:   "Tue, 27 Nov. 2018 18:07:24",
		Updated:   "Tue, 27 Nov. 2018 18:07:24",
		NodeSummary: NodeSummary{
			SystemID:    "r7enqt",
			NodeType:    0,
			FQDN:        "nice-gannet.maas",
			Hostname:    "nice-gannet",
			IsContainer: false,
			Via:         "eth-dMPw46",
		},
		User: "user1",
	},
	IPAddress{
		IP:        net.ParseIP("172.16.2.206"),
		AllocType: 1,
		Created:   "Tue, 27 Nov. 2018 18:07:24",
		Updated:   "Tue, 27 Nov. 2018 18:07:24",
		NodeSummary: NodeSummary{
			SystemID:    "r7enqt",
			NodeType:    0,
			FQDN:        "nice-gannet.maas",
			Hostname:    "nice-gannet",
			IsContainer: false,
			Via:         "eth-6Wz9hw",
		},
		User: "user1",
	},
	IPAddress{
		IP:        net.ParseIP("172.16.2.235"),
		AllocType: 1,
		Created:   "Tue, 27 Nov. 2018 18:07:29",
		Updated:   "Tue, 27 Nov. 2018 18:07:29",
		NodeSummary: NodeSummary{
			SystemID:    "7ghwxs",
			NodeType:    0,
			FQDN:        "alert-lion.sample",
			Hostname:    "alert-lion",
			IsContainer: false,
			Via:         "bond-wnZNKS",
		},
		User: "user2",
	},
	IPAddress{
		IP:        net.ParseIP("172.16.2.252"),
		AllocType: 1,
		Created:   "Tue, 27 Nov. 2018 18:07:32",
		Updated:   "Tue, 27 Nov. 2018 18:07:32",
		NodeSummary: NodeSummary{
			SystemID:    "dnq43f",
			NodeType:    0,
			FQDN:        "square-hornet.maas",
			Hostname:    "square-hornet",
			IsContainer: false,
			Via:         "eth-3GCPHI",
		},
		User: "user2",
	},
}

func TestIPAddress(t *testing.T) {
	addresses := new([]IPAddress)

	// Unmarshal sample data into the types
	if err := helper.TestdataFromJSON("maas/subnets/ipaddresses.json", addresses); err != nil {
		t.Fatal(err)
	}

	// Verify the values are correct
	if diff := cmp.Diff(&sampleIPAddresses, addresses); diff != "" {
		t.Fatalf("json.Decode([]IPAddresses) mismatch (-want +got):\n%s", diff)
	}
}
