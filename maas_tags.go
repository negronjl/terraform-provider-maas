package main

import (
	"log"
	"net/url"

	"github.com/juju/gomaasapi"
)

// tagCreate creates a new tag
func tagCreate(maas *gomaasapi.MAASObject, tag_name string) error {
	log.Printf("[DEBUG] [tagCreate] Creating new tag named %s", tag_name)

	params := url.Values{}
	params.Set("name", tag_name)
	_, err := maas.GetSubObject("tags").CallPost("", params)
	return err
}

// noteTagsUpdate update the tags for a node
func nodeTagsUpdate(maas *gomaasapi.MAASObject, system_id string, tag_name string) error {
	log.Println("[DEBUG] [nodeUpdate] Attempting to update a node's tags")

	// make sure tag exists
	tagObject, err := maas.GetSubObject("tags").GetSubObject(tag_name).Get()
	if err != nil {
		// create tag if it doesn't exist
		log.Println("[ERROR] [nodeTagsUpdate] Tag %s does not exist", tag_name)
		err := tagCreate(maas, tag_name)
		if err != nil {
			return err
		}
		tagObject, err = maas.GetSubObject("tags").GetSubObject(tag_name).Get()
		if err != nil {
			return nil
		}
	}

	params := url.Values{}
	params.Set("add", system_id)
	_, err = tagObject.CallPost("update_nodes", params)
	if err != nil {
		log.Printf("[ERROR] [nodeTagsUpdate] Unable to update node (%s) tag (%s).  Failed withh error (%s)\n", system_id, tag_name, err)
		return err
	}
	return nil
}

// noteTagsRemove remove the deploy tags for a node
func nodeTagsRemove(maas *gomaasapi.MAASObject, system_id string, tag_name string) error {
	log.Println("[DEBUG] [nodeUpdate] Attempting to remove a node's tag")

	// make sure tag exists
	tagObject, err := maas.GetSubObject("tags").GetSubObject(tag_name).Get()
	if err != nil {
		return err
	}

	params := url.Values{}
	params.Set("remove", system_id)
	_, err = tagObject.CallPost("update_nodes", params)
	if err != nil {
		log.Printf("[ERROR] [nodeTagsUpdate] Unable to update node (%s) tag (%s).  Failed withh error (%s)\n", system_id, tag_name, err)
		return err
	}
	return nil
}
