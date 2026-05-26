CREATE TABLE IF NOT EXISTS update (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    body TEXT NOT NULL,
    created_at timestamp(0) with time zone NOT NULL DEFAULT now(),
    updated_at timestamp(0) with time zone NOT NULL DEFAULT now(),
    user_id integer NOT NULL REFERENCES users(id),
    version integer NOT NULL DEFAULT 1
);