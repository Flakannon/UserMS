DROP PROCEDURE IF EXISTS update_user;

CREATE PROCEDURE update_user(
    p_id UUID,
    p_first_name TEXT DEFAULT NULL,
    p_last_name TEXT DEFAULT NULL,
    p_nick_name TEXT DEFAULT NULL,
    p_password TEXT DEFAULT NULL,
    p_email TEXT DEFAULT NULL,
    p_country TEXT DEFAULT NULL
)
LANGUAGE PLPGSQL
AS $$
BEGIN
    -- Validate required inputs
    IF p_id IS NULL THEN
        RAISE EXCEPTION 'Invalid input: id is required.';
    END IF;

    UPDATE users
    SET
        first_name = COALESCE(p_first_name, first_name),
        last_name = COALESCE(p_last_name, last_name),
        nick_name = COALESCE(p_nick_name, nick_name),
        password = COALESCE(p_password, password),
        email = COALESCE(p_email, email),
        country = COALESCE(p_country, country)
    WHERE id = p_id;

    IF NOT FOUND THEN
        RAISE EXCEPTION 'User with id % not found.', p_id;
    END IF;
END;
$$;