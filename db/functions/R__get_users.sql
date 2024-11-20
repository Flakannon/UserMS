CREATE OR REPLACE FUNCTION get_users(
    p_country TEXT DEFAULT NULL,
    p_email TEXT DEFAULT NULL,
    p_first_name TEXT DEFAULT NULL,
    p_last_name TEXT DEFAULT NULL,
    p_nick_name TEXT DEFAULT NULL,
    p_page INT DEFAULT 1,
    p_page_size INT DEFAULT 10
)
RETURNS TABLE (
    id UUID,
    first_name TEXT,
    last_name TEXT,
    nick_name TEXT,
    email TEXT,
    country TEXT,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
)
LANGUAGE PLPGSQL
AS $$
BEGIN
    -- Validate pagination inputs
    IF p_page < 1 THEN
        RAISE EXCEPTION 'Invalid input: page must be >= 1.';
    END IF;

    IF p_page_size < 1 THEN
        RAISE EXCEPTION 'Invalid input: page_size must be >= 1.';
    END IF;

    -- Return with dynamic filtering - partial matching can be applied
    RETURN QUERY
    SELECT *
    FROM users
    WHERE
        (p_country IS NULL OR country = p_country) AND
        (p_email IS NULL OR email = p_email) AND
        (p_first_name IS NULL OR first_name ILIKE '%' || p_first_name || '%') AND
        (p_last_name IS NULL OR last_name ILIKE '%' || p_last_name || '%') AND
        (p_nick_name IS NULL OR nick_name ILIKE '%' || p_nick_name || '%')
    ORDER BY created_at DESC
    LIMIT p_page_size
    OFFSET (p_page - 1) * p_page_size;
END;
$$;