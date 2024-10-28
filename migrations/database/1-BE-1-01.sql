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
    fk_users                 INTEGER REFERENCES "users".users(id_users),
    active                  BOOLEAN DEFAULT FALSE,
    UNIQUE(code_phone_number, phone_number)
);


-- +migrate Down
-- SQL in section 'Down' is executed when this migration is applied


