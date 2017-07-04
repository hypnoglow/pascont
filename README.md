# pascont

[![Build Status](https://travis-ci.org/hypnoglow/pascont.svg?branch=master)](https://travis-ci.org/hypnoglow/pascont)
[![Coverage Status](https://coveralls.io/repos/github/hypnoglow/pascont/badge.svg?branch=master)](https://coveralls.io/github/hypnoglow/pascont?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/hypnoglow/pascont)](https://goreportcard.com/report/github.com/hypnoglow/pascont)
[![codebeat badge](https://codebeat.co/badges/8cf5adbc-fc6b-40ce-8534-a4270a25ee60)](https://codebeat.co/projects/github-com-hypnoglow-pascont-master)


Simple passport for an application or a service. It handles accounting and sessioning.

## Running the server

First of all a config file needs to be prepared.
Copy `resources/config/config.json.dist` to any convenient place, rename and fill it in as you need.

### Docker-compose

    docker-compose up --build

### Manual

For simplicity let's run PostgreSQL in Docker.
The main purpose is to show how to manually run the `pascont` server.

    docker run --rm --name pascont-postgres \
        -p 5432:5432 \
        -e POSTGRES_PASSWORD=123123 \
        -v ~/data/pascont/postgres:/var/lib/postgresql/data \
        postgres:9.6-alpine
        
The database MUST exist when launching the `pascont` server. So if you don't have one,
run a migration:

    psql -h 127.0.0.1 -p 5432 -U postgres -f resources/migrations/init.sql

Export the path to your config file:

    export PASCONT_CONFIG_PATH=/path/to/config_local.json
    
Launch the server:

    go build -o ./pascont && ./pascont

## Server requests examples

Add an account:

    curl -i -X POST \
      http://localhost:9090/accounts \
      -H 'content-type: application/json' \
      -d '{
    	"name": "email@email.com",
    	"password": "password"
      }'
    
Log in:

    curl -i -X POST \
      http://localhost:9090/sessions \
      -H 'content-type: application/json' \
      -d '{
    	"name": "email@email.com",
    	"password": "password"
      }'
    
Get session by token:

    curl -i -X GET \
      http://localhost:9090/sessions/3900cdf3-168b-43be-9034-48bf44bb7c9c \
      -H 'authorization: Bearer MzkwMGNkZjMtMTY4Yi00M2JlLTkwMzQtNDhiZjQ0YmI3YzljMTQ5ODk0MjIyOd4dbUR7plOWVIeYenr9OE_AW36838yOBCrZLsdGNf7apeXGtrube15eHbOv__eMQYFJeW_1YtjA1du6AScA-C4=' \
      -H 'content-type: application/json'
      
Extend session:

    curl -i -X PATCH \
      http://localhost:9090/sessions/3900cdf3-168b-43be-9034-48bf44bb7c9c \
      -H 'authorization: Bearer MzkwMGNkZjMtMTY4Yi00M2JlLTkwMzQtNDhiZjQ0YmI3YzljMTQ5ODk0MjIyOd4dbUR7plOWVIeYenr9OE_AW36838yOBCrZLsdGNf7apeXGtrube15eHbOv__eMQYFJeW_1YtjA1du6AScA-C4=' \
      -H 'content-type: application/json' \
      -d '{
    	"expiresAt": "2017-07-01T21:10:29Z"
      }'

## API communication schema

Resources that respond with a single object, e.g. 
`GET /resources/:id`, `POST /resources`, `PUT /resources/:id`:

    {
        "result": {
            // json object
        },
        // optional meta:
        "meta": {
            // metadata
        }
    }

Resources that respond with multiple objects, e.g. `GET /resources`:

    {
        "results": [
            {
                // json object
            },
            {
                // json object
            }
        ],
        // optional meta:
        "meta": {
            // metadata
        }
    }

Body can also be empty on some requests, e.g `DELETE /resources/:id`.

Any request can result in error. If so, status codes will be in range 
`400-499` (client errors) or `500-599` (server errors).
On error response, body can be either empty or in a form of:

    {
        "errors": [
            {
                "message": "Validation error: ...",
                // optional field:
                "field": "name",
                // optional value:
                "value": "email@email.com"
            }
        ]
    }

## TODO

In order to approach a stable release, this needs to be implemented at least: 

Features:
- [ ] Action: Delete session by ID
- [ ] Action: Delete all account sessions
- [ ] Action: Change account password
- [ ] JWTs
- [ ] Think about scopes?

Technical:
- [ ] Implement Redis session storage
- [ ] Improved Logging (HTTP + errors)
- [ ] Auto documentation
- [ ] Wrap all low-level errors
- [ ] Benchmark to find leaks (potential not closed things)

## License

[Apache 2.0](https://github.com/hypnoglow/pascont/blob/master/LICENSE)