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
It expects two environment variables (these can also be passed in as flags) for the database connection string and for 
the authorisation key (this is a bit of a hack to secure the create and edit endpoints while I flesh out a 
proper auth system)

```bash
xport WEBSITE_DB_DSN='postgres://website:etin@localhost/website?sslmode=disable' && \
export WEBSITE_AUTH_KEY='token-123' && \
go run ./cmd/api
```

### Tests
Not a lot of tests yet, but hopefully changing soon. There are tests written for the `querybuilder` package. To run 
these, run:

```bash
go test -timeout 30s api.etin.dev/pkg/querybuilder
```