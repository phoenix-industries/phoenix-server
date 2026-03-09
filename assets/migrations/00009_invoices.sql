-- +goose Up
CREATE TYPE invoice_status AS ENUM ('pending', 'delivered', 'cancelled');
CREATE TABLE invoices (
	id TEXT PRIMARY KEY,
	user_id TEXT REFERENCES users(id) ON DELETE SET NULL,
	product_id TEXT REFERENCES products(id) ON DELETE SET NULL,
	status invoice_status NOT NULL DEFAULT 'pending',
	note TEXT,
	created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
	deleted_at TIMESTAMP WITH TIME ZONE
);

-- +goose Down
DROP TABLE invoices;
DROP TYPE invoice_status;
