DROP FUNCTION IF EXISTS create_user;

CREATE FUNCTION create_user(
    p_first_name TEXT,
    p_last_name TEXT,
    p_nick_name TEXT,
    p_password TEXT,
    p_email TEXT,
    p_country TEXT
)
RETURNS UUID 
LANGUAGE PLPGSQL
AS $$
DECLARE
    new_user_id UUID;
BEGIN
    -- Validate required inputs
    IF p_email IS NULL THEN
        RAISE EXCEPTION 'Invalid input: email is required.';
    END IF;

    -- Insert into the user table
    INSERT INTO users (
        first_name,
        last_name,
        nick_name,
        password,
        email,
        country
    )
    VALUES (
        p_first_name,
        p_last_name,
        p_nick_name,
        p_password,
        p_email,
        p_country
    )
    RETURNING id INTO new_user_id; -- Capture the generated ID

    RETURN new_user_id; -- Return the captured ID
END;
$$;