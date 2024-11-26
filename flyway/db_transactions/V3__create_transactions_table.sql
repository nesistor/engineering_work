DROP TABLE IF EXISTS transactions;

CREATE TABLE transactions (
  id SERIAL PRIMARY KEY,  -- SERIAL to auto-increment w PostgreSQL
  amount INTEGER NOT NULL,
  currency VARCHAR(255) NOT NULL,
  last_four VARCHAR(255) NOT NULL,
  bank_return_code VARCHAR(255) NOT NULL,
  transaction_status_id INTEGER NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,  -- TIMESTAMPTZ przechowuje datę z czasem i strefą czasową
  updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  expiry_month INTEGER NOT NULL DEFAULT 0,
  expiry_year INTEGER NOT NULL DEFAULT 0,
  payment_intent VARCHAR(255) NOT NULL DEFAULT '',
  payment_method VARCHAR(255) NOT NULL DEFAULT '',
  CONSTRAINT transactions_transaction_statuses_id_fk FOREIGN KEY (transaction_status_id)
    REFERENCES transaction_statuses (id) ON DELETE CASCADE ON UPDATE CASCADE
);

-- Opcjonalnie można dodać indeks dla transaction_status_id, jeśli jest wymagany
CREATE INDEX transactions_transaction_statuses_id_idx ON transactions (transaction_status_id);
