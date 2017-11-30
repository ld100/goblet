# Golang microservice app prototype

Inspired by [Clean Architecture](https://hackernoon.com/golang-clean-archithecture-efd6d7c43047).

Each folder within src is a package usually focused around exact domain, e.g. users, orders, shipping.
Each package may consist on subpackages representing specific architecture levels:

* models - models structs with includes JSON/YAML/XML serialization instructions
* validations models validation
* repository - DAO objects
* services - Service layer AKA Commands AKA Logics
* controllers/rest - REST-specific server configuration, routes, controllers
* controllers/grpc - gRPC-specific server configuration, controllers
* controllers - nuff said
* forms - forms and form validators
* errors - just typed errors

There are few exceptions however, there are still few non-domain packages:

* src/cmd - command-line utilities
* src/util - project-wise utils
* src/config - project-wise configuration
* migrations - database migrations