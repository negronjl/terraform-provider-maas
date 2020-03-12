package provider

import (
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/juju/gomaasapi"
	"github.com/roblox/terraform-provider-maas/pkg/gmaw"
)

// ResourceServer manages global MaaS configuration options ala the MaaS Server endpoint
func ResourceServer() *schema.Resource {
	return &schema.Resource{
		Create: resourceServerCreate,
		Read:   resourceServerRead,
		Update: resourceServerUpdate,
		Delete: resourceServerDelete,

		Schema: map[string]*schema.Schema{
			"ntp_servers": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceServerCreate(d *schema.ResourceData, m interface{}) (err error) {
	err = resourceServerUpdate(d, m)
	if err == nil {
		d.SetId(time.Now().Format(time.RFC3339))
	}
	return
}

func resourceServerRead(d *schema.ResourceData, m interface{}) error {
	mo := m.(*gomaasapi.MAASObject)
	client := gmaw.NewMAASServer(mo)
	res, err := client.Get("ntp_servers")
	if err == nil {
		err = d.Set("ntp_servers", strings.Split(res, ","))
	}
	return err
}

func resourceServerUpdate(d *schema.ResourceData, m interface{}) error {
	mo := m.(*gomaasapi.MAASObject)
	client := gmaw.NewMAASServer(mo)
	val := d.Get("ntp_servers").([]string)
	if err := client.Post("ntp_servers", strings.Join(val, ",")); err != nil {
		return err
	}
	return resourceServerRead(d, m)
}

func resourceServerDelete(d *schema.ResourceData, m interface{}) error {
	d.SetId("")
	return nil
}
