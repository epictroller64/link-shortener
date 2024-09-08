-- Write your migrate up statements here
ALTER TABLE links ADD COLUMN short_id VARCHAR(10) NOT NULL DEFAULT '';
ALTER TABLE clicks ADD COLUMN ip VARCHAR(40) NOT NULL DEFAULT '';
ALTER TABLE clicks ADD COLUMN country VARCHAR(40) NOT NULL DEFAULT '';
ALTER TABLE clicks ADD COLUMN referer VARCHAR(255) NOT NULL DEFAULT '';
ALTER TABLE clicks ADD COLUMN user_agent VARCHAR(255) NOT NULL DEFAULT '';
---- create above / drop below ----

ALTER TABLE links DROP COLUMN short_id;
ALTER TABLE clicks DROP COLUMN ip;
ALTER TABLE clicks DROP COLUMN country;
ALTER TABLE clicks DROP COLUMN referer;
ALTER TABLE clicks DROP COLUMN user_agent;

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
