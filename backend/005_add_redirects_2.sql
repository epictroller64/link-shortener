-- Write your migrate up statements here

-- Create redirects table
CREATE TABLE redirects (
    id SERIAL PRIMARY KEY,
    link_id INT NOT NULL,
    target_type VARCHAR(10) NOT NULL,
    target_method VARCHAR(20) NOT NULL,
    redirect_url TEXT NOT NULL,
    FOREIGN KEY (link_id) REFERENCES links(id) ON DELETE CASCADE
);

-- Create an index on the link_id column for better performance
CREATE INDEX idx_redirects_link_id ON redirects(link_id);

-- Add constraints for target_type and target_method
ALTER TABLE redirects
ADD CONSTRAINT check_target_type
CHECK (target_type IN ('header', 'cookie'));

ALTER TABLE redirects
ADD CONSTRAINT check_target_method
CHECK (target_method IN ('match', 'regex', 'contains', 'startsWith', 'endsWith'));

---- create above / drop below ----

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
