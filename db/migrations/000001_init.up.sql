DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'chair_type') THEN
        CREATE TYPE chair_type AS ENUM ('abc', 'cde');
    END IF;
END
$$;

CREATE TABLE main (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    sub_id INTEGER,
    sub_obj VARCHAR(50),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP
);

CREATE TABLE tools (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    main_id INTEGER REFERENCES main(id),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP
);

CREATE TABLE tables (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    main_id INTEGER REFERENCES main(id),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP
);

CREATE TABLE chairs (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    type chair_type NOT NULL,
    main_id INTEGER REFERENCES main(id),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP
);