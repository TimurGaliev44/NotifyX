deps:
	@go mod tidy
run-check:
	@go run cmd/check/main.go

make-up:
	@docker-compose up -d --build
make-down:
	@docker-compose down
make-check:
	@docker-compose up --build check