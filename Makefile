network:
	docker network create ness-network

init:
	docker volume create ness-api-function
	docker volume create ness-api-data
	docker-compose build

clean:
	make down
	docker volume rm ness-api-function
	docker volume rm ness-api-data

up:
	docker-compose up -d
	docker-compose exec api sh entrypoint.sh

down:
	docker-compose down

start:
	docker-compose start
	docker-compose exec api sh entrypoint.sh

stop:
	docker-compose stop

build:
	GOOS=linux GOARCH=amd64 go build -o bin/main

migrate:
	sql-migrate up

mock:
	mockgen -source ./domain/thread/usecase.go -destination ./domain/thread/usecase_mock.go -package thread

xo:
	rm xogen/*
	xo pgsql://$(DB_USER):$(DB_PASS)@$(DB_HOST)/$(DB_NAME)?sslmode=disable -o xogen

table:
	docker-compose exec aws-cli \
	aws dynamodb create-table \
	--region us-east-1 \
	--endpoint http://db-with-gui:8000 \
	--table-name Thread \
	--attribute-definitions \
		AttributeName=Team,AttributeType=S \
		AttributeName=Title,AttributeType=S \
	--key-schema AttributeName=Team,KeyType=HASH AttributeName=Title,KeyType=RANGE \
	--provisioned-throughput ReadCapacityUnits=1,WriteCapacityUnits=1
