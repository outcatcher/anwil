-- +goose Up

CREATE TYPE "role" AS ENUM ('wisher', 'admin');

CREATE TABLE wishers
(
    "uuid"      UUID PRIMARY KEY        DEFAULT gen_random_uuid(),
    "username"  VARCHAR UNIQUE NOT NULL,
    "password"  VARCHAR        NOT NULL,
    "full_name" VARCHAR        NOT NULL,
    "role"      role           NOT NULL DEFAULT 'wisher',
    "enabled"   BOOLEAN        NOT NULL DEFAULT true
);

-- +goose Down

DROP TABLE wishers;

DROP TYPE "role";
