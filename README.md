# Golang microservice app prototype

Inspired by [Clean Architecture](https://hackernoon.com/golang-clean-archithecture-efd6d7c43047).

Each folder within src is a package usually focused around exact domain, e.g. users, orders, shipping.
Each package may consist on subpackages representing specific architecture levels:

* models
* validations (model validators)
* repository - DAO objects
* services - Service layer AKA Commands AKA Logics
* rest - REST-specific server configuration and routes
* grpc - gRPC-specific server configuration
* controllers - nuff said
* forms - forms and form validators

There are few exceptions however, there are still few non-domain packages:

* src/cmd - command-line utilities
* src/util - project-wise utils
* src/config - project-wise configuration
* migrations - database migrations