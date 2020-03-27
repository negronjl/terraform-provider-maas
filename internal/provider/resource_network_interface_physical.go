package provider

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/roblox/terraform-provider-maas/internal/bridge"
	"github.com/roblox/terraform-provider-maas/internal/tfschema"
)

// ResourceNetworkInterfacePhysical provides a resource to manage physical interfaces
func ResourceNetworkInterfacePhysical() *schema.Resource {
	return &schema.Resource{
		Create: resourceNetworkInterfacePhysicalCreate,
		Read:   resourceNetworkInterfacePhysicalRead,
		Update: resourceNetworkInterfacePhysicalUpdate,
		Delete: resourceNetworkInterfacePhysicalDelete,

		Schema: map[string]*schema.Schema{
			"system_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"interface_id": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
				ForceNew: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"mac_address": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"tags": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"vlan": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"mtu": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"accept_ra": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"autoconf": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
		},

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceNetworkInterfacePhysicalCreate(d *schema.ResourceData, m interface{}) error {
	ifc := bridge.NewNetworkInterface(m)
	sch := tfschema.NewNetworkInterfacePhysical(d)
	if err := ifc.Create(sch); err != nil {
		return err
	}
	if id, err := sch.GetID(); err == nil {
		d.SetId(id)
	}
	return resourceNetworkInterfacePhysicalRead(d, m)
}

func resourceNetworkInterfacePhysicalRead(d *schema.ResourceData, m interface{}) (err error) {
	ifc := bridge.NewNetworkInterface(m)
	sch := tfschema.NewNetworkInterfacePhysical(d)
	if err = ifc.ReadTo(sch); err == nil {
		err = sch.UpdateResource(d)
	}
	return err
}

func resourceNetworkInterfacePhysicalUpdate(d *schema.ResourceData, m interface{}) error {
	ifc := bridge.NewNetworkInterface(m)
	sch := tfschema.NewNetworkInterfacePhysical(d)
	if err := ifc.UpdateFrom(sch); err != nil {
		return err
	}
	return resourceNetworkInterfacePhysicalRead(d, m)
}

func resourceNetworkInterfacePhysicalDelete(d *schema.ResourceData, m interface{}) (err error) {
	ifc := bridge.NewNetworkInterface(m)
	sch := tfschema.NewNetworkInterfacePhysical(d)
	if err = ifc.Delete(sch); err == nil {
		d.SetId("")
	}
	return
}
