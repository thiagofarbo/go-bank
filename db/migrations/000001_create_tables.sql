CREATE TABLE client (
            id SERIAL PRIMARY KEY,
            name VARCHAR(100),
            age INT,
            email VARCHAR(100) UNIQUE NOT NULL,
            created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP;
);

CREATE TABLE accounts (
         id SERIAL PRIMARY KEY,
         client_id INT NOT NULL,
         number VARCHAR(20) UNIQUE NOT NULL,
         balance NUMERIC(10, 2) NOT NULL DEFAULT 0.00,
         status VARCHAR(25) DEFAULT 'active',
         created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
         updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
         FOREIGN KEY (client_id) REFERENCES client(id);
);

CREATE TABLE transactions (
          id SERIAL PRIMARY KEY,
          account_id INTEGER NOT NULL,
          amount NUMERIC(10, 2) NOT NULL DEFAULT 0.00,
          transactionType VARCHAR(50) NOT NULL,
          created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
          description TEXT,
          FOREIGN KEY (account_id) REFERENCES accounts(id);
);

CREATE TABLE loans (
           id SERIAL PRIMARY KEY,
           account_id INTEGER NOT NULL,
           amount DECIMAL(15, 2) NOT NULL,
           interest_rate DECIMAL(5, 2) NOT NULL,
           description VARCHAR(255),
           term INTEGER NOT NULL,
           created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
           updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
           FOREIGN KEY (account_id) REFERENCES accounts(id);
);

CREATE TABLE gross_incomes (
           id SERIAL PRIMARY KEY,
           account_id INTEGER NOT NULL,
           amount DECIMAL(15, 2) NOT NULL,
           created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
           updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
           FOREIGN KEY (account_id) REFERENCES accounts(id);
);

CREATE TABLE users (
           id SERIAL PRIMARY KEY,
           username VARCHAR(100) NOT NULL,
           password VARCHAR(255) NOT NULL,
           email VARCHAR(50) NOT NULL,
           created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
           updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
