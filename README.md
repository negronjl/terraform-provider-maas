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

### Resource Configuration (maas_instance)
This provider is only able to deploy and release nodes already registered and configured in MAAS.  The selection mechanism for the nodes is a subset of criteria described in the MAAS API (https://maas.ubuntu.com/docs/api.html#nodes).  Currently, this provider supports:

- **hostnames**: Host name to try to allocate.
- **architecture**: Architecture of the requested machine: ie: amd64/generic
- **cpu_count**: The minimum number of cpu cores needed for consideration
- **memory**: Minimum amount of RAM neede for consideration
- **tag_names**: List of tags to use in the selection process ( Experimental )

The above constraints parameters can be used to acquire a node that possesses certain characteristics. All the constraints are optional and when multiple constraints are provided, they are combined using ‘AND’ semantics.  In the absence of any constraints, a random node will be selected and deployed.  The examples in the next section attempt to explain how to use the resource.

#### `maas_instance`
##### Deploy a Random node
```
resource "maas_instance" "maas_single_random_node" {
	count = 1
}
```

##### Deploy three random nodes
```
resource "maas_instance" "maas_three_random_nodes" {
	count = 3
}
```

##### Deploy a node named "node-1"
```
resource "maas_instance" "maas_node_1" {
	hostname = "node-1"
}
```

##### Deploy three nodes with at least 8G of RAM
```
resource "maas_instance" "maas_three_nodes_8g" {
	memory = "8G"
	count = 3
}
```
### The future
This is just a basic (proof of concept) provider.  The following are some of the features I would like to see here:

* All of the supported constratins for allocating and deploying a node
* Discover nodes
* Create new nodes
* Accept and commission nodes
