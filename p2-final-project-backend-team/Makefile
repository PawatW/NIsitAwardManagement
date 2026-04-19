wire-gen:
	@echo "Generating code..."

	go fmt ./cmd/api/di

	wire ./...

	@echo "Code generation complete."

start-app:
	@echo "Starting the application..."

	docker-compose up --build -d

migrate-schema:
	@echo "Strating the migrate schema"

	go run ./cmd/migration/main.go
seed-data:
	@echo "Seeding database with mock data..."

	go run ./cmd/seed/main.go

	@echo "Seeding complete ✅"	