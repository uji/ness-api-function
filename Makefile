SHELL=bash

HAS_DOCKER := $(shell\
	if type "docker" > /dev/null 2>&1; then\
			echo true;\
		else\
			echo false;\
	fi\
)

# commands for host with docker command
ifeq ($(HAS_DOCKER),true)
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

serve-with-auth:
	docker-compose exec api go run ./testsrv -teamID $(TEAM_ID) -userID $(USER_ID)

table:
	docker-compose exec api go run ./tools/dbtool/ create
	docker-compose exec api go run ./tools/elsch/ create

destroy-table:
	docker-compose exec api go run ./tools/dbtool/ destroy
	docker-compose exec api go run ./tools/elsch/ delete

status:
	@echo "--elasticsearch--"
	curl -X GET "localhost:9200/_cat/health?v&pretty"
	@echo "--elasticsearch indices--"
	@curl -X GET "localhost:9200/_cat/indices?v&pretty"
endif

# commands for container shell or host without docker command
ifeq ($(HAS_DOCKER),false)
serve:
	go run ./testsrv

serve-with-auth:
	go run ./testsrv -teamID $(TEAM_ID) -userID $(USER_ID)

mock:
	mockgen -source ./domain/thread/usecase.go -destination ./domain/thread/usecase_mock.go -package thread

table:
	go run ./tools/dbtool/ create
	go run ./tools/elsch/ create

destroy-table:
	go run ./tools/dbtool/ destroy
	go run ./tools/elsch/ delete

status:
	@echo "--elasticsearch--"
	@curl -X GET "http://elasticsearch:9200/_cat/health?v&pretty"
	@echo "--elasticsearch indices--"
	@curl -X GET "http://elasticsearch:9200/_cat/indices?v&pretty"

es-threads:
	@curl -X GET "elasticsearch:9200/thread/_search?pretty" -H 'Content-Type: application/json' -d' { "query": {"match_all": {}} }'
endif
