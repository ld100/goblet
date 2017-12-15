# Goblet Roadmap

Prioritized project roadmap.

## Must have

* Authorization for user update/delete;
* Admin role, administrators could update/delete users too;
* Automated testing: integration tests and unit-tests for models;
* Implement [gRPC](https://github.com/grpc/grpc-go) interface for users and sessions. [Middlewares](https://github.com/grpc-ecosystem/go-grpc-middleware) might help;
* Hook-up prometheus monitoring (log the start and end of each request with the elapsed processing time);

## Should have

* Validate HTTP IDs to be integers in requests;
* Update error messages and status code to be 100% appropriate;
* Return user array of errors, not just first error;
* Command-line interface:
  * Start/Stop server on specified ports
  * Run migrations
  * Generate secret key
* Implement REST timeouts & throttling;
* Enable CORS;
* Login security:
  * Limit auth login attempts in time;
  * Overall login rate limit per IP in time;
  * Log all login attempts
* [Audit Trail](https://en.wikipedia.org/wiki/Audit_trail);

## Would have

* Add separate validation method to check that all environment variables are present and valid;
* Multiple e-mails per user, e-mail confirmation;
* Password resets;
* Hook-up [Swagger](https://swagger.io/) for docs;
* Split migrations into separate files per each migration;
* Implement graceful server shutdown;
* Clean-up expired sessions;
* Multiple phone numbers per user, phone number confirmation;
* Abbility to block specific tokens (to log user out instantly);

## Won't have

* Implement admin-only routes;
* Add pagination to users list;
* Build 2-factor auth;
* Hook-up [hydra](https://github.com/ory/hydra) or [dex](https://github.com/coreos/dex) for authentication;
* Gzip compression for clients that accept compressed responses;