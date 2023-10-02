CREATE TABLE IF NOT EXISTS users
(
    "id" SERIAL NOT NULL UNIQUE,
    "email" VARCHAR(64) NOT NULl UNIQUE,
    "username" VARCHAR(64) NOT NULL UNIQUE,
    "password_hash" VARCHAR(256) NOT NULL,
    "registered_at" TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS tokens
(
    "user_id" INT NOT NULL,
    "access_token" VARCHAR(256) NOT NULL,
    "refresh_token" VARCHAR(256) NOT NULL
);