# terraform-provider-maas

## Description
A simple Terraform provider for MAAS.  This is a work in progress and by no means should be considered production quality work.  The current version supports the allocation, power up, power down and release of nodes already registered with MAAS.  I think this is the main usage for MAAS and will cover the majority of use cases out there.  I'll look into adding more functionality in the future.

## Requirements
* [Terraform](https://github.com/hashicorp/terraform)
* [Go MAAS API Library](https://github.com/juju/gomaasapi)

## Usage

### Provider Configuration
The provider requires some variables to be configured in order to gain access to the MAAS server:

* **api_version**:  This is optional and probably only works with 2.0. The defaults to 2.0.
* **api_key**: MAAS API Key (Details: https://maas.ubuntu.com/docs/maascli.html#logging-in)
* **api_url**: URI for your MAAS API server.  ie: http://127.0.0.1:80/MAAS

#### `maas`
```
provider "maas" {
    api_version = "2.0"
    api_key = "YOUR MAAS API KEY"
    api_url = "http://<MAAS_SERVER>[:MAAS_PORT]/MAAS"
}
```

## Resource Configuration
#### `maas_machine`
#### Example
```
resource "maas_machine" "maas_single_random_node" {
	count = 1
}
```
### Argument Reference
The selection mechanism for the machines is a subset of criteria described in the MAAS API (https://docs.maas.io/2.5/en/api#machines).  

- **hostname** - (Optional) Host name to try to allocate
- **system_id** - (Optional) System ID the requested machine
- **architecture** - (Optional) System's archirtecture eg: amd64/generic
- **cpu_count** - (Optional) The minimum number of cpu cores needed for consideration
- **memory** - (Optional) Minimum amount of RAM (expressed in MB) neede for consideration 
- **tags** - (Optional) List of tags to use in the selection process
- **not_tags** - (Optional) Tags the machine must NOT match
- **zone** - (Optional) Availability Zone name the machine must be located in
- **not_in_zones** - (Optional) List of Availability Zones from which the machine must not be acquired.  If multiple zones are specified, the machine must NOT be associated with ANY of them. 
- **comment** - (Optional) Comment for the machine event log
- **deploy** - (Optional) Immediately start deploying machine.  If `False` then only Allocate machine, do not Deploy.  Changing from `False` to `True` will cause the machine to Deploy.  Changing from `True` to `False` will Release the machine and immediately attempt to Allocate.  Default: `True`
- **sticky_hostname** - (Optional) If `True` then the hostname used for Deploy will stick as the machine's name in MAAS dashboard, otherwise the machine's name will revert to the pre-Deploy hostname.  Default: `False`

The above constraints parameters can be used to acquire a machine that possesses certain characteristics. All the constraints are optional and when multiple constraints are provided, they are combined using ‘AND’ semantics.  In the absence of any constraints, a random node will be selected and deployed.  

### Attributes Reference

Currently all the arguments are retained as attributes.  That is, the plug-in does not read the machine's attributes from the MAAS server, and instead simply stores the arguments supplied from the last `apply` into the state file as its attributes.
#### Deploy three random nodes
```
resource "maas_machine" "maas_three_random_nodes" {
	count = 3
}
```

#### Deploy a node named "node-1"
```
resource "maas_machine" "maas_node_1" {
	hostname = "node-1"
}
```

#### Deploy three nodes with at least 8G of RAM
```
resource "maas_machine" "maas_three_nodes_8g" {
	memory = "8000"
	count = 3
}
```

#### Specify user data for nodes

User data can be either a cloud-init script or a bash shell

Header for cloud-init:
```
#cloud-config
```

Header for script (shebang):
```
#!/bin/bash
```

Example (read from file):
```
resource "maas_machine" "maas_single_random_node" {
    count = 1
    
    user_data = "${file("${path.module}/user_data/test_data.txt")}"
}
```

#### Specify a comment in the event log
```
resource "maas_machine" "maas_single_random_node" {
    count = 1
    
    comment = "Platform deployment"
}
```

#### Use tags to restrict deployments to specific nodes
```
resource "maas_machine" "maas_single_random_node" {
    count = 1
    
    tags = ["DELL_R630", "APP_CLASS"]
}
```

#### Specify the hostname for the deployed node
```
resource "maas_machine" "maas_single_random_node" {
    count = 1
    
    deploy_hostname = "freedompants"
}
```

#### Specify tags for the deployed node
```
resource "maas_machine" "maas_single_random_node" {
    count = 1
    
    deploy_tags = ["hostwiththemost", "platform"]
}
```

#### Select distro for a node
Useful for custom OS builds
```
resource "maas_machine" "maas_single_random_node" {
    count = 1

    distro_series = "centos73" 
}
```

## Erasing disks on node release

Maas provides an option to erase the node's disk when releasing the system. By default it will not alter the disk.
This provides a very quick method do release the system back into the pool of nodes. It isn't ideal to leave data on a disk
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

```
resource "maas_machine" "maas_single_random_node" {
    count = 1
}
```

This shows what is set by default in Terraform. You are not required to set this option.

```
resource "maas_machine" "maas_single_random_node" {
    count = 1

    release_erase = true
}
```

How to disable the disk erasure.

```
resource "maas_machine" "maas_single_random_node" {
    count = 1
    
    release_erase = false
}
```

### Secure erase option


```
resource "maas_machine" "maas_single_random_node" {
    count = 1

    release_erase_secure = true
}
```

### Quick erase option

```
resource "maas_machine" "maas_single_random_node" {
    count = 1

    release_erase_quick = true
}
```

If there are conflicting options, such as enabling both secure and quick erase, this is how the Maas API deals with conflicts.

If neither release_secure_erase nor release_quick_erase are specified, MAAS will overwrite the whole disk with null bytes. This can be very slow.

If both release_secure_erase and release_quick_erase are specified and the drive does NOT have a secure erase feature, MAAS will behave as if only quick_erase was specified.

If release_secure_erase is specified and release_quick_erase is NOT specified and the drive does NOT have a secure erase feature, MAAS will behave as if secure_erase was NOT specified, i.e. will overwrite the whole disk with null bytes. This can be very slow.

Source: [Maas API: POST /api/2.0/machines/{system_id}/ op=release](https://docs.ubuntu.com/maas/2.1/en/api)

### The future
This is just a basic (proof of concept) provider.  The following are some of the features I would like to see here:

* All of the supported constratins for allocating and deploying a node
* Discover nodes
* Create new nodes
* Accept and commission nodes
