// Copyright 2016 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package gomaasapi

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/juju/errors"
	"github.com/juju/schema"
	"github.com/juju/version"
)

// Can't use interface as a type, so add an underscore. Yay.
type interface_ struct {
	controller *controller

	resourceURI string

	id      int
	name    string
	type_   string
	enabled bool
	tags    []string

	vlan  *vlan
	links []*link

	macAddress   string
	effectiveMTU int

	parents  []string
	children []string
}

func (i *interface_) updateFrom(other *interface_) {
	i.resourceURI = other.resourceURI
	i.id = other.id
	i.name = other.name
	i.type_ = other.type_
	i.enabled = other.enabled
	i.tags = other.tags
	i.vlan = other.vlan
	i.links = other.links
	i.macAddress = other.macAddress
	i.effectiveMTU = other.effectiveMTU
	i.parents = other.parents
	i.children = other.children
}

// ID implements Interface.
func (i *interface_) ID() int {
	return i.id
}

// Name implements Interface.
func (i *interface_) Name() string {
	return i.name
}

// Parents implements Interface.
func (i *interface_) Parents() []string {
	return i.parents
}

// Children implements Interface.
func (i *interface_) Children() []string {
	return i.children
}

// Type implements Interface.
func (i *interface_) Type() string {
	return i.type_
}

// Enabled implements Interface.
func (i *interface_) Enabled() bool {
	return i.enabled
}

// Tags implements Interface.
func (i *interface_) Tags() []string {
	return i.tags
}

// VLAN implements Interface.
func (i *interface_) VLAN() VLAN {
	if i.vlan == nil {
		return nil
	}
	return i.vlan
}

// Links implements Interface.
func (i *interface_) Links() []Link {
	result := make([]Link, len(i.links))
	for i, link := range i.links {
		result[i] = link
	}
	return result
}

// MACAddress implements Interface.
func (i *interface_) MACAddress() string {
	return i.macAddress
}

// EffectiveMTU implements Interface.
func (i *interface_) EffectiveMTU() int {
	return i.effectiveMTU
}

// UpdateInterfaceArgs is an argument struct for calling Interface.Update.
type UpdateInterfaceArgs struct {
	Name       string
	MACAddress string
	VLAN       VLAN
}

func (a *UpdateInterfaceArgs) vlanID() int {
	if a.VLAN == nil {
		return 0
	}
	return a.VLAN.ID()
}

// Update implements Interface.
func (i *interface_) Update(args UpdateInterfaceArgs) error {
	var empty UpdateInterfaceArgs
	if args == empty {
		return nil
	}
	params := NewURLParams()
	params.MaybeAdd("name", args.Name)
	params.MaybeAdd("mac_address", args.MACAddress)
	params.MaybeAddInt("vlan", args.vlanID())
	source, err := i.controller.put(i.resourceURI, params.Values)
	if err != nil {
		if svrErr, ok := errors.Cause(err).(ServerError); ok {
			switch svrErr.StatusCode {
			case http.StatusNotFound:
				return errors.Wrap(err, NewNoMatchError(svrErr.BodyMessage))
			case http.StatusForbidden:
				return errors.Wrap(err, NewPermissionError(svrErr.BodyMessage))
			}
		}
		return NewUnexpectedError(err)
	}

	response, err := readInterface(i.controller.apiVersion, source)
	if err != nil {
		return errors.Trace(err)
	}
	i.updateFrom(response)
	return nil
}

// Delete implements Interface.
func (i *interface_) Delete() error {
	err := i.controller.delete(i.resourceURI)
	if err != nil {
		if svrErr, ok := errors.Cause(err).(ServerError); ok {
			switch svrErr.StatusCode {
			case http.StatusNotFound:
				return errors.Wrap(err, NewNoMatchError(svrErr.BodyMessage))
			case http.StatusForbidden:
				return errors.Wrap(err, NewPermissionError(svrErr.BodyMessage))
			}
		}
		return NewUnexpectedError(err)
	}
	return nil
}

// InterfaceLinkMode is the type of the various link mode constants used for
// LinkSubnetArgs.
type InterfaceLinkMode string

const (
	// LinkModeDHCP - Bring the interface up with DHCP on the given subnet. Only
	// one subnet can be set to DHCP. If the subnet is managed this interface
	// will pull from the dynamic IP range.
	LinkModeDHCP InterfaceLinkMode = "DHCP"

	// LinkModeStatic - Bring the interface up with a STATIC IP address on the
	// given subnet. Any number of STATIC links can exist on an interface.
	LinkModeStatic InterfaceLinkMode = "STATIC"

	// LinkModeLinkUp - Bring the interface up only on the given subnet. No IP
	// address will be assigned to this interface. The interface cannot have any
	// current DHCP or STATIC links.
	LinkModeLinkUp InterfaceLinkMode = "LINK_UP"
)

// LinkSubnetArgs is an argument struct for passing parameters to
// the Interface.LinkSubnet method.
type LinkSubnetArgs struct {
	// Mode is used to describe how the address is provided for the Link.
	// Required field.
	Mode InterfaceLinkMode
	// Subnet is the subnet to link to. Required field.
	Subnet Subnet
	// IPAddress is only valid when the Mode is set to LinkModeStatic. If
	// not specified with a Mode of LinkModeStatic, an IP address from the
	// subnet will be auto selected.
	IPAddress string
	// DefaultGateway will set the gateway IP address for the Subnet as the
	// default gateway for the machine or device the interface belongs to.
	// Option can only be used with mode LinkModeStatic.
	DefaultGateway bool
}

// Validate ensures that the Mode and Subnet are set, and that the other options
// are consistent with the Mode.
func (a *LinkSubnetArgs) Validate() error {
	switch a.Mode {
	case LinkModeDHCP, LinkModeLinkUp, LinkModeStatic:
	case "":
		return errors.NotValidf("missing Mode")
	default:
		return errors.NotValidf("unknown Mode value (%q)", a.Mode)
	}
	if a.Subnet == nil {
		return errors.NotValidf("missing Subnet")
	}
	if a.IPAddress != "" && a.Mode != LinkModeStatic {
		return errors.NotValidf("setting IP Address when Mode is not LinkModeStatic")
	}
	if a.DefaultGateway && a.Mode != LinkModeStatic {
		return errors.NotValidf("specifying DefaultGateway for Mode %q", a.Mode)
	}
	return nil
}

// LinkSubnet implements Interface.
func (i *interface_) LinkSubnet(args LinkSubnetArgs) error {
	if err := args.Validate(); err != nil {
		return errors.Trace(err)
	}
	params := NewURLParams()
	params.Values.Add("mode", string(args.Mode))
	params.Values.Add("subnet", fmt.Sprint(args.Subnet.ID()))
	params.MaybeAdd("ip_address", args.IPAddress)
	params.MaybeAddBool("default_gateway", args.DefaultGateway)
	source, err := i.controller.post(i.resourceURI, "link_subnet", params.Values)
	if err != nil {
		if svrErr, ok := errors.Cause(err).(ServerError); ok {
			switch svrErr.StatusCode {
			case http.StatusNotFound, http.StatusBadRequest:
				return errors.Wrap(err, NewBadRequestError(svrErr.BodyMessage))
			case http.StatusForbidden:
				return errors.Wrap(err, NewPermissionError(svrErr.BodyMessage))
			case http.StatusServiceUnavailable:
				return errors.Wrap(err, NewCannotCompleteError(svrErr.BodyMessage))
			}
		}
		return NewUnexpectedError(err)
	}

	response, err := readInterface(i.controller.apiVersion, source)
	if err != nil {
		return errors.Trace(err)
	}
	i.updateFrom(response)
	return nil
}

func (i *interface_) linkForSubnet(subnet Subnet) *link {
	for _, link := range i.links {
		if s := link.Subnet(); s != nil && s.ID() == subnet.ID() {
			return link
		}
	}
	return nil
}

// LinkSubnet implements Interface.
func (i *interface_) UnlinkSubnet(subnet Subnet) error {
	if subnet == nil {
		return errors.NotValidf("missing Subnet")
	}
	link := i.linkForSubnet(subnet)
	if link == nil {
		return errors.NotValidf("unlinked Subnet")
	}
	params := NewURLParams()
	params.Values.Add("id", fmt.Sprint(link.ID()))
	source, err := i.controller.post(i.resourceURI, "unlink_subnet", params.Values)
	if err != nil {
		if svrErr, ok := errors.Cause(err).(ServerError); ok {
			switch svrErr.StatusCode {
			case http.StatusNotFound, http.StatusBadRequest:
				return errors.Wrap(err, NewBadRequestError(svrErr.BodyMessage))
			case http.StatusForbidden:
				return errors.Wrap(err, NewPermissionError(svrErr.BodyMessage))
			}
		}
		return NewUnexpectedError(err)
	}

	response, err := readInterface(i.controller.apiVersion, source)
	if err != nil {
		return errors.Trace(err)
	}
	i.updateFrom(response)

	return nil
}

// CreateVLANInterfaceArgs is an argument struct for passing parameters to
// the Machine.CreateVLANInterface method.
type CreateVLANInterfaceArgs struct {
	// VLAN is the untagged VLAN the interface is connected to (required).
	VLAN VLAN
	// Tags to attach to the interface (optional).
	Tags []string
	// MTU - Maximum transmission unit. (optional)
	MTU int
	// AcceptRA - Accept router advertisements. (IPv6 only)
	AcceptRA bool
	// Autoconf - Perform stateless autoconfiguration. (IPv6 only)
	Autoconf bool
}

// Validate checks the required fields are set for the arg structure.
func (a *CreateVLANInterfaceArgs) Validate() error {
	if a.VLAN == nil {
		return errors.NotValidf("missing VLAN")
	}
	return nil
}

// interfacesURI used to add child interfaces to this interface.  We need to strip the interface ID in URI to make post call
func (i *interface_) interfacesURI() string {
	return strings.Replace(i.resourceURI, "/"+strconv.Itoa(i.ID())+"/", "/", 1)
}

// CreateInterface implements Device.
func (i *interface_) CreateVLANInterface(args CreateVLANInterfaceArgs) (Interface, error) {
	if err := args.Validate(); err != nil {
		return nil, errors.Trace(err)
	}
	params := NewURLParams()
	params.Values.Add("vlan", fmt.Sprint(args.VLAN.ID()))
	params.MaybeAddInt("parent", i.ID())
	params.MaybeAdd("tags", strings.Join(args.Tags, ","))
	params.MaybeAddInt("mtu", args.MTU)
	params.MaybeAddBool("accept_ra", args.AcceptRA)
	params.MaybeAddBool("autoconf", args.Autoconf)
	result, err := i.controller.post(i.interfacesURI(), "create_vlan", params.Values)
	if err != nil {
		if svrErr, ok := errors.Cause(err).(ServerError); ok {
			switch svrErr.StatusCode {
			case http.StatusNotFound, http.StatusConflict:
				return nil, errors.Wrap(err, NewBadRequestError(svrErr.BodyMessage))
			case http.StatusForbidden:
				return nil, errors.Wrap(err, NewPermissionError(svrErr.BodyMessage))
			case http.StatusServiceUnavailable:
				return nil, errors.Wrap(err, NewCannotCompleteError(svrErr.BodyMessage))
			}
		}
		return nil, NewUnexpectedError(err)
	}

	iface, err := readInterface(i.controller.apiVersion, result)
	if err != nil {
		return nil, errors.Trace(err)
	}
	iface.controller = i.controller

	// TODO: add to the interfaces for the device when the interfaces are returned.
	// lp:bug 1567213.
	return iface, nil
}

func readInterface(controllerVersion version.Number, source interface{}) (*interface_, error) {
	readFunc, err := getInterfaceDeserializationFunc(controllerVersion)
	if err != nil {
		return nil, errors.Trace(err)
	}

	checker := schema.StringMap(schema.Any())
	coerced, err := checker.Coerce(source, nil)
	if err != nil {
		return nil, WrapWithDeserializationError(err, "interface base schema check failed")
	}
	valid := coerced.(map[string]interface{})
	return readFunc(valid)
}

func readInterfaces(controllerVersion version.Number, source interface{}) ([]*interface_, error) {
	readFunc, err := getInterfaceDeserializationFunc(controllerVersion)
	if err != nil {
		return nil, errors.Trace(err)
	}

	checker := schema.List(schema.StringMap(schema.Any()))
	coerced, err := checker.Coerce(source, nil)
	if err != nil {
		return nil, WrapWithDeserializationError(err, "interface base schema check failed")
	}
	valid := coerced.([]interface{})
	return readInterfaceList(valid, readFunc)
}

func getInterfaceDeserializationFunc(controllerVersion version.Number) (interfaceDeserializationFunc, error) {
	var deserialisationVersion version.Number
	for v := range interfaceDeserializationFuncs {
		if v.Compare(deserialisationVersion) > 0 && v.Compare(controllerVersion) <= 0 {
			deserialisationVersion = v
		}
	}
	if deserialisationVersion == version.Zero {
		return nil, NewUnsupportedVersionError("no interface read func for version %s", controllerVersion)
	}
	return interfaceDeserializationFuncs[deserialisationVersion], nil
}

func readInterfaceList(sourceList []interface{}, readFunc interfaceDeserializationFunc) ([]*interface_, error) {
	result := make([]*interface_, 0, len(sourceList))
	for i, value := range sourceList {
		source, ok := value.(map[string]interface{})
		if !ok {
			return nil, NewDeserializationError("unexpected value for interface %d, %T", i, value)
		}
		read, err := readFunc(source)
		if err != nil {
			return nil, errors.Annotatef(err, "interface %d", i)
		}
		result = append(result, read)
	}
	return result, nil
}

type interfaceDeserializationFunc func(map[string]interface{}) (*interface_, error)

var interfaceDeserializationFuncs = map[version.Number]interfaceDeserializationFunc{
	twoDotOh: interface_2_0,
}

func interface_2_0(source map[string]interface{}) (*interface_, error) {
	fields := schema.Fields{
		"resource_uri": schema.String(),

		"id":      schema.ForceInt(),
		"name":    schema.String(),
		"type":    schema.String(),
		"enabled": schema.Bool(),
		"tags":    schema.OneOf(schema.Nil(""), schema.List(schema.String())),

		"vlan":  schema.OneOf(schema.Nil(""), schema.StringMap(schema.Any())),
		"links": schema.List(schema.StringMap(schema.Any())),

		"mac_address":   schema.OneOf(schema.Nil(""), schema.String()),
		"effective_mtu": schema.ForceInt(),

		"parents":  schema.List(schema.String()),
		"children": schema.List(schema.String()),
	}
	defaults := schema.Defaults{
		"mac_address": "",
	}
	checker := schema.FieldMap(fields, defaults)
	coerced, err := checker.Coerce(source, nil)
	if err != nil {
		return nil, WrapWithDeserializationError(err, "interface 2.0 schema check failed")
	}
	valid := coerced.(map[string]interface{})
	// From here we know that the map returned from the schema coercion
	// contains fields of the right type.

	var vlan *vlan
	// If it's not an attribute map then we know it's nil from the schema check.
	if vlanMap, ok := valid["vlan"].(map[string]interface{}); ok {
		vlan, err = vlan_2_0(vlanMap)
		if err != nil {
			return nil, errors.Trace(err)
		}
	}

	links, err := readLinkList(valid["links"].([]interface{}), link_2_0)
	if err != nil {
		return nil, errors.Trace(err)
	}
	macAddress, _ := valid["mac_address"].(string)
	result := &interface_{
		resourceURI: valid["resource_uri"].(string),

		id:      valid["id"].(int),
		name:    valid["name"].(string),
		type_:   valid["type"].(string),
		enabled: valid["enabled"].(bool),
		tags:    convertToStringSlice(valid["tags"]),

		vlan:  vlan,
		links: links,

		macAddress:   macAddress,
		effectiveMTU: valid["effective_mtu"].(int),

		parents:  convertToStringSlice(valid["parents"]),
		children: convertToStringSlice(valid["children"]),
	}
	return result, nil
}
