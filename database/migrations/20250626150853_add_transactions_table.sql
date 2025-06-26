-- +goose Up
-- +goose StatementBegin
CREATE TABLE transactions (
    id VARCHAR(26) PRIMARY KEY,
    wallet_id VARCHAR(26) NOT NULL,
    amount INTEGER NOT NULL,
    note TEXT NULL,
    type VARCHAR(10) NOT NULL,
    status VARCHAR(20) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (wallet_id) REFERENCES wallets(id)
);

CREATE INDEX idx_transactions_wallet_id_created_at ON transactions(wallet_id, created_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS transactions;
-- +goose StatementEnd
