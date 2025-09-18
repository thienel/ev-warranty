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

install-migrate:
	@which migrate > /dev/null 2>&1 || { \
		echo "Installing golang-migrate..."; \
		go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest; \
	}

wait-for-services:
	@echo "Waiting for services to be ready..."
	@echo "Waiting for PostgreSQL container to start..."
	@timeout 120 bash -c 'until docker compose ps postgres | grep -q "healthy"; do echo "Waiting for postgres health check..."; sleep 3; done' || { \
		echo "PostgreSQL failed to become healthy within 120 seconds"; \
		docker compose logs postgres; \
		exit 1; \
	}
	@echo "Waiting for PostgreSQL connection..."
	@timeout 60 bash -c 'until docker exec auth_postgres pg_isready -U auth_service -d auth_service; do echo "Waiting for pg_isready..."; sleep 2; done' || { \
		echo "PostgreSQL connection failed"; \
		docker compose logs postgres; \
		exit 1; \
	}
	@echo "PostgreSQL is ready"
	@echo "Waiting for auth service to be ready..."
	@timeout 90 bash -c 'until docker compose ps auth-user-service | grep -q "healthy"; do echo "Waiting for auth service health check..."; sleep 3; done' || { \
		echo "Auth service failed to become healthy within 90 seconds"; \
		docker compose logs auth-user-service; \
		exit 1; \
	}
	@echo "All services are ready"

MIGRATE = migrate -path ./auth-user-service/internal/infrastructure/database/migrations -database "postgres://auth_service:password@localhost:5432/auth_service?sslmode=disable"

db-migrate-up: install-migrate wait-for-services
	@echo "Running database migrations..."
	$(MIGRATE) up || { \
		echo "Migration failed, checking migration status..."; \
		$(MIGRATE) version; \
		exit 1; \
	}
	@echo "Migrations completed successfully"

db-migrate-down: install-migrate
	@echo "Rolling back database migrations..."
	$(MIGRATE) down

db-migrate-force: install-migrate
	@echo "Forcing migration version to $(version)..."
	$(MIGRATE) force $(version)

db-migrate-version: install-migrate
	$(MIGRATE) version


ci-test: go-test dotnet-test
	@echo "All CI tests completed"
