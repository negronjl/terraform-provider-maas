/*package entity defines types for the MaaS API endpoints' return types.

Each endpoint returns JSON that describes the object it represents. For example,
GETting the Subnets endpoint will return an array of Subnet, while GETting the
Subnet endpoint (ie subnets/<subnet_id>) will return one Subnet.

Some endpoints expose operations that return metadata about the object, such as as
Subnet's GetStatistics(), which contains statistics about the subnet, but does not
actually describe the subnet: these alternative types can be found in subpackages
named after the endpoint.*/
package entity
