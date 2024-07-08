CREATE TABLE transactions (
      id SERIAL PRIMARY KEY,
      account_id INTEGER NOT NULL,
      amount NUMERIC(10, 2) NOT NULL DEFAULT 0.00,
      transactionType VARCHAR(50) NOT NULL,
      created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
      description TEXT,
      FOREIGN KEY (account_id) REFERENCES accounts(id),
);
