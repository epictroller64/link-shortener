-- Write your migrate up statements here
ALTER TABLE packages
ADD COLUMN features JSONB;

---- create above / drop below ----
ALTER TABLE packages
DROP COLUMN features;

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
