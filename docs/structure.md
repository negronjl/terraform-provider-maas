# Application Structure

This application uses `pkg/` and `internal/` directories to separate public (ie importable) code from private code. The code itself generally falls into one of two categories:

- MaaS API
- Terraform provider

Where the former is code designed to interact with the MaaS API, and the latter is code designed to work with the Terraform provider. All of the MaaS API code is _importable_ and lives in the packages under the `pkg/` directory, and all of the Terraform provider code is specific to this repository and lives in the `internal/` directory.

The `test/` directory contains a package of helper functions to aid in testing, and the `testdata/` subdirectory contains sample API responses taken from the [API documentation](https://maas.io/docs/api).

## MaaS API Code

This repository uses JuJu's [MaaS API library](https://github.com/juju/gomaasapi) under the hood to make API requests. The `pkg/gmaw` package is a wrapper around this library, providing a standardized (and DRY) code API to the web API.

- Each type defined in `pkg/gmaw` corresponds to an API endpoint, and fulfills a similarly-named interface in `pkg/api`.
- Each method in the above types corresponds to an API operation, such as `Commission`ing a `Machine` or `Delete`ing a `Subnet`.

Methods that call operations that return data will provide the data in a structured format as defined by their corresponding type in `pkg/maas/entity`, with per-endpoint subpackages for operations that return metadata such as the `GetIPAddresses()` of the `Subnet` endpoint. Some endpoint actions, such as those that `POST` to create a new resource or `PUT` to update an existing resource, have a significant number of potential options: these are encapsulated in the types in the `pkg/api/params` package.

While the `pkg/gmaw` package provides direct access to the API, it does not manage state, nor does it have a way to determine, for example, when a long-running task (such as commissioning a machine) has completed. The `pkg/maas/manager` package provides types that are analogous to those in `pkg/gmaw` with a slightly abstracted interface for simpler management.

| Package | Purpose |
| ------- | ------- |
| `/pkg/api` | Defines interfaces that provide a standardized `<Endpoint>.<Operation>(<params>...)` interface to the web API
| `/pkg/api/params` | Provides types with attributes for each parameter for operations with many parameters
| `/pkg/gmaw` | Wraps `juju/gomaasapi` to provide an API client that implements the interfaces defined in `/pkg/api`
| `/pkg/maas/entity` | Contains typed schemas for each major type (RackController, Interface, Machine, etc)
| `/pkg/maas/entity/*` | Defines additional, per-endpoint types for operations that return metadata
| `/pkg/maas/manager` | Encapsulates `/pkg/gmaw` with state management and locking capabilities.

## Terraform Provider code

Each Terraform `resource`/`data` is defined in the `/internal/provider` package, with a corresponding type in the `/internal/tfschema` package. The `tfschema` types are necessary because the Terraform schema is not analogous to the MaaS API schema: for example, an `entity.Subnet` embeds the `entity.VLAN` to which it belongs, while a `tfschema.Subnet` only contains the name of the VLAN. Another example is the `Machine` entity, which contains information about the number of CPU cores and amount of memory the chassis contains; these values cannot be changed via IaC - YET!

The `/internal/bridge` package bridges the gap between Terraform and MaaS. This is a Hashicorp [best practice](https://www.terraform.io/docs/extend/writing-custom-providers.html#dedicated-upstream-libraries), and especially necessary to translate between Terraform's state management and MaaS's verbial, operation-oriented REST endpoints.

| Package | Purpose |
| ------- | ------- |
| `/internal/tfschema` | Defines the schema for each Terraform `resource` and `data`
| `/internal/provider` | Defines the Terraform resources
| `/internal/bridge` | Contains types that translate between Terraform and MaaS
