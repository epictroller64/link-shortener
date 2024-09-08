-- Write your migrate up statements here
ALTER TABLE redirects ADD COLUMN target_name TEXT;
---- create above / drop below ----

ALTER TABLE redirects DROP COLUMN target_name;

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
