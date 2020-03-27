package provider

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/juju/gomaasapi"
	"github.com/roblox/terraform-provider-maas/pkg/gmaw"
	"github.com/roblox/terraform-provider-maas/pkg/maas"
)

// resourceInstance provides a resource that correlates to the MaaS Machine endpoint
func resourceInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceInstanceCreate,
		Read:   resourceInstanceRead,
		Update: resourceInstanceUpdate,
		Delete: resourceInstanceDelete,

		Schema: map[string]*schema.Schema{
			"address": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceInstanceCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*gomaasapi.MAASObject)
	var instance Instance

	// Start by allocating the machine
	machinesManager := maas.NewMachinesManager(gmaw.NewMachines(client))
	ap := NewInstance(d).AllocateParams()
	ma, err := machinesManager.Allocate(ap)
	if err != nil {
		return err
	}

	// Go ahead and set the ID if it allocates. That will add the resource to the state.
	instance.FromMachine(ma).UpdateState(d)
	d.SetId(ma.SystemID)

	machineManager, err := maas.NewMachineManager(ma.SystemID, gmaw.NewMachine(client))
	if err != nil {
		return err
	}

	// Get node params and set them along with the deploy
	dp := instance.FromMachine(machineManager.Current()).DeployParams()
	if err := machineManager.Deploy(dp); err != nil {
		machinesManager.Release([]string{machineManager.SystemID()}, "The deploy has broke") // nolint
	}

	// Lock the machine, if necessary
	if instance.Lock {
		machineManager.Lock("Locked by Terraform") // nolint
	}
	return resourceInstanceRead(d, m)
}

func resourceInstanceRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*gomaasapi.MAASObject)
	var instance Instance

	machineManager, err := maas.NewMachineManager(d.Id(), gmaw.NewMachine(client))
	if err != nil {
		// FIXME if the error is that the system ID does not exist, d.SetId("")
		return err
	}

	instance.FromMachine(machineManager.Current()).UpdateState(d)
	return nil
}

func resourceInstanceUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*gomaasapi.MAASObject)
	machineManager, err := maas.NewMachineManager(d.Id(), gmaw.NewMachine(client))
	if err != nil {
		return err
	}

	// Lock the machine, if necessary
	if NewInstance(d).Lock {
		machineManager.Lock("Locked by Terraform") // nolint
	}
	return resourceInstanceRead(d, m)
}

func resourceInstanceDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}
