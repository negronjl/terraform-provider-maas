package main

import (
	"fmt"
	"log"
	"net/url"
	"strconv"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

// resourceMAASInstanceCreate This function doesn't really *create* a new node but, power an already registered
func resourceMAASInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	log.Println("[DEBUG] [resourceMAASInstanceCreate] Launching new maas_instance")

	/*
		According to the MAAS API documentation here: https://maas.ubuntu.com/docs/api.html
		We need to acquire or allocate a node before we start it.  We pass (as url.Values)
		some parameters that could be used to narrow down our selection (cpu_count, memory, etc.)
	*/

	constraints, err := parseConstraints(d)
	if err != nil {
		log.Println("[ERROR] [resourceMAASInstanceCreate] Unable to parse constraints.")
		return err
	}

	nodeObj, err := nodesAllocate(meta.(*Config).MAASObject, constraints)
	if err != nil {
		log.Println("[ERROR] [resourceMAASInstanceCreate] Unable to allocate nodes")
		return err
	}

	// set the node id
	d.SetId(nodeObj.system_id)

	// seperate constraints that are supported for the deploy action
	// parameters to pass when creating a node
	node_params := url.Values{}

	// get user data if defined
	if user_data, ok := d.GetOk("user_data"); ok {
		node_params.Add("user_data", base64encode(user_data.(string)))
	}

	// get comment if defined
	if comment, ok := d.GetOk("comment"); ok {
		node_params.Add("comment", comment.(string))
	}

	// get distro_series if defined
	distro_series, ok := d.GetOk("distro_series")
	if ok {
		node_params.Add("distro_series", distro_series.(string))
	}

	if err := nodeDo(meta.(*Config).MAASObject, d.Id(), "deploy", node_params); err != nil {
		log.Printf("[ERROR] [resourceMAASInstanceCreate] Unable to power up node: %s\n", d.Id())
		// unable to perform action, release the node
		if err := nodeRelease(meta.(*Config).MAASObject, d.Id(), url.Values{}); err != nil {
			log.Printf("[DEBUG] Unable to release node")
		}
		return err
	}

	log.Printf("[DEBUG] [resourceMAASInstanceCreate] Waiting for instance (%s) to become active\n", d.Id())
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"9:"},
		Target:     []string{"6:"},
		Refresh:    getNodeStatus(meta.(*Config).MAASObject, d.Id()),
		Timeout:    25 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	nodeObjRaw, err := stateConf.WaitForState()
	if err != nil {
		if err := nodeRelease(meta.(*Config).MAASObject, d.Id(), url.Values{}); err != nil {
			log.Printf("[DEBUG] Unable to release node")
		}
		return fmt.Errorf("[ERROR] [resourceMAASInstanceCreate] Error waiting for instance (%s) to become deployed: %s", d.Id(), err)
	}

	nodeObj = nodeObjRaw.(*NodeInfo)

	// update node
	params := url.Values{}
	if hostname, ok := d.GetOk("deploy_hostname"); ok {
		params.Add("hostname", hostname.(string))
	}

	// only updating hostname
	err = nodeUpdate(meta.(*Config).MAASObject, d.Id(), params)
	if err != nil {
		log.Println("[DEBUG] Unable to update node")
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

	d.Set("hostname", nodeObj.fqdn)
	d.Set("ip_addresses", nodeObj.ip_addresses)
	d.SetConnInfo(map[string]string{
		"type": "ssh",
		"host": nodeObj.fqdn,
	})

	return resourceMAASInstanceUpdate(d, meta)
}

// resourceMAASInstanceRead read instance information from a maas node
// TODO: remove or do something
func resourceMAASInstanceRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Reading instance (%s) information.\n", d.Id())
	return nil
}

// resourceMAASInstanceUpdate update an instance in terraform state
func resourceMAASInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] [resourceMAASInstanceUpdate] Modifying instance %s\n", d.Id())

	d.Partial(true)

	d.Partial(false)

	log.Printf("[DEBUG] Done Modifying instance %s", d.Id())
	return resourceMAASInstanceRead(d, meta)
}

// resourceMAASInstanceDelete This function doesn't really *delete* a maas managed instance but releases (read, turns off) the node.
func resourceMAASInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Deleting instance %s\n", d.Id())
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
		Pending:    []string{"6:", "12:", "14:"},
		Target:     []string{"4:"},
		Refresh:    getNodeStatus(meta.(*Config).MAASObject, d.Id()),
		Timeout:    30 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf(
			"[ERROR] [resourceMAASInstanceCreate] Error waiting for instance (%s) to become ready: %s", d.Id(), err)
	}

	params := url.Values{}
	params.Set("hostname", "")
	// remove custom hostname
	err := nodeUpdate(meta.(*Config).MAASObject, d.Id(), params)
	if err != nil {
		log.Println("[DEBUG] Unable to update node")
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

	log.Printf("[DEBUG] [resourceMAASInstanceDelete] Node (%s) released", d.Id())

	d.SetId("")

	return nil
}
