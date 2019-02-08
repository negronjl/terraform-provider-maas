package main

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"log"
	"net/url"
	"regexp"
	"strconv"
	"time"

	"github.com/juju/gomaasapi"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

// resourceMAASMachine creates a new terraform schema resource
func resourceMAASMachine() *schema.Resource {
	log.Println("[DEBUG] [resourceMAASMachine] Initializing data structure")
	return &schema.Resource{
		Create: resourceMAASMachineCreate,
		Read:   resourceMAASMachineRead,
		Update: resourceMAASMachineUpdate,
		Delete: resourceMAASMachineDelete,

		SchemaVersion: 1,

		Schema: map[string]*schema.Schema{
			"architecture": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"boot_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"cpu_count": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},

			"disable_ipv4": {
				Type:     schema.TypeBool,
				Optional: true,
			},

			"distro_series": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"hostname": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"deploy_hostname": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"deploy_tags": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"deploy_interface": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"type": {
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: validation.StringMatch(regexp.MustCompile(`^\s*vlan\s*$|^\s*physical\s*$|^\s*bond\s*$`), "Must be vlan, physical, dhcp or bond."),
						},
						"fabric": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"vlan": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
						},
						"subnet": {
							Type:         schema.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: validation.CIDRNetwork(1, 32),
						},
						"ip_mode": {
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: validation.StringMatch(regexp.MustCompile(`^\s*auto\s*$|^\s*static\s*$|^\s*dhcp\s*$|^\s*unconfigured\s*$`), "Must be auto, static, dhcp or unconfigured."),
						},
						"ip": {
							Type:         schema.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: validation.SingleIP(),
						},
					},
				},
			},

			"tags": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"not_tags": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"release_erase": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: false,
				Default:  true,
			},

			"release_erase_secure": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: false,
				Default:  false,
			},

			"release_erase_quick": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: false,
				Default:  false,
			},

			"ip_addresses": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"macaddress_set": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"mac_address": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"resource_uri": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},

			"memory": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},

			"netboot": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},

			"osystem": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"owner": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"physicalblockdevice_set": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"block_size": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"id": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"id_path": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"model": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"path": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"serial": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"size": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"tags": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},

			"power_state": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"power_type": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"pxe_mac": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"mac_address": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"resource_uri": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},

			"resource_uri": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"routers": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"status": {
				Type:     schema.TypeInt,
				Optional: true,
			},

			"volumes": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"label": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"size": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"tags": {
							Type:     schema.TypeList,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},

			"swap_size": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"system_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"tag_names": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"zone": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"description": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"resource_uri": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},

			"not_zones": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"user_data": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				StateFunc: func(v interface{}) string {
					switch v.(type) {
					case string:
						hash := sha1.Sum([]byte(v.(string)))
						return hex.EncodeToString(hash[:])
					default:
						return ""
					}
				},
			},

			"hwe_kernel": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"comment": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

// resourceMAASMachineCreate This function doesn't really *create* a new node but, power an already registered
func resourceMAASMachineCreate(d *schema.ResourceData, meta interface{}) error {
	log.Println("[DEBUG] [resourceMAASMachineCreate] Launching new maas_machine")

	/*
		According to the MAAS API documentation here: https://maas.ubuntu.com/docs/api.html
		We need to acquire or allocate a node before we start it.  We pass (as url.Values)
		some parameters that could be used to narrow down our selection (cpu_count, memory, etc.)
	*/

	controller := meta.(*Config).controller
	acquireParams := convertConstraints(d)

	machine, _, err := controller.AllocateMachine(acquireParams)
	if err != nil {
		log.Println("[ERROR] [resourceMAASMachineCreate] Unable to allocate machine.")
		return err
	}

	// set the node id
	d.SetId(machine.SystemID())

	// update machine attributes during allocated state
	params := url.Values{}
	if hostname, ok := d.GetOk("deploy_hostname"); ok {
		params.Add("hostname", hostname.(string))
	}

	err = nodeUpdate(meta.(*Config).MAASObject, d.Id(), params)
	if err != nil {
		log.Println("[DEBUG] Unable to update node")
	}

	// Update networking interfaces, if set
	if v, ok := d.GetOk("deploy_interface"); ok {
		subnets, err := getSubnets(controller)
		if err != nil {
			log.Println("[ERROR] [resourceMAASMachineCreate] Unable to get subnets")
			if err := resourceMAASMachineDelete(d, meta); err != nil {
				log.Printf("[DEBUG] Unable to release node: %s", err.Error())
			}
			return err
		}

		for _, interfaces_ := range v.(*schema.Set).List() {
			i := interfaces_.(map[string]interface{})

			if cidr, ok := i["subnet"].(string); ok {
				if subnet, ok := subnets[cidr]; ok {
					nics := machine.InterfaceSet()
					for _, nic := range nics {
						if nic.Name() == i["name"].(string) {
							err := nic.LinkSubnet(gomaasapi.LinkSubnetArgs{
								Mode:      gomaasapi.LinkModeStatic,
								Subnet:    subnet,
								IPAddress: i["ip"].(string),
							})
							if err != nil {
								log.Println("[ERROR] [resourceMAASMachineCreate] Unable to link subnet")
								if err := resourceMAASMachineDelete(d, meta); err != nil {
									log.Printf("[DEBUG] Unable to release node: %s", err.Error())
								}
								return err
							}
						}
					}
				}
			}
		}
	}

	startArgs := gomaasapi.StartArgs{
		UserData:     base64encode(d.Get("user_data").(string)),
		DistroSeries: d.Get("distro_series").(string),
		Kernel:       d.Get("hwe_kernel").(string),
		Comment:      d.Get("comment").(string),
	}

	if err := machine.Start(startArgs); err != nil {
		log.Printf("[ERROR] [resourceMAASMachineCreate] Unable to power up node: %s\n", d.Id())
		// unable to perform action, release the node
		if err := resourceMAASMachineDelete(d, meta); err != nil {
			log.Printf("[DEBUG] Unable to release node: %s", err.Error())
		}
		return err
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"Deploying", "Ready"},
		Target:     []string{"Deployed"},
		Refresh:    getNodeStatus(meta.(*Config).MAASObject, d.Id()),
		Timeout:    25 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		if err := resourceMAASMachineDelete(d, meta); err != nil {
			log.Printf("[DEBUG] Unable to release node: %s", err.Error())
		}
		return fmt.Errorf("[ERROR] [resourceMAASMachineCreate] Error waiting for machine (%s) to become deployed: %s", machine.SystemID(), err)
	}

	// update node tags
	if tags, ok := d.GetOk("deploy_tags"); ok {
		for i := range tags.([]interface{}) {
			err := nodeTagsUpdate(meta.(*Config).MAASObject, d.Id(), tags.([]interface{})[i].(string))
			if err != nil {
				log.Printf("[ERROR] Unable to update node (%s) with tag (%s)", d.Id(), tags.([]interface{})[i].(string))
			}
		}
	}

	return resourceMAASMachineUpdate(d, meta)

}

// resourceMAASMachineRead read machine information from a maas node
// TODO: remove or do something
func resourceMAASMachineRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Reading machine (%s) information.\n", d.Id())
	return nil
}

// resourceMAASMachineUpdate update machine in terraform state
func resourceMAASMachineUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] [resourceMAASMachineUpdate] Modifying machine %s\n", d.Id())

	d.Partial(true)

	d.Partial(false)

	log.Printf("[DEBUG] Done Modifying machine %s", d.Id())
	return resourceMAASMachineRead(d, meta)
}

// resourceMAASMachineDelete This function doesn't really *delete* a maas managed machine but releases (read, turns off) the node.
//     TODO: this should implement gomaasapi.Controller.ReleaseMachines, but currently ReleaseMachines doesn't support erase release
func resourceMAASMachineDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Deleting machine %s\n", d.Id())

	controller := meta.(*Config).controller
	var ids []string
	machines, err := controller.Machines(gomaasapi.MachinesArgs{
		SystemIDs: append(ids, d.Id()),
	})
	if err != nil {
		log.Printf("[ERROR] [resourceMAASMachineDetele] cannnot list machines")
		return err
	}
	if len(machines) != 1 {
		return fmt.Errorf("[ERROR] [resourceMAASMachineDetele] machine no longer exists")
	}
	release_params := url.Values{}

	if release_erase, ok := d.GetOk("release_erase"); ok {
		release_params.Add("erase", strconv.FormatBool(release_erase.(bool)))
	}

	if release_erase_secure, ok := d.GetOk("release_erase_secure"); ok {
		// setting erase to true in the event a user didn't set both options
		release_params.Add("erase", strconv.FormatBool(true))
		release_params.Add("secure_erase", strconv.FormatBool(release_erase_secure.(bool)))
	}

	if release_erase_quick, ok := d.GetOk("release_erase_quick"); ok {
		// setting erase to true in the event a user didn't set both options
		release_params.Add("erase", strconv.FormatBool(true))
		release_params.Add("quick_erase", strconv.FormatBool(release_erase_quick.(bool)))
	}

	if err := nodeRelease(meta.(*Config).MAASObject, d.Id(), release_params); err != nil {
		return err
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"Deployed", "Releasing", "Disk erasing"},
		Target:     []string{"Ready"},
		Refresh:    getNodeStatus(meta.(*Config).MAASObject, d.Id()),
		Timeout:    30 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf(
			"[ERROR] [resourceMAASMachineDelete] Error waiting for machine (%s) to become ready: %s", d.Id(), err)
	}

	//unlink any subnets
	if v, ok := d.GetOk("deploy_interface"); ok {
		subnets, err := getSubnets(controller)
		if err != nil {
			log.Println("[WARN] [resourceMAASMachineDelete] Unable to get subnets.")
		} else {
			for _, netInterfaces := range v.(*schema.Set).List() {
				i := netInterfaces.(map[string]interface{})

				if cidr, ok := i["subnet"].(string); ok {
					if subnet, ok := subnets[cidr]; ok {
						nics := machines[0].InterfaceSet()
						for _, nic := range nics {
							if nic.Name() == i["name"].(string) {
								err := nic.UnlinkSubnet(subnet)
								if err != nil {
									log.Println("[ERROR] [resourceMAASMachineDelete] Unable to unlink subnet.")
									return err
								}
							}
						}
					}
				}
			}
		}
	}

	// remove deploy hostname if set
	if _, ok := d.GetOk("deploy_hostname"); ok {
		params := url.Values{}
		params.Set("hostname", "")
		err := nodeUpdate(meta.(*Config).MAASObject, d.Id(), params)
		if err != nil {
			log.Println("[DEBUG] Unable to reset hostname: %s", err)
		}
	}

	// remove deployed tags
	if tags, ok := d.GetOk("deploy_tags"); ok {
		for i := range tags.([]interface{}) {
			err := nodeTagsRemove(meta.(*Config).MAASObject, d.Id(), tags.([]interface{})[i].(string))
			if err != nil {
				log.Printf("[ERROR] Unable to update node (%s) with tag (%s)", d.Id(), tags.([]interface{})[i].(string))
			}
		}
	}

	log.Printf("[DEBUG] [resourceMAASMachineDelete] Node (%s) released", d.Id())

	d.SetId("")

	return nil
}

func convertConstraints(d *schema.ResourceData) gomaasapi.AllocateMachineArgs {
	args := gomaasapi.AllocateMachineArgs{}

	if hostname, ok := d.GetOk("hostname"); ok {
		args.Hostname = hostname.(string)
	}

	if systemID, ok := d.GetOk("system_id"); ok {
		args.SystemId = systemID.(string)
	}

	if architecture, ok := d.GetOk("architecture"); ok {
		args.Architecture = architecture.(string)
	}

	if minCPUCount, ok := d.GetOk("cpu_count"); ok {
		args.MinCPUCount = minCPUCount.(int)
	}

	if minRAM, ok := d.GetOk("memory"); ok {
		args.MinMemory = minRAM.(int)
	}

	if tags, ok := d.GetOk("tags"); ok {
		args.Tags = expandStringList(tags.([]interface{}))
	}

	if notTags, ok := d.GetOk("not_tags"); ok {
		args.NotTags = expandStringList(notTags.([]interface{}))
	}

	zone := d.Get("zone").(*schema.Set).List()

	if len(zone) > 0 {
		//use first zone only
		z := zone[0].(map[string]interface{})
		args.Zone = z["name"].(string)
	}

	if notZones, ok := d.GetOk("not_zones"); ok {
		args.NotInZone = expandStringList(notZones.([]interface{}))
	}

	volumes := d.Get("volumes").(*schema.Set).List()

	for _, vol := range volumes {
		v := vol.(map[string]interface{})
		args.Storage = append(args.Storage, gomaasapi.StorageSpec{
			Label: v["label"].(string),
			Size:  v["size"].(int),
			Tags:  v["tags"].([]string),
		})
	}

	if comment, ok := d.GetOk("comment"); ok {
		args.Comment = comment.(string)
	}

	return args
}

func getSubnets(controller gomaasapi.Controller) (map[string]gomaasapi.Subnet, error) {
	// Get all the spaces
	spaces, err := controller.Spaces()
	if err != nil {
		return nil, err
	}

	// Get all the subnets
	subnets := map[string]gomaasapi.Subnet{}
	for _, space := range spaces {
		for _, subnet := range space.Subnets() {
			subnets[subnet.Name()] = subnet
		}
	}

	return subnets, nil
}
