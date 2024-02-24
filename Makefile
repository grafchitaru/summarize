test:
	go test -v ./...

up:
	docker compose up -d

run:
	go run cmd/summarize/main.go

migrate:
	goose -dir migrations postgres "postgres://root:root@localhost:54322/app" up