CREATE TABLE IF NOT EXISTS users
(
    "id" SERIAL NOT NULL UNIQUE,
    "email" SERIAL NOT NULl UNIQUE,
    "username" SERIAL NOT NULL UNIQUE,
    "password_hash" SERIAL NOT NULL,
    "registered_at" TIMESTAMP NOT NULL DEFAULT NOW()
)