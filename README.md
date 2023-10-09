Authorization microservice

Other libraries: swaggo/swag, gorilla/mux, spf13/viper, redis/go-redis, golang-jwt/jwt
Databases: redis, postgresql

.env file:
SERVER_HOST=web
SERVER_PORT=8080

PG_HOST=postgres
PG_PORT=5432
PG_USER=postgres
PG_PASSWORD=postgres
PG_BASE=postgres
PG_SSL_MODE=disable

REDIS_HOST=redis
REDIS_PORT=6379
REDIS_DB=0

USER_PASSWORD_SALT=salt

TOKENS_SIGNING_KEY=secret

ACCESS_TOKEN_TTL=900
REFRESH_TOKEN_TTL=2592000



SERVER_EXTERNAL_PORT=8080

PG_EXTERNAL_HOST=0.0.0.0
PG_EXTERNAL_PORT=5433

REDIS_EXTERNAL_PORT=6378
