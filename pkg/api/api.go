/*package api defines an interface to each MaaS API endpoint.

Each interface correlates to one endpoint, such as Subnets for the Subnets
endpoint (ie /subnets) and Subnet for the Subnet endpoint (eg subnets/<subnet_id>).
API clients are expected to implement these interfaces to provide a normalized way
of accessing the MaaS API with normalized results (eg the types defined in the
entity package).

Some endpoint operations require multiple parameters, such as the Rack Controllers
GET operation, which takes a number of QSP that can be used to filter results. These
parameters are encapsulated in the params subpackage, providing a quick reference
for performing API operations.*/
package api
