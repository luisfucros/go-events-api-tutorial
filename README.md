# go-events-api-tutorial

## Docker & Migrations

- The Docker image runs an entrypoint script that waits for the database, runs migrations using `golang-migrate`, then starts the app.
- The migration tool is pinned in the `Dockerfile` (`MIGRATE_VERSION=v4.16.0`) to avoid surprises from `latest` releases.

### Local development

- Build and run the stack: `docker-compose up --build`
- The API container will wait for the database and run `migrate up` automatically on start.

### Notes

- For local dev we still use the root password from `.env` (do not commit secrets). For production, prefer secrets managers or Docker secrets.
- To run migrations manually: `docker-compose run --rm events-api /usr/local/bin/migrate -path=/migrations -database "mysql://root:events_password@tcp(db:3306)/events" up`