# api.etin.dev
## Overview
This is the backend service for my new personal website.
Overkill? Certainly. I'm doing this as an excuse to learn Go and to experiment with a variety of things, including data modelling, content management systems and CI/CD processes.

As the primary purpose of this is experimentation (and maybe eventually writing some blog posts) a lot of the tooling and basics will be handrolled rather than reaching for off the shelf components. I don't claim to be able to build a better ORM than anyone else, but I will hopefully enjoy it.

### Stack
**Server** - Golang server

**Database** - Postgres

**CI/CD** - GitHub Actions

### Roadmap
I'm hoping that this backend can power everything I need on my personal site
including blog posts, work history and a projects showcase.

## Contributing
Clone the repository and navigate to the folder.

### Database setup
To start the server locally, you will need a Postgres database.
If you already have one set up, you can create a database for this project and set it up following the scripts 
specified [in the internal data folder](internal/data/README.md).

**TODO**: Add seeders so that database can be easily spun up.

Alternatively, you can spin up a docker container with Postgres.

```bash
docker run --name my_postgres \        
  -e POSTGRES_USER=user \
  -e POSTGRES_PASSWORD=password \
  -e POSTGRES_DB=website \
  -p 5432:5432 \
  -v pgdata:/var/lib/postgresql/data \
  -d postgres:17
```

You can connect to this database and run the SQL to seed the database

```bash
psql -h localhost -U user -d website
```

### Running the server
To run the server, you will need to build or run the `cmd/api` package.
It expects environment variables (these can also be passed in as flags) for the database connection string
and for the admin login credentials used to access write operations in the API.

```bash
export WEBSITE_DB_DSN='postgres://website:etin@localhost/website?sslmode=disable' && \
export WEBSITE_ADMIN_EMAIL='admin@example.com' && \
export WEBSITE_ADMIN_PASSWORD='super-secret-password' && \
export WEBSITE_CORS_TRUSTED_ORIGINS='http://localhost:3000 https://admin.example.com' && \
export WEBSITE_CLOUDINARY_CLOUD_NAME='your-cloud-name' && \
export WEBSITE_CLOUDINARY_API_KEY='your-api-key' && \
export WEBSITE_CLOUDINARY_API_SECRET='your-api-secret' && \
export WEBSITE_CLOUDINARY_FOLDER='website-assets' && \
go run ./cmd/api
```

The optional `WEBSITE_CORS_TRUSTED_ORIGINS` environment variable (or the `-cors-trusted-origins` flag) accepts a space separated
list of origins that should receive CORS headers. When unset, the API will automatically trust `https://etin.dev` and
`https://admin.etin.dev`.

With the server running you can exchange the admin credentials for a bearer token by
posting to the login endpoint:

```bash
curl -X POST http://localhost:4000/v1/admin/login \
  -H 'Content-Type: application/json' \
  -d '{"email":"admin@example.com","password":"super-secret-password"}'
```

The JSON response contains the bearer token and expiry timestamp. Supply this token in the
`Authorization` header when calling any create, update or delete endpoints. To invalidate the
token, call `POST /v1/admin/logout` with the same header.

### Asset storage pipeline

Uploads are handled through Cloudinary. Set the following environment variables (or equivalent
flags) so the API can initialise the uploader:

* `WEBSITE_CLOUDINARY_CLOUD_NAME`
* `WEBSITE_CLOUDINARY_API_KEY`
* `WEBSITE_CLOUDINARY_API_SECRET`
* `WEBSITE_CLOUDINARY_FOLDER` *(optional)*

`internal/assets` wraps the Cloudinary SDK and exposes an `Uploader` interface so other packages
can remain storage agnostic. Metadata about each upload is persisted to the `assets` table (see
`internal/data/assets.go`) which unlocks future join tables for connecting images to notes,
projects or roles.

When creating or updating a project (and future models that support rich images), the frontend
should request an upload first, then include the resulting Cloudinary URL in the payload's optional
`imageUrl` field.

### Tests
Not a lot of tests yet, but hopefully changing soon. There are tests written for the `querybuilder` package. To run
these, run:

```bash
go test -timeout 30s api.etin.dev/pkg/querybuilder
```

### OpenAPI specification
The `/swagger` endpoint serves a generated OpenAPI document embedded in the binary. To regenerate the specification after
changing handlers or schemas, run:

```bash
go run ./cmd/openapi
```

This command rewrites `cmd/api/openapi.json`. Commit the updated file to keep the embedded schema in sync. Pull requests
automatically verify that the committed document matches the generated output, and you can comment “generate openapi spec” on
an open PR to trigger CI to regenerate and push the latest document for you.
