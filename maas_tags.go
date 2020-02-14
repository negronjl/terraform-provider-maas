package main

import (
	"log"
	"net/url"

	"github.com/juju/gomaasapi"
)

// tagCreate creates a new tag
func tagCreate(maas *gomaasapi.MAASObject, tagName string) error {
	log.Printf("[DEBUG] [tagCreate] Creating new tag named %s", tagName)

	params := url.Values{}
	params.Set("name", tagName)
	_, err := maas.GetSubObject("tags").CallPost("", params)
	return err
}

// noteTagsUpdate update the tags for a node
func nodeTagsUpdate(maas *gomaasapi.MAASObject, systemID, tagName string) error {
	log.Println("[DEBUG] [nodeUpdate] Attempting to update a node's tags")

	// make sure tag exists
	tagObject, err := maas.GetSubObject("tags").GetSubObject(tagName).Get()
	if err != nil {
		// create tag if it doesn't exist
		log.Printf("[ERROR] [nodeTagsUpdate] Tag %s does not exist", tagName)
		err = tagCreate(maas, tagName)
		if err != nil {
			return err
		}
		tagObject, err = maas.GetSubObject("tags").GetSubObject(tagName).Get()
		if err != nil {
			return nil
		}
	}

	params := url.Values{}
	params.Set("add", systemID)
	_, err = tagObject.CallPost("update_nodes", params)
	if err != nil {
		log.Printf("[ERROR] [nodeTagsUpdate] Unable to update node (%s) tag (%s).  Failed with error (%s)\n",
			systemID, tagName, err)
		return err
	}
	return nil
}

// noteTagsRemove remove the deploy tags for a node
func nodeTagsRemove(maas *gomaasapi.MAASObject, systemID, tagName string) error {
	log.Println("[DEBUG] [nodeUpdate] Attempting to remove a node's tag")

	// make sure tag exists
	tagObject, err := maas.GetSubObject("tags").GetSubObject(tagName).Get()
	if err != nil {
		return err
	}

	params := url.Values{}
	params.Set("remove", systemID)
	_, err = tagObject.CallPost("update_nodes", params)
	if err != nil {
		log.Printf("[ERROR] [nodeTagsUpdate] Unable to update node (%s) tag (%s).  Failed with error (%s)\n",
			systemID, tagName, err)
		return err
	}
	return nil
}
