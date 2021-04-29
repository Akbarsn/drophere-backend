# Drophere

Drophere is a platform for submitting your task quickly.

## How to run locally

1. Copy the env file `.env.example` to `.env`.
2. Execute db.sql from this [file][migration-file].
3. Run `go run server/main.go` to start the app.
4. Browse to `localhost:8080` by your own browser.

## How to run using docker-compose

1. Run this command

   ```bash
    docker-compose up
   ```

## Contributing to this project

Interested in contributing? please check out [the Contributing Guide](CONTRIBUTING.MD) to get started

[migration-file]: files/sql/migrate.sql
