run:
	docker-compose -f docker-compose.yaml up --force-recreate

swag:
	swag init -g cmd/app/main.go

migup:
	migrate -path ./migration -database 'postgres://postgres:postgres@0.0.0.0:5433/postgres?sslmode=disable' up

migdown:
	echo y | migrate -path ./migration -database 'postgres://postgres:postgres@0.0.0.0:5433/postgres?sslmode=disable' down

rm:
	docker-compose -f docker-compose.yaml stop \
	&& docker-compose -f docker-compose.yaml rm