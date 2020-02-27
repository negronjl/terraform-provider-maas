package entity_test

import (
	"net"
	"testing"

	"github.com/google/go-cmp/cmp"

	. "github.com/roblox/terraform-provider-maas/pkg/maas/entity"
	"github.com/roblox/terraform-provider-maas/test/helper"
)

var sampleSubnet Subnet = Subnet{
	Name: "172.16.5.0/24",
	VLAN: VLAN{
		VID:           0,
		MTU:           1500,
		DHCPOn:        false,
		FabricID:      0,
		SecondaryRack: "76y7pg",
		ID:            5001,
		Fabric:        "fabric-0",
		Name:          "untagged",
		Space:         "management",
		PrimaryRack:   "7xtf67",
		ResourceURI:   "/MAAS/api/2.0/vlans/5001/",
	},
	CIDR:            "172.16.5.0/24",
	RDNSMode:        2,
	DNSServers:      []net.IP{},
	AllowDNS:        true,
	AllowProxy:      true,
	ActiveDiscovery: false,
	Managed:         true,
	ID:              9,
	Space:           "management",
	ResourceURI:     "/MAAS/api/2.0/subnets/9/",
}

var sampleSubnets []Subnet = []Subnet{
	Subnet{
		Name: "name-rLI3eq",
		VLAN: VLAN{
			VID:           0,
			MTU:           1500,
			DHCPOn:        false,
			ExternalDHCP:  "",
			RelayVLAN:     0,
			SecondaryRack: "76y7pg",
			FabricID:      0,
			Space:         "management",
			Fabric:        "fabric-0",
			ID:            5001,
			Name:          "untagged",
			PrimaryRack:   "7xtf67",
			ResourceURI:   "/MAAS/api/2.0/vlans/5001/",
		},
		CIDR:      "172.16.1.0/24",
		RDNSMode:  2,
		GatewayIP: net.ParseIP("172.16.1.1"),
		DNSServers: []net.IP{
			net.ParseIP("fd89:8724:81f1:5512:557f:99c3:6967:8d63"),
		},
		AllowDNS:        true,
		AllowProxy:      true,
		ActiveDiscovery: false,
		Managed:         true,
		Space:           "management",
		ID:              1,
		ResourceURI:     "/MAAS/api/2.0/subnets/1/",
	},
	Subnet{
		Name: "name-v5djzQ",
		VLAN: VLAN{
			VID:           0,
			MTU:           1500,
			DHCPOn:        false,
			ExternalDHCP:  "",
			RelayVLAN:     0,
			SecondaryRack: "76y7pg",
			FabricID:      1,
			Space:         "management",
			Fabric:        "fabric-1",
			ID:            5003,
			Name:          "untagged",
			PrimaryRack:   "7xtf67",
			ResourceURI:   "/MAAS/api/2.0/vlans/5003/",
		},
		CIDR:      "172.16.2.0/24",
		RDNSMode:  2,
		GatewayIP: net.ParseIP("172.16.2.1"),
		DNSServers: []net.IP{
			net.ParseIP("fcb0:c682:8c15:817d:7d80:2713:e225:5624"),
			net.ParseIP("fd66:86c9:6a50:27cd:de13:3f1c:40d1:8aac"),
			net.ParseIP("120.129.237.29"),
		},
		AllowDNS:        true,
		AllowProxy:      true,
		ActiveDiscovery: false,
		Managed:         true,
		Space:           "management",
		ID:              2,
		ResourceURI:     "/MAAS/api/2.0/subnets/2/",
	},
	Subnet{
		Name: "name-zznp45",
		VLAN: VLAN{
			VID:           10,
			MTU:           1500,
			DHCPOn:        false,
			ExternalDHCP:  "",
			RelayVLAN:     0,
			SecondaryRack: "76y7pg",
			FabricID:      0,
			Space:         "internal",
			Fabric:        "fabric-0",
			ID:            5002,
			Name:          "10",
			PrimaryRack:   "7xtf67",
			ResourceURI:   "/MAAS/api/2.0/vlans/5002/",
		},
		CIDR:      "172.16.3.0/24",
		RDNSMode:  2,
		GatewayIP: net.ParseIP("172.16.3.1"),
		DNSServers: []net.IP{
			net.ParseIP("fd98:8601:90d0:c8c:dd2e:ba51:fa5a:dcfa"),
			net.ParseIP("11.209.150.208"),
			net.ParseIP("fde6:f9ef:3ee9:c5de:2a66:1582:cc83:abaf"),
		},
		AllowDNS:        true,
		AllowProxy:      true,
		ActiveDiscovery: false,
		Managed:         true,
		Space:           "internal",
		ID:              3,
		ResourceURI:     "/MAAS/api/2.0/subnets/3/",
	},
	Subnet{
		Name: "name-c2ULe1",
		VLAN: VLAN{
			VID:           10,
			MTU:           1500,
			DHCPOn:        false,
			ExternalDHCP:  "",
			RelayVLAN:     0,
			SecondaryRack: "76y7pg",
			FabricID:      0,
			Space:         "internal",
			Fabric:        "fabric-0",
			ID:            5002,
			Name:          "10",
			PrimaryRack:   "7xtf67",
			ResourceURI:   "/MAAS/api/2.0/vlans/5002/",
		},
		CIDR:      "172.16.4.0/24",
		RDNSMode:  2,
		GatewayIP: net.ParseIP("172.16.4.1"),
		DNSServers: []net.IP{
			net.ParseIP("fd08:fef7:5c1f:a2e6:3d8e:6c3b:89f9:80cb"),
			net.ParseIP("fc67:ad6a:88fe:9192:62f9:e882:8bcc:339e"),
			net.ParseIP("255.59.162.158"),
		},
		AllowDNS:        true,
		AllowProxy:      true,
		ActiveDiscovery: false,
		Managed:         true,
		Space:           "internal",
		ID:              4,
		ResourceURI:     "/MAAS/api/2.0/subnets/4/",
	},
	Subnet{
		Name: "name-m3vYqT",
		VLAN: VLAN{
			VID:           42,
			MTU:           1500,
			DHCPOn:        false,
			ExternalDHCP:  "",
			RelayVLAN:     0,
			SecondaryRack: "",
			FabricID:      1,
			Space:         "ipv6-testbed",
			Fabric:        "fabric-1",
			ID:            5004,
			Name:          "42",
			PrimaryRack:   "",
			ResourceURI:   "/MAAS/api/2.0/vlans/5004/",
		},
		CIDR:      "2001:db8:42::/64",
		RDNSMode:  2,
		GatewayIP: net.IP{},
		DNSServers: []net.IP{
			net.ParseIP("fd15:6cb0:a55c:235f:e78f:ba4f:2eb4:6b3"),
			net.ParseIP("fcc5:8b5e:c55b:90e0:8be:6b87:eb5:f4c7"),
		},
		AllowDNS:        true,
		AllowProxy:      true,
		ActiveDiscovery: false,
		Managed:         true,
		Space:           "ipv6-testbed",
		ID:              5,
		ResourceURI:     "/MAAS/api/2.0/subnets/5/",
	},
	Subnet{
		Name: "192.168.122.0/24",
		VLAN: VLAN{
			VID:           0,
			MTU:           1500,
			DHCPOn:        false,
			ExternalDHCP:  "",
			RelayVLAN:     0,
			SecondaryRack: "76y7pg",
			FabricID:      0,
			Space:         "management",
			Fabric:        "fabric-0",
			ID:            5001,
			Name:          "untagged",
			PrimaryRack:   "7xtf67",
			ResourceURI:   "/MAAS/api/2.0/vlans/5001/",
		},
		CIDR:            "192.168.122.0/24",
		RDNSMode:        2,
		GatewayIP:       net.ParseIP("192.168.122.1"),
		DNSServers:      []net.IP{},
		AllowDNS:        true,
		AllowProxy:      true,
		ActiveDiscovery: false,
		Managed:         true,
		Space:           "management",
		ID:              6,
		ResourceURI:     "/MAAS/api/2.0/subnets/6/",
	},
	Subnet{
		Name: "172.16.99.0/24",
		VLAN: VLAN{
			VID:           0,
			MTU:           1500,
			DHCPOn:        false,
			ExternalDHCP:  "",
			RelayVLAN:     0,
			SecondaryRack: "76y7pg",
			FabricID:      0,
			Space:         "management",
			Fabric:        "fabric-0",
			ID:            5001,
			Name:          "untagged",
			PrimaryRack:   "7xtf67",
			ResourceURI:   "/MAAS/api/2.0/vlans/5001/",
		},
		CIDR:            "172.16.99.0/24",
		RDNSMode:        2,
		GatewayIP:       net.IP{},
		DNSServers:      []net.IP{},
		AllowDNS:        true,
		AllowProxy:      true,
		ActiveDiscovery: false,
		Managed:         true,
		Space:           "management",
		ID:              7,
		ResourceURI:     "/MAAS/api/2.0/subnets/7/",
	},
}

func TestSubnet(t *testing.T) {
	subnet := new(Subnet)
	subnets := new([]Subnet)

	// Unmarshal sample data into the types
	if err := helper.TestdataFromJSON("maas/subnet.json", subnet); err != nil {
		t.Fatal(err)
	}
	if err := helper.TestdataFromJSON("maas/subnets.json", subnets); err != nil {
		t.Fatal(err)
	}

	// Verify the values are correct
	if diff := cmp.Diff(&sampleSubnet, subnet); diff != "" {
		t.Fatalf("json.Decode(Subnet) mismatch (-want +got):\n%s", diff)
	}
	if diff := cmp.Diff(&sampleSubnets, subnets); diff != "" {
		t.Fatalf("json.Decode([]Subnet) mismatch (-want +got):\n%s", diff)
	}
}
