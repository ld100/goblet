# Golang microservice app prototype

Inspired by [Clean Architecture](https://hackernoon.com/golang-clean-archithecture-efd6d7c43047).

Application consists of configuration and Golang code. Configuration is concentrated in the root folder and `config`.
Go code consists of common use packages and domain-specific ones. Common use packages are:

* cmd − command-line tools
* server − entry point for server-related code. Consists of REST and gRPC servers.
* util − project-wide utility code, including database, cryptography, data migration, environment, etc

`domain` includes domain-specific packages, e.g. users, orders, articles. Each domain package may consist of the following parts:

* models − models structs with JSON/YAML/XML serialization instructions. Usually model is a GORM db model with validations.
* services − Service layer AKA Commands AKA Logic
* rest − REST-specific routes controllers, middlewares. Uses CHI framework.
* grpc − gRPC-specific controllers and middlewares.
* forms − forms and form validators
* errors − just typed errors


`util` packages consists of many (not always related) subpackages:

* database − create & drop Postgres database.
* environment − global environment object for holding configuration and DB handles. Currently used mostly as a GORM DB handler.
* hash − different hash utils, for example base64 operations.
* log − Logrus library wrapper to be used in all parts of the app for logging.
* migrate − database migrations and seeds. Probably would be moved out of util later.
* securerandom − Golang's copy of Ruby's securerandom package. Used mostly for UUIDs generation.
* Web-related code like custom error handlers or fileserver code, probably would be moved to the top server/rest package.