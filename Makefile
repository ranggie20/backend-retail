default: 

stack-up:
	docker-compose up -d

stack-down:
	docker-compose down --rmi local

migrate:
	go run main.go migrate

generate-queries:
	rm -rf repo/generated & sqlc generate

run-local:
	go run main.go start

run-consumer:
	go run main.go consumer

compile:
	env GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/backend main.go