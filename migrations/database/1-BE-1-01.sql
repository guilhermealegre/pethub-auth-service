-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied


CREATE TABLE "auth".auth
(
    id_auth                 SERIAL PRIMARY KEY,
    email                   VARCHAR(100) UNIQUE,
    code_phone_number       VARCHAR(6),
    phone_number            VARCHAR(20),
    password                BYTEA,
    created_at              timestamptz NOT NULL DEFAULT now(),
    updated_at              timestamptz NOT NULL DEFAULT now(),
    active                  BOOLEAN DEFAULT FALSE,
    uuid_user                    UUID NOT NULL UNIQUE,
    UNIQUE(code_phone_number, phone_number)
);

CREATE INDEX idx_auth_uuid ON "auth".auth(uuid_user);



-- +migrate Down
-- SQL in section 'Down' is executed when this migration is applied

DROP TABLE "auth".auth;