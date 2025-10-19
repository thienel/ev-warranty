
.PHONY: frontend-build
frontend-build:
	cd frontend && npm ci && npm run build

.PHONY: go-build
go-build:
	cd backend-go && go build ./...

.PHONY: go-test
go-test:
	cd backend-go && go test ./... -v

.PHONY: dotnet-restore
dotnet-restore:
	cd backend-dotnet && dotnet restore

.PHONY: dotnet-build
dotnet-build:
	cd backend-dotnet && dotnet build --no-restore --configuration Release

.PHONY: dotnet-test
dotnet-test:
	cd backend-dotnet && dotnet test --no-build --verbosity normal

.PHONY: docker-compose-down
docker-compose-down:
	docker compose down --remove-orphans || true

.PHONY: docker-compose-down-clean
docker-compose-down-clean:
	docker compose down -v --remove-orphans || true

.PHONY: docker-compose-up
docker-compose-up:
	docker compose up -d --build

.PHONY: docker-compose-update
docker-compose-update:
	docker compose pull
	docker compose up -d --build

.PHONY: ci-test
ci-test: go-test dotnet-test
	@echo "All CI tests completed"

