include .envrc

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

.PHONY: db/migrations/new
db/migrations/new:
	@echo 'Creating migration files for ${name}'
	migrate create -seq -ext .sql -dir ./migrations ${name}

.PHONY: db/migrations/up
db/migrations/up: confirm
	@echo 'running up migrations...'
	migrate -path ./migrations -database ${DATABASE_URL} up

.PHONY: db/migrations/down
db/migrations/down: confirm
	@echo 'running down migrations...'
	migrate -path ./migrations -database ${DATABASE_URL} down

.PHONY: db/migrations/force
db/migrations/force: confirm
	@echo 'forcing migration version to ${version}...'
	migrate -path ./migrations -database ${DATABASE_URL} force ${version}

.PHONY: db/migrations/version
db/migrations/version:
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

.PHONY: test/mocks
test/mocks:
	mockery --all --with-expecter --case=snake
	@rm -rf ./pkg/mocks
	@mv ./mocks ./pkg/mocks

.PHONY: test/cover/html
test/cover/html:
	ginkgo -cover ./...
	go tool cover -html=coverprofile.out

.PHONY: test/cover/func
test/cover/func:
	ginkgo -cover ./...
	go tool cover -func=coverprofile.out
