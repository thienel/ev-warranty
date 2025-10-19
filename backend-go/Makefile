include .envrc

.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

.PHONY: confirm
confirm:
	@printf 'Are you sure? [y/N]' && read ans && [ $${ans:-N} = y ]

.PHONY: run
run:
	go run ./cmd/server/

.PHONY: check
check:
	go vet ./...
	go build -o /dev/null ./cmd/server/

.PHONY: db/psql
db/psql:
	psql ${DATABASE_URL}

.PHONY: install-migrate
install-migrate:
	@which migrate > /dev/null 2>&1 || { \
		echo "Installing golang-migrate..."; \
		go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest; \
	}

.PHONY: db/migrations/new
db/migrations/new: install-migrate
	@echo 'Creating migration files for ${name}'
	migrate create -seq -ext .sql -dir ./migrations ${name}

.PHONY: db/migrations/up
db/migrations/up: confirm install-migrate
	@echo 'running up migrations...'
	migrate -path ./migrations -database ${DATABASE_URL} up

.PHONY: db/migrations/down
db/migrations/down: confirm install-migrate
	@echo 'running down migrations...'
	migrate -path ./migrations -database ${DATABASE_URL} down

.PHONY: db/migrations/force
db/migrations/force: confirm install-migrate
	@echo 'forcing migration version to ${version}...'
	migrate -path ./migrations -database ${DATABASE_URL} force ${version}

.PHONY: db/migrations/version
db/migrations/version: install-migrate
	migrate -path ./migrations -database ${DATABASE_URL} version

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

.PHONY: mocks
mocks:
	mockery --all --with-expecter --case=snake
	@rm -rf ./pkg/mocks
	@mv ./mocks ./pkg/mocks