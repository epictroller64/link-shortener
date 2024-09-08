-- Write your migrate up statements here
CREATE TABLE IF NOT EXISTS clicks (
    id SERIAL PRIMARY KEY,
    link_id INTEGER NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (link_id) REFERENCES links(id) ON DELETE CASCADE
);

CREATE INDEX idx_clicks_link_id ON clicks(link_id);
CREATE INDEX idx_clicks_created_at ON clicks(created_at);

-- Update the link_id column to use INTEGER type
ALTER TABLE clicks ALTER COLUMN link_id TYPE INTEGER;

-- Update the foreign key constraint to match the new type
ALTER TABLE clicks DROP CONSTRAINT clicks_link_id_fkey;
ALTER TABLE clicks ADD CONSTRAINT clicks_link_id_fkey 
    FOREIGN KEY (link_id) REFERENCES links(id) ON DELETE CASCADE;

-- Recreate the index with the updated column type
DROP INDEX IF EXISTS idx_clicks_link_id;
CREATE INDEX idx_clicks_link_id ON clicks(link_id);

---- create above / drop below ----

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
