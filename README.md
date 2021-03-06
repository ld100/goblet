# Golang microservice app prototype

Inspired by [Clean Architecture](https://hackernoon.com/golang-clean-archithecture-efd6d7c43047).

Application consists of configuration and Golang code. Configuration is concentrated in the root folder and `config`.
Go code consists of common use packages and domain-specific ones. Common use packages are:

* cmd − command-line tools
* pkg/server − entry point for server-related code. Consists of REST and gRPC servers and its tools (utils, error handlers).
* pkg/persistence − database & ORM handlers, database setup, migration & seed operations.
* pkg/util − project-wide utility code, including database, cryptography, data migration, environment, etc

`pkg/domain` includes domain-specific packages, e.g. users, orders, articles. Each domain package may consist of the following parts:

* model − models structs with JSON/YAML/XML serialization instructions. Usually model is a GORM db model with validations.
* service − Service layer AKA Commands AKA Logic
* rest − REST-specific routes controllers, middlewares. Uses CHI framework.
* grpc − gRPC-specific controllers and middlewares.
* form − forms and form validators
* error − just typed errors


`pkg/util` packages consists of many (not always related) subpackages:

* hash − different hash utils, for example base64 operations.
* logger − Logrus library wrapper to be used in all parts of the app for logging.
* securerandom − Golang's copy of Ruby's securerandom package. Used mostly for UUIDs generation.

## Monitoring

Monitoring to be done via Prometheus/Grafana

## Distributed logging

Logging is cloud-adapted and using ELK stack.
https://github.com/deviantony/docker-elk/blob/master/docker-compose.yml
