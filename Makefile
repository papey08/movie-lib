test:
	go mod download && docker-compose --profile test up -d && go test ./... ./...

run:
	docker-compose --profile release up
