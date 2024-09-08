-- Write your migrate up statements here
CREATE TABLE links (
    id SERIAL PRIMARY KEY,
    original TEXT NOT NULL,
    short TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by TEXT NOT NULL
);

---- create above / drop below ----

DROP TABLE links;
-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
