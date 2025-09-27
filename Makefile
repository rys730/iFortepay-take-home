migrate-up:
# 	goose -dir db/migrations postgres "postgres://user:pass@localhost:5432/app?sslmode=disable" up
	goose -dir db/migrations postgres "postgres://postgres:postgres@localhost:3033/iforte?sslmode=disable" up

run-app: ## Run main service
	go build -o ./bin/app ./cmd/
	./bin/app

gen-docs: ## Generate Swagger docs
	swag init -g cmd/app.go -o docs -d ./,./internal/handler