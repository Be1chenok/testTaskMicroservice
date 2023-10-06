ifeq ("$(wildcard .env)",".env")
	include .env
endif

DATABASE_URL="postgres://${PG_USER}:${PG_PASSWORD}@${PG_INTERNAL_HOST}:${PG_INTERNAL_PORT}/${PG_BASE}?sslmode=${PG_SSL_MODE}"

run-containers:
	docker-compose -f docker-compose.yaml up --force-recreate

swag-init:
	swag init -g cmd/app/main.go

migrate-up:
	migrate -path ./migration -database ${DATABASE_URL} up

migrate-down:
	echo y | migrate -path ./migration -database ${DATABASE_URL} down

rm-containers:
	echo docker-compose -f docker-compose.yaml stop \
	&& echo y | docker-compose -f docker-compose.yaml rm