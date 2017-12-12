# Heroku Go Database Example

This is a simple example of a PostgreSQL-backed Go API application running on Heroku

## Components

```
heroku-go-db-example         -- Directory / repository name
├── .env.dev                 -- Environment variables for local development
├── Makefile                 -- Targets for local development (optional)
├── Procfile                 -- Defines a single `web` process
├── README.md                -- This file
├── bin
│   └── heroku-go-db-example -- The locally-built binary (gitignored)
├── main.go                  -- The main Go application
└── vendor
    ├── github.com
    │   └── lib
    │       └── pq           -- Vendored dependency [1]
    └── vendor.json          -- Buildpack identification / Heroku configuration [2]
```

## Running Locally

When the application is deployed, the [Heroku Go buildpack](https://github.com/heroku/heroku-buildpack-go) will run `go install .` to build and install your application to `/app/bin`.  We want to mimic this in local development, so the `run` target in the `Makefile` will:

1. Copy `.env.dev` to `.env` to make environment variables available to the application
1. Create the local database
1. Compile `main.go` to `bin/heroku-go-db-example` -- this is important because Go will use the pacakge name (`heroku-go-db-example`) when creating the binary.
1. Prepend the `./bin` directory to the `$PATH` and run `heroku local`.  This allows the `web` process in the `Procfile` to find the `heroku-go-db-example` binary.

This is all done for you by running `make` since `run` is the default target.

## Deploying

```
$ heroku apps:create heroku-go-db-example
$ heroku addons:create heroku-postgresql
$ git push heroku master
```

## Testing

```
GET:  curl -s https://heroku-go-db-example.herokuapp.com | jq .
POST: curl -s -d '{"username":"reagent"}' https://heroku-go-db-example.herokuapp.com | jq .
404:  curl -i -X OPTIONS https://heroku-go-db-example.herokuapp.com
404:  curl -i https://heroku-go-db-example.herokuapp.com/unknown
```

## Links

* [1] https://github.com/lib/pq
* [2] https://github.com/kardianos/govendor
