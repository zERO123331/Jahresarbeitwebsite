CREATE TABLE IF NOT EXISTS shopentry (
    id BIGSERIAL PRIMARY KEY,
    created_at timestamp(0) with time zone NOT NULL DEFAULT now(),
    title text NOT NULL,
    description text NOT NULL,
    price integer NOT NULL,
    quantity integer NOT NULL,
    image_urls text[] NOT NULL,
    categories text[] NOT NULL,
    user_id integer NOT NULL REFERENCES users(id)
);