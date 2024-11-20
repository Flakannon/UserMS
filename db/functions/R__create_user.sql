-- DROP PROCEDURE IF EXISTS create_user(
--     UUID,
--     TEXT,
--     TEXT,
--     TEXT,
--     TEXT,
--     TEXT,
--     TEXT
--     TIMESTAMP WITH TIME ZONE,
--     TIMESTAMP WITH TIME ZONE
-- );

DROP PROCEDURE IF EXISTS create_user;

CREATE PROCEDURE create_user(
    p_id UUID,
    p_first_name TEXT,
    p_last_name TEXT,
    p_nick_name TEXT,
    p_password TEXT,
    p_email TEXT,
    p_country TEXT
)
LANGUAGE PLPGSQL
AS $$
BEGIN
    -- Validate required inputs
    IF p_id IS NULL OR p_email IS NULL THEN
        RAISE EXCEPTION 'Invalid input: id and email are required.';
    END IF;

    -- Insert into the user table
    INSERT INTO users (
        id,
        first_name,
        last_name,
        nick_name,
        password,
        email,
        country
    )
    VALUES (
        p_id,
        p_first_name,
        p_last_name,
        p_nick_name,
        p_password,
        p_email,
        p_country
    )
    ON CONFLICT (email) DO NOTHING; -- Optional conflict handling
END;
$$;