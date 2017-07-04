-- psql -h 127.0.0.1 -p 5432 -U postgres -f resources/migrations/init.sql

CREATE DATABASE pascont;

\c pascont;

CREATE TABLE account (
  id            BIGSERIAL,
  name          VARCHAR(64) UNIQUE       NOT NULL,
  password_hash BYTEA                    NOT NULL,
  created_at    TIMESTAMP WITH TIME ZONE NOT NULL,
  updated_at    TIMESTAMP WITH TIME ZONE NOT NULL,
  PRIMARY KEY (id)
);

CREATE TABLE session (
  id         UUID UNIQUE              NOT NULL,
  account_id BIGINT                   NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL,
  expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
  PRIMARY KEY (id),
  FOREIGN KEY (account_id) REFERENCES account (id) ON DELETE CASCADE
);