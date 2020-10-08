SHELL=bash

IN_DOCKER := $(shell\
	if type "docker" > /dev/null 2>&1; then\
			echo false;\
		else\
			echo true;\
	fi\
)

# commands for host machine shell
ifeq ($(IN_DOCKER),false)
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

serve:
	docker-compose exec api go run ./testsrv

table:
	docker-compose exec api go run ./tools/dbtool/ create

destroy-table:
	docker-compose exec api go run ./tools/dbtool/ destroy
endif

# commands for container shell
ifeq ($(IN_DOCKER),true)
serve:
	go run ./testsrv

mock:
	mockgen -source ./domain/thread/usecase.go -destination ./domain/thread/usecase_mock.go -package thread

table:
	go run ./tools/dbtool/ create

destroy-table:
	go run ./tools/dbtool/ destroy
endif
