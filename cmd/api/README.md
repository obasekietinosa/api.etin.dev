# `cmd/api`

This package wires together the HTTP API server. Key files include:

- `main.go`: bootstraps configuration, database connections, and starts the server.
- `routes.go`: declares the routing table and middleware stack.
- `handlers.go` and `handler_*.go`: provide request handlers for domain entities such as companies and roles.
- `helper.go`: shared utilities for responding to requests and managing dependencies.
- `handler_swagger.go`: serves the generated OpenAPI specification at `/swagger`.

The command generates the OpenAPI document on startup so the latest routes and schemas are always available from the `/swagger` endpoint.

The command expects the database DSN and admin credentials to be supplied either as flags or via the environment variables `WEBSITE_DB_DSN`, `WEBSITE_ADMIN_EMAIL`, and `WEBSITE_ADMIN_PASSWORD`.

You can optionally provide `WEBSITE_DEPLOY_WEBHOOK_URL` (or the `-deploy-webhook-url` flag) to ping an external deployment
service whenever write operations succeed. When configured, the server asynchronously issues a `POST` request containing an
empty JSON object to the supplied URL.

## Asset uploads

Authenticated administrators can push files to Cloudinary through the `/v1/assets` endpoint. Send a
`multipart/form-data` request containing a single `file` part; payloads are limited to 10 MiB. The
API streams the file to Cloudinary, persists the returned metadata, and responds with the database
record containing the secure delivery URL.

```bash
curl -X POST http://localhost:4000/v1/assets \
  -H "Authorization: Bearer ${TOKEN}" \
  -F "file=@./banner.png"
```

Requests missing the `file` part or exceeding the size limit are rejected with `400` or `413`
responses. Upload failures from Cloudinary surface as `502`, and database persistence issues return
`500`.
