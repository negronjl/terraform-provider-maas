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
func resourceMAASInstanceCreate(d *schema.ResourceData, meta interface{}) error { // nolint: funlen, gocognit
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
	d.SetId(nodeObj.systemID)

	// separate constraints that are supported for the deploy action
	// parameters to pass when creating a node
	nodeParams := url.Values{}

	// get user data if defined
	if userData, ok := d.GetOk("user_data"); ok {
		nodeParams.Add("user_data", base64encode(userData.(string)))
	}

	// get comment if defined
	if comment, ok := d.GetOk("comment"); ok {
		nodeParams.Add("comment", comment.(string))
	}

	// get distro_series if defined
	distroSeries, ok := d.GetOk("distro_series")
	if ok {
		nodeParams.Add("distro_series", distroSeries.(string))
	}

	// get hwe_kernel if defined
	if HWEKernel, ok := d.GetOk("hwe_kernel"); ok {
		nodeParams.Add("hwe_kernel", HWEKernel.(string))
	}

	// install kvm and register the server as a kvm server if requested
	if installKVM, ok := d.GetOk("install_kvm"); ok {
		log.Printf("[INFO] Adding KVM packages and configuration: %s", installKVM)
		nodeParams.Add("install_kvm", strconv.FormatBool(true))
	}

	// install rackd if requested
	if installRackD, ok := d.GetOk("install_rackd"); ok {
		log.Printf("[INFO] Adding maas rack controller packages: %s", installRackD)
		nodeParams.Add("install_rackd", strconv.FormatBool(true))
	}

	if err = nodeDo(meta.(*Config).MAASObject, d.Id(), "deploy", nodeParams); err != nil {
		log.Printf("[ERROR] [resourceMAASInstanceCreate] Unable to power up node: %s\n", d.Id())
		// unable to perform action, release the node
		if err = nodeRelease(meta.(*Config).MAASObject, d.Id(), url.Values{}); err != nil {
			log.Printf("[DEBUG] Unable to release node")
		}
		return err
	}

	log.Printf("[DEBUG] [resourceMAASInstanceCreate] Waiting for instance (%s) to become active\n", d.Id())
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"9:"},
		Target:     []string{"6:"},
		Refresh:    getNodeStatus(meta.(*Config).MAASObject, d.Id()),
		Timeout:    25 * time.Minute, // nolint: gomnd
		Delay:      10 * time.Second, // nolint: gomnd
		MinTimeout: 3 * time.Second,  // nolint: gomnd
	}

	if _, err = stateConf.WaitForState(); err != nil {
		if err = nodeRelease(meta.(*Config).MAASObject, d.Id(), url.Values{}); err != nil {
			log.Printf("[DEBUG] Unable to release node")
		}
		return fmt.Errorf("[ERROR] [resourceMAASInstanceCreate] Error waiting for instance (%s) to become deployed: %s",
			d.Id(), err)
	}

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

// resourceMAASInstanceDelete
// This function doesn't really *delete* a maas managed instance but releases (read, turns off) the node.
func resourceMAASInstanceDelete(d *schema.ResourceData, meta interface{}) error { // nolint: funlen
	log.Printf("[DEBUG] Deleting instance %s\n", d.Id())
	releaseParams := url.Values{}

	if releaseErase, ok := d.GetOk("release_erase"); ok {
		releaseParams.Add("erase", strconv.FormatBool(releaseErase.(bool)))
	}

	if releaseEraseSecure, ok := d.GetOk("release_erase_secure"); ok {
		// setting erase to true in the event a user didn't set both options
		releaseParams.Add("erase", strconv.FormatBool(true))
		releaseParams.Add("secure_erase", strconv.FormatBool(releaseEraseSecure.(bool)))
	}

	if releaseEraseQuick, ok := d.GetOk("release_erase_quick"); ok {
		// setting erase to true in the event a user didn't set both options
		releaseParams.Add("erase", strconv.FormatBool(true))
		releaseParams.Add("quick_erase", strconv.FormatBool(releaseEraseQuick.(bool)))
	}

	if err := nodeRelease(meta.(*Config).MAASObject, d.Id(), releaseParams); err != nil {
		return err
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"6:", "12:", "14:"},
		Target:     []string{"4:"},
		Refresh:    getNodeStatus(meta.(*Config).MAASObject, d.Id()),
		Timeout:    30 * time.Minute, // nolint: gomnd
		Delay:      10 * time.Second, // nolint: gomnd
		MinTimeout: 3 * time.Second,  // nolint: gomnd
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf(
			"[ERROR] [resourceMAASInstanceCreate] Error waiting for instance (%s) to become ready: %s", d.Id(), err)
	}

	// remove deploy hostname if set
	if _, ok := d.GetOk("deploy_hostname"); ok {
		params := url.Values{}
		params.Set("hostname", "")
		err := nodeUpdate(meta.(*Config).MAASObject, d.Id(), params)
		if err != nil {
			log.Printf("[DEBUG] Unable to reset hostname: %s", err)
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

	log.Printf("[DEBUG] [resourceMAASInstanceDelete] Node (%s) released", d.Id())

	d.SetId("")

	return nil
}
