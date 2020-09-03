init:
	docker volume create ness-api-function
	docker volume create ness-api-data
	docker network create ness-api-network
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
