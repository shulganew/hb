# Migrations

.PHONY: pg
pg:
	docker run --rm \
		--name=happydb_v1 \
		-v $(abspath ./docker/init/):/docker-entrypoint-initdb.d \
		-e POSTGRES_PASSWORD="postgres" \
		-d \
		-p 5439:5432 \
		postgres:15.3
	
	sleep 3
	
	goose -dir ./migrations  up

.PHONY: pg-stop
pg-stop:
	docker stop happydb_v1

.PHONY: swagger
swagger:
	go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config=internal/config/oapi.yaml --package oapi api/api.yaml
	