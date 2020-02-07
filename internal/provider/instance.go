package provider

import (
	"crypto/sha1" // nolint: gosec
	"encoding/hex"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/roblox/terraform-provider-maas/pkg/maas"
)

// MACAddress is used by the maas_instance resource
type MACAddress maas.MACAddress

// Instance models the maas_instance Terraform resource
type Instance struct {
	DeployTags             []string      `optional:"true" forcenew:"true"`
	Tags                   []string      `optional:"true" forcenew:"true"`
	IPAddresses            []string      `optional:"true" forcenew:"true"`
	MACAddressSet          []MACAddress  `optional:"true" forcenew:"true"`
	PhysicalBlockDeviceSet []BlockDevice `optional:"true" forcenew:"true" name:"physicalblockdevice_set"`
	PXEMac                 []MACAddress  `optional:"true" type:"Set"`
	Routers                []string      `optional:"true"`
	TagNames               []string      `optional:"true"`
	Zone                   []Zone        `optional:"true" type:"Set"`
	Architecture           string        `optional:"true" forcenew:"true"`
	BootType               string        `optional:"true" forcenew:"true"`
	DistroSeries           string        `optional:"true" forcenew:"true"`
	Hostname               string        `optional:"true" forcenew:"true"`
	DeployHostname         string        `optional:"true" forcenew:"true"`
	OSystem                string        `optional:"true" forcenew:"true"`
	Owner                  string        `optional:"true" forcenew:"true"`
	PowerState             string        `optional:"true"`
	PowerType              string        `optional:"true"`
	ResourceURI            string        `optional:"true" forcenew:"true"`
	SystemID               string        `optional:"true" forcenew:"true"`
	UserData               string        `optional:"true" forcenew:"true" statefunc:"true"`
	HWEKernel              string        `optional:"true" forcenew:"true"`
	Comment                string        `optional:"true"`
	CPUCount               int           `optional:"true" forcenew:"true"`
	Memory                 int           `optional:"true" forcenew:"true"`
	Status                 int           `optional:"true"`
	Storage                int           `optional:"true"`
	SwapSize               int           `optional:"true"`
	DisableIPv4            bool          `optional:"true" name:"disable_ipv4"`
	ReleaseErase           bool          `optional:"true" forcenew:"true" default:"false"`
	ReleaseEraseSecure     bool          `optional:"true" forcenew:"true" default:"false"`
	ReleaseEraseQuick      bool          `optional:"true" forcenew:"true" default:"false"`
	Netboot                bool          `optional:"true" forcenew:"true"`
	InstallKVM             bool          `optional:"true" forcenew:"true" default:"false"`
	Lock                   bool          `optional:"true" default:"false"`
}

// NewInstance creates a new instance from the value of the Terraform resource
func NewInstance(resource *schema.ResourceData) *Instance {
	var instance Instance
	st := reflect.TypeOf(instance)
	sv := reflect.ValueOf(instance)

	for i := 0; i < st.NumField(); i++ {
		// Get the name of the schema field
		key := st.Field(i).Name
		if tag, ok := st.Field(i).Tag.Lookup("name"); ok {
			if tag == "-" {
				continue
			}
			key = tag
		}
		key = strings.ToLower(key)

		// Set the value if one exists
		if schemaVal, ok := resource.GetOk(key); ok {
			field := sv.FieldByName(key)
			reflectVal := reflect.ValueOf(schemaVal)
			field.Set(reflectVal)
		}
	}
	return &instance
}

// FromMachine updates the instance to reflect the state of a Machine
func (i *Instance) FromMachine(m *maas.Machine) *Instance {
	return i
}

// GetMetadata implements the Endpoint interface
func (i *Instance) GetMetadata() interface{} {
	return maas.Machine{}
}

// UpdateState updates the Terraform state to match the Instance state
func (i *Instance) UpdateState(resource *schema.ResourceData) {
	st := reflect.TypeOf(i)
	sv := reflect.ValueOf(i)

	for i := 0; i < st.NumField(); i++ {
		// Get the name of the schema field
		key := st.Field(i).Name
		if tag, ok := st.Field(i).Tag.Lookup("name"); ok {
			if tag == "-" {
				continue
			}
			key = tag
		}
		key = strings.ToLower(key)

		// Set the value on the resource if it is not a zero value
		field := sv.FieldByName(key)
		if field.IsValid() {
			resource.Set(key, field.Interface()) // nolint
		}
	}
}

// AllocateParams creates parameters based on the current value of the Instance
func (i *Instance) AllocateParams() *maas.MachinesAllocateParams {
	var params maas.MachinesAllocateParams
	if i.SystemID != "" {
		params.SystemID = i.SystemID
		return &params
	}
	return nil
}

// DeployParams creates parameters based on the current value of the Instance
func (i *Instance) DeployParams() *maas.MachineDeployParams {
	var params maas.MachineDeployParams
	return &params
}

func (i *Instance) UserDataStateFunc(v interface{}) string {
	switch val := v.(type) {
	case string:
		hash := sha1.Sum([]byte(val)) // nolint: gosec
		return hex.EncodeToString(hash[:])
	default:
		return ""
	}
}
