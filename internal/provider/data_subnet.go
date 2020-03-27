package provider

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/juju/gomaasapi"
	"github.com/roblox/terraform-provider-maas/pkg/gmaw"
	"github.com/roblox/terraform-provider-maas/pkg/maas/entity"
)

// DataSubnet provides a lookup for a MaaS Subnet
func DataSubnet() *schema.Resource {
	return &schema.Resource{
		Read: dataSubnetRead,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"vlan": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"cidr": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"rdns_mode": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"gateway_ip": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"dns_servers": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSubnetRead(d *schema.ResourceData, m interface{}) error {
	mo := m.(*gomaasapi.MAASObject)
	res, err := gmaw.NewSubnets(mo).Get()
	if err != nil {
		return err
	}

	for idx := range res {
		if !dataSubnetIsMatch(d, &res[idx]) {
			continue
		}
		if err := d.Set("name", res[idx].Name); err != nil {
			return err
		}
		if err := d.Set("vlan", res[idx].VLAN.ID); err != nil {
			return err
		}
		if err := d.Set("cidr", res[idx].CIDR); err != nil {
			return err
		}
		if err := d.Set("rdns_mode", res[idx].RDNSMode); err != nil {
			return err
		}
		if err := d.Set("gateway_ip", res[idx].GatewayIP); err != nil {
			return err
		}
		if err := d.Set("dns_servers", res[idx].DNSServers); err != nil {
			return err
		}
		d.SetId(string(res[idx].ID))
		return nil
	}
	return fmt.Errorf("could not find matching subnet")
}

func dataSubnetIsMatch(d *schema.ResourceData, res *entity.Subnet) bool {
	name, vlan, cidr := d.Get("name").(string), d.Get("vlan").(int), d.Get("cidr").(string)
	if !(name == "" || name == res.Name) {
		return false
	}
	if !(vlan == 0 || vlan == res.VLAN.ID) {
		return false
	}
	return (cidr == "" || cidr == res.CIDR)
}
