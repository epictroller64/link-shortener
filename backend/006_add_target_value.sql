-- Write your migrate up statements here
ALTER TABLE redirects
ADD COLUMN target_value TEXT;

---- create above / drop below ----

ALTER TABLE redirects
DROP COLUMN target_value;

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
