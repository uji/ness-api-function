network:
	docker network create ness-network

init:
	docker volume create ness-api-function
	docker volume create ness-api-data
	docker-compose build

clean:
	make down
	docker volume rm ness-api-function

up:
	docker-compose up -d
	docker-compose exec app sh entrypoint.sh

down:
	docker-compose down

start:
	docker-compose start
	docker-compose exec app sh entrypoint.sh

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
