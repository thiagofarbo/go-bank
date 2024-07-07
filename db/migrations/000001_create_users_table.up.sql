CREATE TABLE users (
   id SERIAL PRIMARY KEY,
   name VARCHAR(100),
   age INT,
   email VARCHAR(100) UNIQUE NOT NULL,
   created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
