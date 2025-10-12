# `cmd/api`

This package wires together the HTTP API server. Key files include:

- `main.go`: bootstraps configuration, database connections, and starts the server.
- `routes.go`: declares the routing table and middleware stack.
- `handlers.go` and `handler_*.go`: provide request handlers for domain entities such as companies and roles.
- `helper.go`: shared utilities for responding to requests and managing dependencies.
- `handler_swagger.go`: serves the generated OpenAPI specification at `/swagger`.

The command generates the OpenAPI document on startup so the latest routes and schemas are always available from the `/swagger` endpoint.

The command expects the database DSN and authorisation token to be supplied either as flags or via the environment variables `WEBSITE_DB_DSN` and `WEBSITE_AUTH_KEY`.
