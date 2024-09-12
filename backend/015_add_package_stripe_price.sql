-- Write your migrate up statements here
ALTER TABLE packages ADD COLUMN price_id VARCHAR(255);
---- create above / drop below ----
ALTER TABLE packages DROP COLUMN price_id;
-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
