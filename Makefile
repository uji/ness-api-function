init:
	docker volume create ness-function
	docker-compose build

clean:
	make down
	docker volume rm ness-function

up:
	docker-compose up -d

down:
	docker-compose down

build:
	GOOS=linux GOARCH=amd64 go build -o bin/main

