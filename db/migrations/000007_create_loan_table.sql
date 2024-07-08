CREATE TABLE loans (
       id SERIAL PRIMARY KEY,
       account_id INTEGER NOT NULL,
       amount DECIMAL(15, 2) NOT NULL,
       interest_rate DECIMAL(5, 2) NOT NULL,
       description VARCHAR(255),
       term INTEGER NOT NULL,
       created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
       updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
       FOREIGN KEY (account_id) REFERENCES accounts (id) ON DELETE CASCADE
);

CREATE TABLE gross_incomes (
       id SERIAL PRIMARY KEY,
       account_id INTEGER NOT NULL,
       amount DECIMAL(15, 2) NOT NULL,
       created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
       updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
       FOREIGN KEY (account_id) REFERENCES accounts (id) ON DELETE CASCADE,
);