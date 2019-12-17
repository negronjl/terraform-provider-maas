package main

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/juju/gomaasapi"
)

// NodeInfo detailed information from a node
type NodeInfo struct {
	system_id     string
	hostname      string
	url           string
	power_state   string
	cpu_count     uint16
	architecture  string
	distro_series string
	hwe_kernel    string
	memory        uint64
	osystem       string
	status        uint16
	tag_names     []string
	data          map[string]interface{}
}

// toNodeInfo Convenience function to convert a MAASObject to NodeInfo
func toNodeInfo(nodeObject *gomaasapi.MAASObject) (*NodeInfo, error) {
	log.Println("[DEBUG] [toNodeInfo] Attempting to convert node information from MAASObject to NodeInfo")

	nodeMap := nodeObject.GetMap()

	system_id, err := nodeMap["system_id"].GetString()
	if err != nil {
		log.Printf("[ERROR] [toNodeInfo] Unable to get node (%s)\n", system_id)
		return nil, err
	}

	hostname, err := nodeMap["hostname"].GetString()
	if err != nil {
		log.Printf("[ERROR] [toNodeInfo] Unable to get the node (%s) hostname\n", system_id)
		return nil, err
	}

	node_url := nodeObject.URL().String()
	if len(node_url) == 0 {
		return nil, errors.New("[ERROR] [toNodeInfo] Empty URL for node")
	}

	power_state, err := nodeMap["power_state"].GetString()
	if err != nil {
		log.Printf("[ERROR] [toNodeInfo] Unable to get the power_state for node: %s\n", system_id)
		return nil, err
	}

	cpu_count_float, err := nodeMap["cpu_count"].GetFloat64()
	if err != nil {
		log.Printf("[ERROR] [toNodeInfo] Unable to get the cpu_count for node: %s\n", system_id)
		log.Printf("[ERROR] [toNodeInfo] Error: %s\n", err)
		log.Printf("[ERROR] [toNodeInfo] cpu_count_float: %v\n", cpu_count_float)
		log.Println("[ERROR] [toNodeInfo] Defaulting cpu_count to 0")
		cpu_count_float = 0
	}
	cpu_count := uint16(cpu_count_float)

	architecture, err := nodeMap["architecture"].GetString()
	if err != nil {
		log.Printf("[ERROR] [toNodeInfo] Unable to get the node (%s) architecture\n", system_id)
		return nil, err
	}

	distro_series, err := nodeMap["distro_series"].GetString()
	if err != nil {
		log.Printf("[ERROR] [toNodeInfo] Unable to get the distro_series for node: %s\n", system_id)
		return nil, err
	}

	memory_float, err := nodeMap["memory"].GetFloat64()
	if err != nil {
		log.Printf("[ERROR] [toNodeInfo] Unable to get the memory for node: %s\n", system_id)
		log.Printf("[ERROR] [toNodeInfo] Error: %s\n", err)
		log.Printf("[ERROR] [toNodeInfo] memory_float: %v\n", memory_float)
		log.Print("[ERROR] [toNodeInfo] Defaulting memory to 0")
		memory_float = 0
	}
	memory := uint64(memory_float)

	osystem, err := nodeMap["osystem"].GetString()
	if err != nil {
		log.Printf("[ERROR] [toNodeInfo] Unable to get the OS for node: %s\n", system_id)
		return nil, err
	}

	status_float, err := nodeMap["status"].GetFloat64()
	if err != nil {
		log.Printf("[ERROR] [toNodeInfo] Unable to get the status for node: %s\n", system_id)
		return nil, err
	}
	status := uint16(status_float)

	tag_names, err := nodeMap["tag_names"].GetArray()
	if err != nil {
		log.Printf("[ERROR] [toNodeInfo] Unable to get the tags for node: %s\n", system_id)
		return nil, err
	}

	tag_array := make([]string, 0, 1)
	for _, tag_object := range tag_names {
		tag_name, err := tag_object.GetString()
		if err != nil {
			log.Printf("[ERROR] [toNodeInfo] Unable to parse tag information (%v) for node (%s)", tag_object, system_id)
			return nil, err
		}
		tag_array = append(tag_array, tag_name)
	}

	prettyJSON, err := json.MarshalIndent(nodeObject, "", "    ")
	if err != nil {
		log.Printf("[ERROR] [toNodeInfo] Unable to convert node (%s) information to JSON\n", system_id)
		return nil, err
	}

	log.Printf("[DEBUG] [toNodeInfo] Node (%s) JSON:\n%s\n", system_id, prettyJSON)

	var raw_data map[string]interface{}
	if err := json.Unmarshal(prettyJSON, &raw_data); err != nil {
		log.Printf("[ERROR] [toNodeInfo] Unable to Unmarshal JSON data for node: %s\n", system_id)
		return nil, err
	}

	return &NodeInfo{system_id: system_id,
		hostname:      hostname,
		url:           node_url,
		power_state:   power_state,
		cpu_count:     uint16(cpu_count),
		architecture:  architecture,
		distro_series: distro_series,
		memory:        memory,
		osystem:       osystem,
		status:        uint16(status),
		tag_names:     tag_array,
		data:          raw_data}, nil
}

// Config provider configuration
type Config struct {
	APIKey     string
	APIURL     string
	APIver     string
	MAASObject *gomaasapi.MAASObject
}

// Client authenticate to MAAS and create a session
func (c *Config) Client() (interface{}, error) {
	log.Println("[DEBUG] [Config.Client] Configuring the MAAS API client")
	authClient, err := gomaasapi.NewAuthenticatedClient(
		gomaasapi.AddAPIVersionToURL(c.APIURL, c.APIver), c.APIKey)
	if err != nil {
		log.Printf("[ERROR] [Config.Client] Unable to authenticate against the MAAS Server (%s)", c.APIURL)
		return nil, err
	}
	c.MAASObject = gomaasapi.NewMAAS(*authClient)
	return c, nil
}
