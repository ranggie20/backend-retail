# Backend for Cerita Kaos

## Pre-requisites
1. Install [sqlc](https://github.com/sqlc-dev/sqlc/blob/main/docs/overview/install.md)

    Quick install command:
    ```
    go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
    ```

## Development

1. Start database server: `docker-compose up -d`
2. Copy `config.example.yaml` to `config.yaml`
3. Generate sqlc: `make generate-queries`
4. Update database: `go run main.go migrate`
5. Start server: `go run main.go start`

Whenever there is changes in queries, re-run step number 3.

Whenever there is changes in db_schema, re-run step number 4.
