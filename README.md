# Development environment

## Environment variables

### DATABASE_URL

Connection string to database, e.g. 
```
$ export DATABASE_URL=postgres://postgres:abc123@0.0.0.0/sandbox
```

### JWT_SECRET

Sign key for token string generation. Any string will work for development environment, e.g.

```
$ openssl rand -base64 60
```

## Local Postgresql

Run the following command to start docker container with PostgreSQL:

```
$ podman run -p 5432:5432 --name localpg -e POSTGRESS_PASSWORD=<some password> -d postgres
```

Access to the database

```
$ podman run -it --rm postgres psql -h host.containers.internal -U postgres
```