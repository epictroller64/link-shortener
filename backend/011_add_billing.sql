-- Write your migrate up statements here
CREATE TABLE payments (
    id VARCHAR(255) PRIMARY KEY,
    amount BIGINT NOT NULL,
    currency VARCHAR(3) NOT NULL,
    status VARCHAR(10) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

-- Add an index on the status column for faster queries
CREATE INDEX idx_payments_status ON payments(status);

-- Add a check constraint to ensure valid status values
ALTER TABLE payments
ADD CONSTRAINT chk_payment_status
CHECK (status IN ('pending', 'success', 'failed'));


CREATE TABLE subscriptions (
    id VARCHAR(255) PRIMARY KEY,
    customer_id VARCHAR(255) NOT NULL,
    status VARCHAR(20) NOT NULL,
    current_period_end TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

-- Add an index on the customer_id column for faster queries
CREATE INDEX idx_subscriptions_customer_id ON subscriptions(customer_id);

-- Add a check constraint to ensure valid status values
ALTER TABLE subscriptions
ADD CONSTRAINT chk_subscription_status
CHECK (status IN ('incomplete', 'incomplete_expired', 'trialing', 'active', 'past_due', 'canceled', 'unpaid', 'paused'));


---- create above / drop below ----
DROP TABLE IF EXISTS subscriptions;
DROP TABLE IF EXISTS payments;

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
