DROP PROCEDURE IF EXISTS delete_user;

CREATE PROCEDURE delete_user(p_id UUID)
LANGUAGE PLPGSQL
AS $$
BEGIN
    -- Validate input
    IF p_id IS NULL THEN
        RAISE EXCEPTION 'Invalid input: id is required.';
    END IF;

    DELETE FROM users
    WHERE id = p_id;

    IF NOT FOUND THEN
        RAISE EXCEPTION 'User with id % not found.', p_id;
    END IF;
END;
$$;