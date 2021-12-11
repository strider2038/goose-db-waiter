FROM golang:1-alpine AS build-env

COPY . /project

RUN set -xe \
    && cd /project \
    && go mod download \
    && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o db-waiter

FROM scratch

ENV DATABASE_URL="" \
    MIGRATION_VERSION="" \
    MIGRATIONS_TABLE=""

COPY --from=build-env /project/db-waiter /app/

ENTRYPOINT ["/app/db-waiter"]
