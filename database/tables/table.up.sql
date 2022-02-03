CREATE TABLE IF NOT EXISTS account_data(
id SERIAL PRIMARY KEY,
name TEXT NOT NULL,
currency VARCHAR(3), /* Currency follows ISO 4217 standard */
balance NUMERIC,
created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP

constraint balance_non_negative check (balance >= 0.0)
);
