# Include variables from the .envrc file
include .envrc

# ==================================================================================== #
# HELPERS
# ==================================================================================== #

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

.PHONY: confirm
confirm:
	@printf 'Are you sure? [y/N]' && read ans && [ $${ans:-N} = y ]

# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #

## run/api: run the cmd/api application
.PHONY: run/server
run/server:
	go run ./cmd/server/main.go

## db/psql: connect to the database using psql
.PHONY: db/psql
db/psql:
	psql ${DATABASE_URL}

# Install migrate tool if not present (for CI environments)
.PHONY: install-migrate
install-migrate:
	@which migrate > /dev/null 2>&1 || { \
		echo "Installing golang-migrate..."; \
		go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest; \
	}

## db/migrations/new name=$1: create a new database migration
.PHONY: db/migrations/new
db/migrations/new: install-migrate
	@echo 'Creating migration files for ${name}'
	migrate create -seq -ext .sql -dir ./internal/infrastructure/database/migrations ${name}

## db/migrations/up: apply all up database migrations
.PHONY: db/migrations/up
db/migrations/up: confirm install-migrate
	@echo 'running up migrations...'
	migrate -path ./internal/infrastructure/database/migrations -database ${DATABASE_URL} up

## db/migrations/up-ci: apply all up database migrations without confirmation (for CI)
.PHONY: db/migrations/up-ci
db/migrations/up-ci: install-migrate
	@echo 'running up migrations in CI mode...'
	migrate -path ./internal/infrastructure/database/migrations -database ${DATABASE_URL} up || { \
		echo "Migration failed, checking current version..."; \
		migrate -path ./internal/infrastructure/database/migrations -database ${DATABASE_URL} version; \
		exit 1; \
	}

## db/migrations/down: apply all down database migrations
.PHONY: db/migrations/down
db/migrations/down: confirm install-migrate
	@echo 'running down migrations...'
	migrate -path ./internal/infrastructure/database/migrations -database ${DATABASE_URL} down

## db/migrations/force version=$1: set the migration version without running migrations
.PHONY: db/migrations/force
db/migrations/force: confirm install-migrate
	@echo 'forcing migration version to ${version}...'
	migrate -path ./internal/infrastructure/database/migrations -database ${DATABASE_URL} force ${version}

## db/migrations/version: show current migration version
.PHONY: db/migrations/version
db/migrations/version: install-migrate
	migrate -path ./internal/infrastructure/database/migrations -database ${DATABASE_URL} version

# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #

## audit: tidy dependencies and format, vet and test all code
.PHONY: audit
audit:
	@echo 'Tidying and verifying  module dependencies...'
	go mod tidy
	go mod verify
	@echo 'Formatting code...'
	go fmt ./...
	@echo 'Vetting code...'
	staticcheck ./...
	@echo 'Running tests...'
	go test -race -vet=off ./...
