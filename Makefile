.PHONY: frontend-build go-build go-test dotnet-build dotnet-test docker-compose-up

frontend-build:
	cd frontend && npm ci && npm run build

go-build:
	cd auth-user-service && go build ./...

go-test:
	cd auth-user-service && go test ./... -v

dotnet-restore:
	cd vehicle-customer-service && dotnet restore

dotnet-build:
	cd vehicle-customer-service && dotnet build --no-restore --configuration Release

dotnet-test:
	cd vehicle-customer-service && dotnet test --no-build --verbosity normal

docker-compose-down:
	docker compose down --remove-orphans || true

docker-compose-up:
	docker compose up -d --build

wait-for-services:
	@echo "Waiting for services to be ready..."
	@timeout 60 bash -c 'until docker exec auth_postgres pg_isready -U auth_service -d auth_service; do sleep 2; done'
	@echo "PostgreSQL is ready"

MIGRATE = migrate -path ./auth-user-service/internal/infrastructure/database/migrations -database "postgres://auth_service:password@localhost:5432/auth_service?sslmode=disable"

db-migrate-up: wait-for-services
	$(MIGRATE) up

db-migrate-down:
	$(MIGRATE) down

db-migrate-force:
	$(MIGRATE) force $(version)

db-migrate-version:
	$(MIGRATE) version