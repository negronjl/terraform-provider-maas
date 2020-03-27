package provider

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/juju/gomaasapi"
	"github.com/roblox/terraform-provider-maas/pkg/api/params"
	"github.com/roblox/terraform-provider-maas/pkg/gmaw"
)

// DataRackController provides a lookup for MaaS Rack Controllers
func DataRackController() *schema.Resource {
	return &schema.Resource{
		Read: dataRackControllerRead,

		Schema: map[string]*schema.Schema{
			"hostname": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"mac_address": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"system_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"domain": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"zone": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"pool": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"agent_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataRackControllerRead(d *schema.ResourceData, m interface{}) error {
	mo := m.(*gomaasapi.MAASObject)
	criteria := &params.RackControllerSearch{
		Hostname:   d.Get("hostname").(string),
		MACAddress: d.Get("mac_address").(string),
		SystemID:   d.Get("system_id").(string),
		Domain:     d.Get("domain").(string),
		Zone:       d.Get("zone").(string),
		Pool:       d.Get("pool").(string),
		AgentName:  d.Get("agent_name").(string),
	}
	ctrls, err := gmaw.NewRackControllers(mo).Get(criteria)
	if err != nil {
		return err
	}
	if len(ctrls) == 0 {
		return fmt.Errorf("no matching rack controllers found")
	}

	if err := d.Set("hostname", ctrls[0].Hostname); err != nil {
		return err
	}
	if err := d.Set("system_id", ctrls[0].SystemID); err != nil {
		return err
	}
	if err := d.Set("domain", ctrls[0].Domain.Name); err != nil {
		return err
	}
	if err := d.Set("zone", ctrls[0].Zone.Name); err != nil {
		return err
	}
	if err := d.Set("pool", ctrls[0].Pool.Name); err != nil {
		return err
	}
	d.SetId(ctrls[0].SystemID)
	return nil
}
