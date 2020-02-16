package gmaw

import (
	"encoding/json"
	"net/url"
	"strconv"

	"github.com/juju/gomaasapi"
	"github.com/roblox/terraform-provider-maas/pkg/api/params"
	"github.com/roblox/terraform-provider-maas/pkg/maas"
	"github.com/roblox/terraform-provider-maas/pkg/maas/entity"
	"github.com/roblox/terraform-provider-maas/pkg/maas/entity/subnet"
)

// Subnet provides methods for the Subnet operations in the MaaS API.
// This type should be instantiated via NewSubnet(). It fulfills the
// api.Subnet interface.
type Subnet struct {
	c Client
}

// NewSubnet configures a new Subnet.
func NewSubnet(client *gomaasapi.MAASObject) *Subnet {
	c := client.GetSubObject("subnets")
	return &Subnet{c: Client{&c}}
}

// client returns a Client (ie wrapped MAASOBject) for the subnet with the given ID
func (s *Subnet) client(id int) Client {
	return s.c.GetSubObject(strconv.Itoa(id))
}

// Delete removes a subnet.
// This function returns an error if the gomaasapi returns an error.
func (s *Subnet) Delete(id int) error {
	return s.client(id).Delete()
}

// Get returns information about a subnet.
// This function returns an error if the gomaasapi returns an error or if
// the response cannot be decoded.
func (s *Subnet) Get(id int) (sn *entity.Subnet, err error) {
	sn = new(entity.Subnet)
	err = s.client(id).Get("", url.Values{}, func(data []byte) error {
		return json.Unmarshal(data, sn)
	})
	return
}

// GetIPAddresses returns information about the IP addresses in the subnet.
// The withUsername and withSummary parameters can be set to false to suppress
// the display of usernames and nodes/BMCs/DNSRRs, respectively, associated with
// each IP. This function returns an error if the gomaasapi returns an error or
// if the response cannot be decoded.
func (s *Subnet) GetIPAddresses(id int, withUsername, withSummary bool) (addrs []subnet.IPAddress, err error) {
	v := make(url.Values)
	if !withUsername {
		v.Set("with_username", "0")
	}
	if !withSummary {
		v.Set("with_summary", "0")
	}
	err = s.client(id).Get("ip_addresses", v, func(data []byte) error {
		return json.Unmarshal(data, &addrs)
	})
	return
}

// GetReservedIPRanges returns a list of reserved IP ranges.
// This function returns an error if the gomaasapi returns an error or if
// the response cannot be decoded.
func (s *Subnet) GetReservedIPRanges(id int) (ranges []subnet.ReservedIPRange, err error) {
	err = s.client(id).Get("reserved_ip_ranges", url.Values{}, func(data []byte) error {
		return json.Unmarshal(data, &ranges)
	})
	return
}

// GetStatistics returns some statistics about the subnet.
// If includeRanges is true, the response will include detailed information
// about the usage of the range. Setting includeSuggestions includes the
// suggested gateway and dynamic range for this subnet, if it were
// to be configured.
// This function returns an error if the gomaasapi returns an error or if
// the response cannot be decoded.
func (s *Subnet) GetStatistics(id int, includeRanges, includeSuggestions bool) (stats *subnet.Statistics, err error) {
	stats = new(subnet.Statistics)
	v := make(url.Values)
	if includeRanges {
		v.Set("include_ranges", "1")
	}
	if includeSuggestions {
		v.Set("include_suggestions", "1")
	}
	err = s.client(id).Get("statistics", v, func(data []byte) error {
		return json.Unmarshal(data, stats)
	})
	return
}

// GetUnreservedIPRanges returns a list of unreserved IP ranges.
// This function returns an error if the gomaasapi returns an error or if
// the response cannot be decoded.
func (s *Subnet) GetUnreservedIPRanges(id int) (ranges []subnet.IPRange, err error) {
	err = s.client(id).Get("unreserved_ip_ranges", url.Values{}, func(data []byte) error {
		return json.Unmarshal(data, &ranges)
	})
	return
}

// Put updates the configuration for a subnet.
// This function returns an error if the gomaasapi returns an error or if
// the response cannot be decoded.
func (s *Subnet) Put(id int, p *params.Subnet) (sn *entity.Subnet, err error) {
	qsp := maas.ToQSP(p)
	sn = new(entity.Subnet)
	err = s.client(id).Put(qsp, func(data []byte) error {
		return json.Unmarshal(data, sn)
	})
	return
}
