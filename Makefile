include .env.tests
export $(shell sed 's/=.*//' .env.tests)

run:
	go run ./cmd/server/main.go

docker:
	docker-compose up

gen:
	go gen ./...

swag:
	swag init -g cmd/server/main.go

test:
	APP_LEVEL=test go test -v ./... -coverprofile=cover.out

cover:
	go tool cover -html=cover.out

mocks:
	mockgen -source=./internal/client/statistics/client.go -destination=./pkg/mocks/client/statistics/mock_client.go
	mockgen -source=./internal/repository/repository.go -destination=./pkg/mocks/repository/person/mock_person.go
