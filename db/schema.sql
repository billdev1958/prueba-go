CREATE TABLE IF NOT EXISTS comercios (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    comission_rate NUMERIC(19, 4) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS transactions (
    id UUID PRIMARY KEY,
    comercio_id UUID NOT NULL,
    amount NUMERIC(19, 4) NOT NULL,
    applied_rate NUMERIC(19, 4) NOT NULL,
    commission NUMERIC(19, 4) NOT NULL,
    net_amount NUMERIC(19, 4) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_comercio FOREIGN KEY (comercio_id) REFERENCES comercios(id)
);

CREATE TABLE IF NOT EXISTS audit_logs (
    log_id UUID PRIMARY KEY,
    action VARCHAR(100) NOT NULL,
    actor VARCHAR(100) NOT NULL,
    resource_id UUID NOT NULL,
    timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_transactions_comercio_id ON transactions(comercio_id);