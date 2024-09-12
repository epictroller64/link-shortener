-- Write your migrate up statements here
CREATE TABLE packages (
    id VARCHAR(255) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    price INTEGER NOT NULL,
    max_links INTEGER NOT NULL,
    max_clicks INTEGER NOT NULL,
    custom_domains INTEGER NOT NULL,
    is_default BOOLEAN NOT NULL DEFAULT FALSE
);

-- Add an index on the is_default column for faster queries
CREATE INDEX idx_packages_is_default ON packages(is_default);

-- Add a check constraint to ensure non-negative values for numeric fields
ALTER TABLE packages
ADD CONSTRAINT chk_packages_positive_values
CHECK (
    price >= 0 AND
    max_links >= 0 AND
    max_clicks >= 0 AND
    custom_domains >= 0
);

-- Add a unique constraint to ensure only one default package
CREATE UNIQUE INDEX idx_packages_default
ON packages ((is_default IS TRUE))
WHERE is_default IS TRUE;

-- Add foreign key to subscriptions table
ALTER TABLE subscriptions
ADD COLUMN package_id VARCHAR(255) REFERENCES packages(id);

-- Add an index on the package_id column in subscriptions table
CREATE INDEX idx_subscriptions_package_id ON subscriptions(package_id);

---- create above / drop below ----
DROP TABLE IF EXISTS packages;
-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
