-- Write your migrate up statements here
-- Modify the created_by column in the links table



ALTER TABLE links
ALTER COLUMN created_by TYPE UUID USING created_by::uuid;

-- Update the foreign key constraint if it exists
ALTER TABLE links
DROP CONSTRAINT IF EXISTS links_created_by_fkey,
ADD CONSTRAINT links_created_by_fkey
FOREIGN KEY (created_by) REFERENCES users(id);



---- create above / drop below ----

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
