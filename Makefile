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

MIGRATE = migrate -path ./auth-user-service/internal/infrastructure/database/migrations -database "postgres://auth_service:password@postgres:5432/auth_service?sslmode=disable"

db-migrate-up:
	$(MIGRATE) up

db-migrate-down:
	$(MIGRATE) down

db-migrate-force:
	$(MIGRATE) force

db-migrate-version:
	$(MIGRATE) version
