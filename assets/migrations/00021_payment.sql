-- +goose Up

DROP TABLE IF EXISTS invoices;
DROP TABLE IF EXISTS shippings;

CREATE TABLE wallets (
	id TEXT PRIMARY KEY,
	user_id TEXT NOT NULL REFERENCES users(id),
	name TEXT NOT NULL,
	balance INTEGER NOT NULL,
	currency TEXT NOT NULL,
	created_at TIMESTAMP NOT NULL,
	updated_at TIMESTAMP NOT NULL,
	deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE transactions (
	id TEXT PRIMARY KEY,
	user_id TEXT NOT NULL REFERENCES users(id),
	wallet_id TEXT NOT NULL REFERENCES wallets(id),
	amount INTEGER NOT NULL,
	currency TEXT NOT NULL,
	description TEXT,
	created_at TIMESTAMP NOT NULL,
	updated_at TIMESTAMP NOT NULL,
	deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE invoices (
	id TEXT PRIMARY KEY,
	user_id TEXT NOT NULL REFERENCES users(id),
	transaction_id TEXT NOT NULL REFERENCES transactions(id),
	subtotal INTEGER NOT NULL,
	discount INTEGER NOT NULL,
	amount INTEGER NOT NULL,
	currency TEXT NOT NULL,
	description TEXT,
	created_at TIMESTAMP NOT NULL,
	updated_at TIMESTAMP NOT NULL,
	deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE invoice_items (
	id TEXT PRIMARY KEY,
	invoice_id TEXT NOT NULL REFERENCES invoices(id),
	product_id TEXT NOT NULL REFERENCES products(id),
	quantity INTEGER NOT NULL,
	price INTEGER NOT NULL,
	discount INTEGER NOT NULL,
	amount INTEGER NOT NULL,
	created_at TIMESTAMP NOT NULL,
	updated_at TIMESTAMP NOT NULL,
	deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE shippings (
	id TEXT PRIMARY KEY,
	user_id TEXT REFERENCES users(id) ON DELETE SET NULL,
	invoice_id TEXT REFERENCES invoices(id) ON DELETE SET NULL,
	fee INTEGER NOT NULL,
	address TEXT NOT NULL,
	note TEXT,
	created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
	deleted_at TIMESTAMP WITH TIME ZONE
);

ALTER TABLE products ADD COLUMN donated BOOLEAN NOT NULL DEFAULT false;

-- +goose Down
DROP TABLE wallets;
DROP TABLE transactions;
DROP TABLE invoices;
DROP TABLE invoice_items;
DROP TABLE shippings;
ALTER TABLE products DROP COLUMN donated;
