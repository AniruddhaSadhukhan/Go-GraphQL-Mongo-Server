# Go-GraphQL-Mongo-Server

A sample server in Go language that exposes GraphQL APIs and connects to a MongoDB with features like GraphiQl IDE, OIDC/OAUTH & Personal Access Token Authentication, Rate Limiting, Prometheus Metrics, Health Check, Dynamic Configuration, Logging, Database Migration, Cron Jobs, Negroni Middlewares, HTTPS Support, Telemetry, GraphQl Sanitizer, Cache etc.

## Major Features

### Mongo DB

This server uses MongoDB as the database and contains GoLang code for establishing connection & various CRUD operations.
The file [models/db.go](./models/db.go) contains the database connection code and various utility function for database operations.

### GraphQL

The server exposes GraphQL APIs that can be queried. The file [gqlhandler/graphqlHandler.go](./gqlhandler/graphqlHandler.go) contains the GraphQL handlers.

The schema is defined in [gqlhandler/schema](./gqlhandler/schema/) folder. The mutations and queries with their resolvers are defined in [gqlhandler/mutation](./gqlhandler/mutation/) and [gqlhandler/query](./gqlhandler/query/) respectively.

### GraphiQl

The server exposes a GraphiQl webapp that is a graphical interactive in-browser GraphQL IDE with documentation of various queries and mutations. It has a very easy to use plugin (Explorer Plugin) that helps in creating different GraphQl queries and mutations with just mouse clicks. It has custom Header support, history etc. More details can be found in [GraphiQl GitHub Page](https://github.com/graphql/graphiql#graphiql).

Here the file [gqlhandler/graphiqlHandler.go](./gqlhandler/graphiqlHandler.go) contains the GraphiQl handlers and the file [resources/graphiql.html](./resources/graphiql.html) contains the GraphiQl webapp.

### OIDC/OAUTH Authentication Middleware

The server has support for OpenID Connect (OIDC) for authentication. The file [auth/middleware.go](./auth/middleware.go) contains the OIDC authentication middleware. To setup the configuration for OIDC, please set the required environment variables like `OIDC_URL`, `CLIENT_ID`, `CLIENT_SECRET`.

This does local Access Token verification ie. it gets the RSA Public Key from the OIDC JWKS URL and verifies the signature of the Access Token against it. Then it checks the validity of the Access Token along with some of the claims. You can read more about it [here](https://developer.okta.com/docs/guides/validate-id-tokens/main/#what-to-check-when-validating-an-id-token).

### Personal Access Token (PAT) Authentication Middleware

This server has support for Personal Access Token (PAT) authentication. The file [auth/middleware.go](./auth/middleware.go) contains the PAT authentication middleware. To setup the configuration for PAT, please set the required environment variables like `JWT_PRIVATE_KEY`.

This is a in-house PAT implementation. This creates a JWT PAT with custom expiry date, signs it with the RSA private key and returns it to the user. It only stores a hash (SHA256 checksum) of the generated token in the database. For validation, it compares the hash with the one stored in the database and verifies the signature of the JWT PAT. The code related to these can be found in [gqlhandler/mutation/tokenMt.go](./gqlhandler/mutation/tokenMt.go) and [gqlhandler/query/tokenQl.go](./gqlhandler/query/tokenQl.go).

### API Rate Limiting

The server has support for API rate limiting. The file [routes/limiter.go](./routes/limiter.go) contains the API rate limiting middleware. To set this up, provide the environment variable `API_LIMIT_PER_SECOND`, whose default value is 500, meaning it will allow 500 requests per second from a particular IP address.

### Prometheus Metrics

This server exposes a `/metrics' endpoint that gives the list of Prometheus metrics for this server instance. You can find the code in [routes/routes.go](./routes/routes.go)

### Health Check

This server exposes a `/health' endpoint that checks the health of the server. It will also regularly check the database connection and if it fails for some authentication related reason, it will fail the health check. You can find the code in [routes/healthchecks.go](./routes/healthchecks.go)

### Dynamic Configuration

This server reads various configuration from the environment variables. You can find the code in [config/config.go](./config/config.go)

### Logger

This server uses [zap logger](https://github.com/uber-go/zap) for logging. You can find the code in [logger/zapLogger.go](./logger/zapLogger.go).

It logs the time, level, caller (file, function, line), message, and stacktrace (for errors or warnings) of the message.

In case of local development, it will output the log in Console Format with color-coded output.
In case of production, it will output the log in JSON Format.

### Database Migration

This server uses [golang-migrate/migrate](https://github.com/golang-migrate/migrate) for database migration. You can find the code in [dbmigration/schema_migration.go](./dbmigration/schema_migration.go) and the migration scripts in [resources/schema_migrations](./resources/schema_migrations/).

This is useful to create indexes, changes some schemas in the DB etc and the DB will also remember the version of the migrations.

### Cron Jobs

The server has support for cron jobs. The file [main.go](./main.go) contains the cron jobs. Currently the pinging of the DB every 5 minutes is implemented as a cron job.

### Negroni Middleware

The server has [Negroni](https://github.com/urfave/negroni) middleware added. You can find the code in [main.go](./main.go).
Currently the following are used:

- **Panic Recovery Middleware** : This middleware catches panics and responds with a 500 response code
- **Request/Response Logger Middleware** : This middleware logs each incoming request and response with the response code, time taken, etc.

### HTTPS Support

The server has support for HTTPS. You can find the code in [main.go](./main.go). To set this up for HTTPS, please set the required environment variables like `HTTPS_ENABLED`, `HTTPS_CERT_FILE_PATH`, `HTTPS_KEY_FILE_PATH`.

### Telemetry

The server has support for sending telemetry for all the GraphQl calls. It can log the GraphQl calls, errors, user who called it, the device information where the server is running etc . You can find the code in [telemetry/telemetry.go](./telemetry/telemetry.go).

### GraphQl Sanitizer

The server has support for sanitizing the user inputs in GraphQl queries & mutations. You can find the code in [sanitizer/sanitizer.go](./common/sanitize.go). It recursively sanitizes the user inputs. In case of string inputs, it uses [bluemonday](https://github.com/microcosm-cc/bluemonday) to sanitize the input.

### GraphQl Client

It has some utility functions for making GraphQL requests. You can find the code in [common/httpUtils.go](./common/httpUtils.go).

### Cache Manager

This has a cache package that helps in caching any type of data. You can find the code in [cache/cache.go](./cache/cache.go). It supports custom TTL for each cache and user implementable cache update function.

---

## Update dependency

`go get -u ./...`

## Install Linter

`go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest`

## Run Lint

`golangci-lint run`
