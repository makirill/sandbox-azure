
SET client_min_messages TO warning;

BEGIN;

CREATE TYPE public.status AS ENUM (
    'Running',
    'Stopped',
    'Expired',
    'Pending',
    'Failed',
    'Deleted'
);

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE sandboxes (
    id uuid DEFAULT uuid_generate_v4() CONSTRAINT sandboxes_pk PRIMARY KEY,
    name varchar(100) NOT NULL CHECK (name <> ''),
    created_at timestamp NOT NULL DEFAULT now(),
    updated_at timestamp NOT NULL DEFAULT now(),
    expires_at timestamp NOT NULL,
    status public.status NOT NULL
);

/*
TODO: add permissions to sandboxes table
e.g.: GRANT SELECT ON TABLE sandboxes TO public;
*/


COMMIT;
