SA_PASSWORD ?= $(shell grep '^SA_PASSWORD=' .env 2>/dev/null | cut -d '=' -f2)

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

.PHONY: seed-check
seed-check:
	@echo "Checking seed data in database..."
	docker exec -it db_sqlserver_dotnet /opt/mssql-tools18/bin/sqlcmd -S localhost -U sa -P "${SA_PASSWORD}" -d WarrantyDb -Q "SELECT COUNT(*) as Customers FROM customers; SELECT COUNT(*) as VehicleModels FROM vehicle_models; SELECT COUNT(*) as Vehicles FROM vehicles;" -C

.PHONY: check-workorders
check-workorders:
	@echo "Checking work orders in database..."
	docker exec -it db_sqlserver_dotnet /opt/mssql-tools18/bin/sqlcmd -S localhost -U sa -P "${SA_PASSWORD}" -d WarrantyDb -Q "SELECT COUNT(*) as TotalWorkOrders FROM work_orders; SELECT TOP 10 * FROM work_orders ORDER BY created_at DESC;" -C

.PHONY: check-all-tables
check-all-tables:
	@echo "Checking all tables record count..."
	docker exec -it db_sqlserver_dotnet /opt/mssql-tools18/bin/sqlcmd -S localhost -U sa -P "${SA_PASSWORD}" -d WarrantyDb -Q "SELECT 'work_orders' AS TableName, COUNT(*) AS RecordCount FROM work_orders UNION ALL SELECT 'customers', COUNT(*) FROM customers UNION ALL SELECT 'vehicles', COUNT(*) FROM vehicles UNION ALL SELECT 'vehicle_models', COUNT(*) FROM vehicle_models UNION ALL SELECT 'warranty_policies', COUNT(*) FROM warranty_policies UNION ALL SELECT 'parts', COUNT(*) FROM parts UNION ALL SELECT 'part_categories', COUNT(*) FROM part_categories UNION ALL SELECT 'policy_coverage_parts', COUNT(*) FROM policy_coverage_parts;" -C

.PHONY: seed-up
seed-up:
	@echo "Running seed data..."
	@docker cp backend-dotnet/seed_data_up.sql db_sqlserver_dotnet:/tmp/seed_data_up.sql
	docker exec -i db_sqlserver_dotnet /opt/mssql-tools18/bin/sqlcmd -S localhost -U sa -P "${SA_PASSWORD}" -d WarrantyDb -C -i /tmp/seed_data_up.sql

.PHONY: seed-down
seed-down:
	@echo "Removing seed data..."
	@docker cp backend-dotnet/seed_data_down.sql db_sqlserver_dotnet:/tmp/seed_data_down.sql
	docker exec -i db_sqlserver_dotnet /opt/mssql-tools18/bin/sqlcmd -S localhost -U sa -P "${SA_PASSWORD}" -d WarrantyDb -C -i /tmp/seed_data_down.sql

.PHONY: debug-workorders
debug-workorders:
	@echo "Debugging work orders..."
	@docker cp backend-dotnet/debug_workorders.sql db_sqlserver_dotnet:/tmp/debug_workorders.sql
	docker exec -i db_sqlserver_dotnet /opt/mssql-tools18/bin/sqlcmd -S localhost -U sa -P "${SA_PASSWORD}" -d WarrantyDb -C -i /tmp/debug_workorders.sql

.PHONY: seed-reset
seed-reset: seed-down seed-up
	@echo "Seed data reset completed"

.PHONY: dotnet-logs
dotnet-logs:
	docker logs -f backend-dotnet

.PHONY: dotnet-restart
dotnet-restart:
	docker-compose restart backend-dotnet
	@echo "Backend .NET restarted"

.PHONY: dotnet-rebuild
dotnet-rebuild:
	docker-compose up -d --build backend-dotnet
	@echo "Backend .NET rebuilt and restarted"