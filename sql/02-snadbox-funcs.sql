SET client_min_messages TO warning;

BEGIN;

-- TODO: Check input values
CREATE OR REPLACE FUNCTION insert_sandbox(in_name varchar)
    RETURNS uuid
    LANGUAGE 'plpgsql'
AS
$$
DECLARE
    sandbox_id uuid;
BEGIN
    INSERT INTO sandboxes (name, status)
    VALUES (in_name, 'Pending')
    RETURNING id INTO sandbox_id;

    RETURN sandbox_id;
END;
$$;

CREATE OR REPLACE FUNCTION update_sandbox_status(in_sandbox_id uuid, in_status public.status)
    RETURNS boolean
    LANGUAGE 'plpgsql'
AS
$$
BEGIN
    UPDATE sandboxes
    SET status = in_status,
        updated_at = now()
    WHERE id = in_sandbox_id;

    RETURN FOUND;
END;
$$;

CREATE OR REPLACE FUNCTION update_sandbox_expires_at(in_sandbox_id uuid, in_expires_at timestamp)
    RETURNS boolean
    LANGUAGE 'plpgsql'
AS
$$
BEGIN
    UPDATE sandboxes
    SET expires_at = in_expires_at, 
        updated_at = now()
    WHERE id = in_sandbox_id;

    RETURN FOUND;
END;
$$;


CREATE OR REPLACE FUNCTION delete_sandbox(in_sandbox_id uuid)
    RETURNS boolean
    LANGUAGE 'plpgsql'
AS
$$
BEGIN
    DELETE FROM sandboxes
    WHERE id = in_sandbox_id;

    RETURN FOUND;
END;
$$;

CREATE OR REPLACE FUNCTION get_sandbox_by_id(in_sandbox_id uuid)
    RETURNS table
    (
        id uuid,
        name varchar,
        created_at timestamp,
        updated_at timestamp,
        expires_at timestamp,
        status public.status
    )
    LANGUAGE 'plpgsql'
AS
$$
BEGIN
    RETURN QUERY
    SELECT
        s.id,
        s.name,
        s.created_at,
        s.updated_at,
        s.expires_at,
        s.status
    FROM
        sandboxes s
    WHERE
        s.id = in_sandbox_id;
END;
$$;

CREATE OR REPLACE FUNCTION get_sandbox_by_name(in_name varchar)
    RETURNS table
    (
        id uuid,
        name varchar,
        created_at timestamp,
        updated_at timestamp,
        expires_at timestamp,
        status public.status
    )
    LANGUAGE 'plpgsql'
AS
$$
BEGIN
    RETURN QUERY
    SELECT
        s.id,
        s.name,
        s.created_at,
        s.updated_at,
        s.expires_at,
        s.status
    FROM
        sandboxes s
    WHERE
        s.name = in_name;
END;
$$;

CREATE OR REPLACE FUNCTION get_sandbox_all()
    RETURNS table
    (
        id uuid,
        name varchar,
        created_at timestamp,
        updated_at timestamp,
        expires_at timestamp,
        status public.status
    )
    LANGUAGE 'plpgsql'
AS
$$
BEGIN
    RETURN QUERY
    SELECT
        s.id,
        s.name,
        s.created_at,
        s.updated_at,
        s.expires_at,
        s.status
    FROM
        sandboxes s;
END;
$$;

CREATE OR REPLACE FUNCTION get_sandbox_by_status(in_status public.status)
    RETURNS table
    (
        id uuid,
        name varchar,
        created_at timestamp,
        updated_at timestamp,
        expires_at timestamp,
        status public.status
    )
    LANGUAGE 'plpgsql'
AS
$$
BEGIN
    RETURN QUERY
    SELECT
        s.id,
        s.name,
        s.created_at,
        s.updated_at,
        s.expires_at,
        s.status
    FROM
        sandboxes s
    WHERE
        s.status = in_status;
END;
$$;


COMMIT;