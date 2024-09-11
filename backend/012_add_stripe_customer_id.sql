-- Write your migrate up statements here
ALTER TABLE users ADD COLUMN stripe_customer_id VARCHAR(255);
---- create above / drop below ----
ALTER TABLE users DROP COLUMN stripe_customer_id;
-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
