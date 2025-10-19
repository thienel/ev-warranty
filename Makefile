.PHONY: ci-frontend
ci-frontend:
	cd frontend && npm ci && npm run build
	@echo "Frontend CI completed"

.PHONY: ci-go
ci-go:
	cd backend-go && go build ./...
	cd backend-go && go test ./... -v
	@echo "Go CI completed"

.PHONY: ci-dotnet
ci-dotnet:
	cd backend-dotnet && dotnet restore
	cd backend-dotnet && dotnet build --no-restore --configuration Release
	cd backend-dotnet && dotnet test --no-build --verbosity normal
	@echo ".NET CI completed"