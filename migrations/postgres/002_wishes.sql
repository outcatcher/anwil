-- +goose Up

CREATE TYPE visibility AS ENUM ('private', 'direct_link', 'public');

CREATE TABLE wishlists
(
    "uuid"        UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "wisher_uuid" UUID REFERENCES wishers (uuid),
    "name"        VARCHAR NOT NULL DEFAULT 'default',
    "visibility"  VARCHAR NOT NULL DEFAULT 'private',
    "position"    INT4    NOT NULL DEFAULT 0,

    UNIQUE (wisher_uuid, name),
    UNIQUE (wisher_uuid, position)
);

CREATE TABLE wishes
(
    "uuid"          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "wishlist_uuid" UUID    NOT NULL REFERENCES wishlists (uuid),
    "description"   VARCHAR NOT NULL,
    "position"      INT4    NOT NULL,
    "fulfilled"     BOOLEAN NOT NULL DEFAULT FALSE
);

-- +goose Down

DROP TABLE wishes;
DROP TABLE wishlists;
DROP TYPE visibility;
