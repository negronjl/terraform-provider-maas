package main

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/juju/gomaasapi"
)

// NodeInfo detailed information from a node
type NodeInfo struct {
	tagNames     []string
	systemID     string
	hostname     string
	url          string
	powerState   string
	architecture string
	distroSeries string
	// hwe_kernel    string
	osystem  string
	memory   uint64
	data     map[string]interface{}
	cpuCount uint16
	status   uint16
}

// toNodeInfo Convenience function to convert a MAASObject to NodeInfo
func toNodeInfo(nodeObject *gomaasapi.MAASObject) (*NodeInfo, error) { // nolint: funlen
	log.Println("[DEBUG] [toNodeInfo] Attempting to convert node information from MAASObject to NodeInfo")

	nodeMap := nodeObject.GetMap()

	systemID, err := nodeMap["system_id"].GetString()
	if err != nil {
		log.Printf("[ERROR] [toNodeInfo] Unable to get node (%s)\n", systemID)
		return nil, err
	}

	hostname, err := nodeMap["hostname"].GetString()
	if err != nil {
		log.Printf("[ERROR] [toNodeInfo] Unable to get the node (%s) hostname\n", systemID)
		return nil, err
	}

	nodeURL := nodeObject.URL().String()
	if nodeURL == "" {
		return nil, errors.New("[ERROR] [toNodeInfo] Empty URL for node")
	}

	powerState, err := nodeMap["power_state"].GetString()
	if err != nil {
		log.Printf("[ERROR] [toNodeInfo] Unable to get the power state for node: %s\n", systemID)
		return nil, err
	}

	cpuCountFloat, err := nodeMap["cpu_count"].GetFloat64()
	if err != nil {
		log.Printf("[ERROR] [toNodeInfo] Unable to get the cpu_count for node: %s\n", systemID)
		log.Printf("[ERROR] [toNodeInfo] Error: %s\n", err)
		log.Printf("[ERROR] [toNodeInfo] cpuCountFloat: %v\n", cpuCountFloat)
		log.Println("[ERROR] [toNodeInfo] Defaulting cpuCount to 0")
		cpuCountFloat = 0
	}
	cpuCount := uint16(cpuCountFloat)

	architecture, err := nodeMap["architecture"].GetString()
	if err != nil {
		log.Printf("[ERROR] [toNodeInfo] Unable to get the node (%s) architecture\n", systemID)
		return nil, err
	}

	distroSeries, err := nodeMap["distro_series"].GetString()
	if err != nil {
		log.Printf("[ERROR] [toNodeInfo] Unable to get the distro_series for node: %s\n", systemID)
		return nil, err
	}

	memoryFloat, err := nodeMap["memory"].GetFloat64()
	if err != nil {
		log.Printf("[ERROR] [toNodeInfo] Unable to get the memory for node: %s\n", systemID)
		log.Printf("[ERROR] [toNodeInfo] Error: %s\n", err)
		log.Printf("[ERROR] [toNodeInfo] memory_float: %v\n", memoryFloat)
		log.Print("[ERROR] [toNodeInfo] Defaulting memory to 0")
		memoryFloat = 0
	}
	memory := uint64(memoryFloat)

	osystem, err := nodeMap["osystem"].GetString()
	if err != nil {
		log.Printf("[ERROR] [toNodeInfo] Unable to get the OS for node: %s\n", systemID)
		return nil, err
	}

	statusFloat, err := nodeMap["status"].GetFloat64()
	if err != nil {
		log.Printf("[ERROR] [toNodeInfo] Unable to get the status for node: %s\n", systemID)
		return nil, err
	}
	status := uint16(statusFloat)

	tagNames, err := nodeMap["tag_names"].GetArray()
	if err != nil {
		log.Printf("[ERROR] [toNodeInfo] Unable to get the tags for node: %s\n", systemID)
		return nil, err
	}

	tagArray := make([]string, 0, 1)
	var tagName string
	for _, tagObject := range tagNames {
		tagName, err = tagObject.GetString()
		if err != nil {
			log.Printf("[ERROR] [toNodeInfo] Unable to parse tag information (%v) for node (%s)", tagObject, systemID)
			return nil, err
		}
		tagArray = append(tagArray, tagName)
	}

	prettyJSON, err := json.MarshalIndent(nodeObject, "", "    ")
	if err != nil {
		log.Printf("[ERROR] [toNodeInfo] Unable to convert node (%s) information to JSON\n", systemID)
		return nil, err
	}

	log.Printf("[DEBUG] [toNodeInfo] Node (%s) JSON:\n%s\n", systemID, prettyJSON)

	var rawData map[string]interface{}
	if err := json.Unmarshal(prettyJSON, &rawData); err != nil {
		log.Printf("[ERROR] [toNodeInfo] Unable to Unmarshal JSON data for node: %s\n", systemID)
		return nil, err
	}

	return &NodeInfo{
		systemID:     systemID,
		hostname:     hostname,
		url:          nodeURL,
		powerState:   powerState,
		cpuCount:     cpuCount,
		architecture: architecture,
		distroSeries: distroSeries,
		memory:       memory,
		osystem:      osystem,
		status:       status,
		tagNames:     tagArray,
		data:         rawData}, nil
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
