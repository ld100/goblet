# Golang microservice app prototype

Inspired by [Clean Architecture](https://hackernoon.com/golang-clean-archithecture-efd6d7c43047).

Each folder within src is a package usually focused around exact domain, e.g. users, orders, shipping.
Each package may consist on subpackages representing specific architecture levels:

* models - models structs with includes JSON/YAML/XML serialization instructions. Usually model is a GORM db model, with custom validations, that are run on GORM BeforeSave callback.
* services - Service layer AKA Commands AKA Logic
* controllers/rest - REST-specific server configuration, routes, controllers. Uses CHI framework.
* controllers/grpc - gRPC-specific server configuration, controllers
* forms - forms and form validators
* errors - just typed errors

There are few exceptions however, there are still few non-domain packages:

* src/cmd - command-line utilities
* src/util - project-wise utils
* src/config - project-wise configuration
* migrations - manual database migrations