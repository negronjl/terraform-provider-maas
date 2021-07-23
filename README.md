# terraform-provider-maas

## Description

A simple Terraform provider for MAAS.  This is a work in progress and by no means should be considered production quality work.  The current version supports the allocation, power up, power down and release of nodes already registered with MAAS.  I think this is the main usage for MAAS and will cover the majority of use cases out there.  I'll look into adding more functionality in the future.

## Requirements

- [Terraform](https://github.com/hashicorp/terraform)
- [Go MAAS API Library](https://github.com/juju/gomaasapi)

## Usage

### Provider Configuration

The provider requires some variables to be configured in order to gain access to the MAAS server:

- **api_version**:  This is optional and probably only works with 2.0. The defaults to 2.0.
- **api_key**: [MAAS API Key](https://maas.ubuntu.com/docs/maascli.html#logging-in), specify with the `MAAS_API_KEY` environment variable
- **api_url**: URI for your MAAS API server (eg <http://127.0.0.1:80/MAAS>)

#### `maas`

For most setups, you should save your API token as the `MAAS_API_KEY` in your shell, then
just hardcode the server's URI in the provider definition:

```hcl
provider "maas" {
    api_url = "http://<MAAS_SERVER>[:MAAS_PORT]/MAAS"
}
```

To hardcode everything:

```hcl
provider "maas" {
  api_version = "2.0"
  api_key = "YOUR MAAS API KEY"
  api_url = "http://<MAAS_SERVER>[:MAAS_PORT]/MAAS"
}
```

### Resource Configuration (maas_instance)

This provider is only able to deploy and release nodes already registered and configured in MAAS.  The selection mechanism for the nodes is a subset of criteria described in the [MAAS API]<https://maas.ubuntu.com/docs/api.html#nodes>.  Currently, this provider supports:

- **hostnames**: Host name to try to allocate.
- **architecture**: Architecture of the requested machine: ie: amd64/generic
- **cpu_count**: The minimum number of cpu cores needed for consideration
- **memory**: Minimum amount of RAM neede for consideration
- **tags**: List of tags to use in the selection process

The above constraints parameters can be used to acquire a node that possesses certain characteristics. All the constraints are optional and when multiple constraints are provided, they are combined using ‘AND’ semantics.  In the absence of any constraints, a random node will be selected and deployed.  The examples in the next section attempt to explain how to use the resource.

#### `maas_instance`

##### Deploy a Random node

```hcl
resource "maas_instance" "maas_single_random_node" {
  count = 1
}
```

##### Deploy three random nodes

```hcl
resource "maas_instance" "maas_three_random_nodes" {
  count = 3
}
```

##### Deploy a node named "node-1"

```hcl
resource "maas_instance" "maas_node_1" {
  hostname = "node-1"
}
```

##### Deploy three nodes with at least 8G of RAM

```hcl
resource "maas_instance" "maas_three_nodes_8g" {
  memory = "8G"
  count = 3
}
```

#### maas_interface_physical

Configures a physical interface on a system, where a system is anything with a system ID.

```hcl
resource "maas_interface_physical" "myserver_eth0" {
  system_id   = maas_instance.myserver.system_id
  name        = "eth0"
  mac_address = "01:23:45:67:89:0A"
  tags        = ["foo", "bar"]
  vlan        = "my_vlan"
  mtu         = 1500
  accept_ra   = true
  autoconf    = true
}
```

##### Available Parameters

| Name | Type | Description
| ---- | ---- | -----------
| `system_id` | `string` | The system ID of the system to which this interface belongs
| `name` | `string` | Name of the interface (such as `eth0`)
| `mac_address` | `string` | MAC Address of the interface. This can be changed.
| `tags` | `list(string)` | Tags to apply to the interface
| `vlan` | `string` | Name of the VLAN to which this interface is connected. If empty or `undefined`, the interface is assumed to be disconnected.
| `mtu` | `int` | MTU of the interface
| `accept_ra` | `bool` | Accept router advertisements (IPv6 only). Default false.
| `autoconf` | `bool` | Use autoconf (IPv6 only). Default false.

The `system_id` and `mac_address` parameters are required.

##### Importing

An interface can be uniquely identified by its system ID and interface id (which is an integer).

```bash
terraform import maas_interface_physical.my_eth 3xtkyg:23
```

#### maas_interface_link

Configures a link between physical interface on a system and a subnet.

```hcl
resource "maas_interface_link" "eth0_sn123" {
  system_id    = maas_interface_physical.myserver.system_id
  interface_id = maas_interface_physical.myserver.interface_id
  subnet_id    = 123
  mode         = "STATIC"
  ip_address   = "::1"
  force        = true
}
```

##### Available Parameters

| Name | Type | Description
| ---- | ---- | -----------
| `system_id` | `string` | The system ID of the system to which this interface belongs
| `interface_id` | `int` | ID of the interface being linked
| `subnet_id` | `int` | ID of the subnet being linked
| `mode` | `string` | One of AUTO, DHCP, STATIC or LINK_UP connection to subnet.
| `ip_address` | `string` | Optional IP address for the interface in subnet. Only used when mode is STATIC. If not provided an IP address from subnet will be auto selected.
| `force` | `bool` | If True, allows LINK_UP to be set on the interface even if other links already exist. Also allows the selection of any VLAN, even a VLAN MAAS does not believe the interface to currently be on. Using this option will cause all other links on the interface to be deleted. (Defaults to False.)
| `default_gateway` | `string` | True sets the gateway IP address for the subnet as the default gateway for the node this interface belongs to. Option can only be used with the AUTO and STATIC modes.

Mode definitions:

- **AUTO**: Assign this interface a static IP address from the provided subnet. The subnet must be a managed subnet. The IP address will not be assigned until the node goes to be deployed.
- **DHCP**: Bring this interface up with DHCP on the given subnet. Only one subnet can be set to DHCP. If the subnet is managed this interface will pull from the dynamic IP range.
- **STATIC**: Bring this interface up with a static IP address on the given subnet. Any number of static links can exist on an interface.
- **LINK_UP**: Bring this interface up only on the given subnet. No IP address will be assigned to this interface. The interface cannot have any current AUTO, DHCP or STATIC links.

The `system_id`, `interface_id`, and `subnet_id` parameters are required.

##### Importing

A link can be uniquely identified by a combination of its system ID, interface ID, and the ID of the subnet connected by the link.

```bash
terraform import maas_interface_link.my_link 3xtkyg:23:42
```

#### maas_server

Configure MaaS server parameters.

| Name | Type | Description
| ---- | ---- | -----------
| `ntp_servers` | `list(string)` | Addresses of NTP servers. NTP servers, specified as IP addresses or hostnames delimited by commas and/or spaces, to be used as time references for MAAS itself, the machines MAAS deploys, and devices that make use of MAAS's DHCP services.

This resource currently only supports configuring NTP servers.

#### data.maas_subnet

Search the MaaS API for a subnet. If there are multiple matches, the first one will be returned.

```hcl
data "maas_subnet" "mynet" {
  name = "the_subnet"
  vlan = 123
  cidr = "192.168.0.0/24"
}
```

##### Available Parameters

| Name | Type | Description
| ---- | ---- | -----------
| `name` | `string` | Name of the subnet
| `vlan` | `int` | ID of a VLAN
| `cidr` | `int` | Subnet of the subnet in CIDR notation

All parameters are optional. The first subnet that matches all specified parameters will be returned.

##### Additional Properties

Besides the properties defined above, the following properties are also available:

| Name | Type | Description
| ---- | ---- | -----------
| `id` | `int` | The ID of the subnet
| `rdns_mode` | `int` | How reverse DNS is handled for this subnet
| `gateway_ip` | `string` | The gateway IP address for this subnet.
| `dns_servers` | `list(string` | List of DNS servers for this subnet.

rdns_mode values:

- `0` Disabled: No reverse zone is created.
- `1` Enabled: Generate reverse zone.
- `2` RFC2317: Extends '1' to create the necessary parent zone with the appropriate CNAME resource records for the network, if the the network is small enough to require the support described in RFC2317.

#### data.maas_rack_controller

Search the MaaS API for a rack controller. If there are multiple matches, the first one will be returned.

```hcl
data "maas_rack_controller" "ctrl" {
  domain = ".example.com"
  zone = "the_zone"
  pool = "foo"
}
```

##### Available Parameters

| Name | Type | Description
| ---- | ---- | -----------
| `hostname` | `string` | Only nodes relating to the node with the matching hostname will be returned.
| `mac_address` | `string` | Only nodes relating to the node owning the specified MAC address will be returned.
| `system_id` | `string` | Only nodes relating to the nodes with matching system ids will be returned.
| `domain` | `string` | Only nodes relating to the nodes in the domain will be returned.
| `zone` | `string` | Only nodes relating to the nodes in the zone will be returned.
| `pool` | `string` | Only nodes belonging to the pool will be returned.
| `agent_name` | `string` | Only nodes relating to the nodes with matching agent names will be returned.

All parameters are optional. The first subnet that matches all specified parameters will be returned.

### Specify user data for nodes

User data can be either a cloud-init script or a bash shell

Header for cloud-init:

```plain
#cloud-config
```

Header for script (shebang):

```plain
#!/bin/bash
```

Example (read from file):

```hcl
resource "maas_instance" "maas_single_random_node" {
  count = 1
  user_data = "${file("${path.module}/user_data/test_data.txt")}"
}
```

### Specify a comment in the event log

```hcl
resource "maas_instance" "maas_single_random_node" {
  count = 1
  comment = "Platform deployment"
}
```

### Use tags to restrict deployments to specific nodes

```hcl
resource "maas_instance" "maas_single_random_node" {
  count = 1
  tags = ["DELL_R630", "APP_CLASS"]
}
```

### Specify the hostname for the deployed node

```hcl
resource "maas_instance" "maas_single_random_node" {
  count = 1
  deploy_hostname = "freedompants"
}
```

### Specify tags for the deployed node

```hcl
resource "maas_instance" "maas_single_random_node" {
  count = 1
  deploy_tags = ["hostwiththemost", "platform"]
}
```

### Select distro for a node

Useful for custom OS builds

```hcl
resource "maas_instance" "maas_single_random_node" {
  count = 1
  distro_series = "centos73"
}
```

### Select ubuntu kernel release (HWE kernel version)

Useful for choosing a particular kernel release train builds

```hcl
resource "maas_instance" "maas_single_random_node" {
  count = 1
  distro_series = "bionic"
  hwe_kernel = "hwe-18.04"
}
```

## Erasing disks on node release

Maas provides an option to erase the node's disk when releasing the system. By default it will not alter the disk.
This provides a very quick method to release the system back into the pool of nodes. It isn't ideal to leave data on a disk
as this may lead to data loss or even booting a system that may cause a service outage. With this in mind the
Terraform provider is set to erase the disk on release. This ensures that the machine will be released into the pool with a clean state.

There are a few options when releasing a system:

- erase
  - The default setting
  - MAAS will overwrite the whole disk with null bytes. This can be very slow.
  - Estimated 20min
- secure erase
  - Requires the disk to support a secure erase option.
  - If the disk does not support secure erase it will default the erase option. MAAS will overwrite the whole disk with null bytes. This can be very slow.
  - Estimated 20min
- quick erase
  - Wipe 1MiB at the start and at the end of the drive to make data recovery inconvenient and unlikely to happen by accident. This is not secure.
  - Estimated 3min

### Using the erase feature

### Default erase option

The default option is to always perform an erase.

```hcl
resource "maas_instance" "maas_single_random_node" {
  count = 1
}
```

This shows what is set by default in Terraform. You are not required to set this option.

```hcl
resource "maas_instance" "maas_single_random_node" {
  count = 1
  release_erase = true
}
```

How to disable the disk erasure.

```hcl
resource "maas_instance" "maas_single_random_node" {
  count = 1
  release_erase = false
}
```

### Secure erase option

```hcl
resource "maas_instance" "maas_single_random_node" {
  count = 1
  release_erase_secure = true
}
```

### Quick erase option

```hcl
resource "maas_instance" "maas_single_random_node" {
  count = 1
  release_erase_quick = true
}
```

### Build kvm server

```hcl
resource "maas_instance" "maas_single_random_node" {
  count = 1
  install_kvm = true
}
```

### Build rack controller

```hcl
resource "maas_instance" "maas_single_random_node" {
  count = 1
  install_rackd = true
}
```

If there are conflicting options, such as enabling both secure and quick erase, this is how the Maas API deals with conflicts.

If neither release_secure_erase nor release_quick_erase are specified, MAAS will overwrite the whole disk with null bytes. This can be very slow.

If both release_secure_erase and release_quick_erase are specified and the drive does NOT have a secure erase feature, MAAS will behave as if only quick_erase was specified.

If release_secure_erase is specified and release_quick_erase is NOT specified and the drive does NOT have a secure erase feature, MAAS will behave as if secure_erase was NOT specified, i.e. will overwrite the whole disk with null bytes. This can be very slow.

Source: [Maas API: POST /api/2.0/machines/{system_id}/ op=release](https://docs.ubuntu.com/maas/2.1/en/api)

### The future

This is just a basic (proof of concept) provider.  The following are some of the features I would like to see here:

- All of the supported constratins for allocating and deploying a node
- Discover nodes
- Create new nodes
- Accept and commission nodes
