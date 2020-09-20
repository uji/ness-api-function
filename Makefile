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

mock:
	mockgen -source ./domain/thread/usecase.go -destination ./domain/thread/usecase_mock.go -package thread

table:
	docker-compose exec aws-cli \
	aws dynamodb create-table \
	--region us-east-1 \
	--endpoint http://db-with-gui:8000 \
	--table-name Thread \
	--attribute-definitions \
		AttributeName=PK,AttributeType=S \
		AttributeName=SK,AttributeType=S \
	--key-schema AttributeName=PK,KeyType=HASH AttributeName=SK,KeyType=RANGE \
	--provisioned-throughput ReadCapacityUnits=1,WriteCapacityUnits=1
