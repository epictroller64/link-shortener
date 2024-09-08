-- Write your migrate up statements here
-- Add clicks column to links table
ALTER TABLE links
ADD COLUMN clicks INTEGER NOT NULL DEFAULT 0;

-- Create an index on the clicks column for better performance
CREATE INDEX idx_links_clicks ON links(clicks);

---- create above / drop below ----

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
