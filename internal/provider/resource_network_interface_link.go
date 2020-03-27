package provider

import (
	"fmt"
	"net"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/roblox/terraform-provider-maas/internal/bridge"
	"github.com/roblox/terraform-provider-maas/internal/tfschema"
)

// ResourceNetworkInterfaceLink provides a resource that can be used to manage links between interfaces
func ResourceNetworkInterfaceLink() *schema.Resource {
	return &schema.Resource{
		Create: resourceNetworkInterfaceLinkCreate,
		Read:   resourceNetworkInterfaceLinkRead,
		Delete: resourceNetworkInterfaceLinkDelete,

		Schema: map[string]*schema.Schema{
			"system_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"interface_id": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"subnet_id": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"mode": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(string)
					if !(v == "AUTO" || v == "DHCP" || v == "STATIC" || v == "LINK_UP") {
						errs = append(errs, fmt.Errorf("%q must be 'AUTO', 'DHCP', 'STATIC', or 'LINK_UP' (got '%s')", key, v))
					}
					return
				},
			},
			"ip_address": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(string)
					if ip := net.ParseIP(v); ip == nil {
						errs = append(errs, fmt.Errorf("%q must be a valid IP address (got '%s')", key, v))
					}
					return
				},
			},
			"force": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"default_gateway": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(string)
					if ip := net.ParseIP(v); ip == nil {
						errs = append(errs, fmt.Errorf("%q must be a valid IP address (got '%s')", key, v))
					}
					return
				},
			},
		},

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceNetworkInterfaceLinkCreate(d *schema.ResourceData, m interface{}) (err error) {
	ifc := bridge.NewNetworkInterface(m)
	sch := tfschema.NewNetworkInterfaceLink(d)
	if err = ifc.LinkSubnet(sch); err != nil {
		return
	}
	if id, err := sch.GetID(); err == nil {
		d.SetId(id)
	}
	return
}

func resourceNetworkInterfaceLinkRead(d *schema.ResourceData, m interface{}) (err error) {
	ifc := bridge.NewNetworkInterface(m)
	sch := tfschema.NewNetworkInterfaceLink(d)
	if err = ifc.ReadLink(sch); err == nil {
		err = sch.UpdateResource(d)
	}
	return err
}

func resourceNetworkInterfaceLinkDelete(d *schema.ResourceData, m interface{}) (err error) {
	ifc := bridge.NewNetworkInterface(m)
	sch := tfschema.NewNetworkInterfaceLink(d)
	if err = ifc.UnlinkSubnet(sch); err == nil {
		d.SetId("")
	}
	return err
}
