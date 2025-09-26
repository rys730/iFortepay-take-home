-- +goose Up
-- +goose StatementBegin
CREATE TABLE products (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMPTZ NOT NULL DEFAULT (CURRENT_TIMESTAMP),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT (CURRENT_TIMESTAMP),
  deleted_at TIMESTAMPTZ,
  sku VARCHAR(8) UNIQUE NOT NULL,
  name VARCHAR(128) NOT NULL,
  price DECIMAL(10, 2) NOT NULL,
  quantity INT NOT NULL DEFAULT 0
);

INSERT INTO products (id, sku, name, price, quantity) VALUES
(1, '120P90', 'Google Home', 49.99, 10),
(2, '43N23P', 'MacBook Pro', 5399.99, 5),
(3, 'A304SD', 'Alexa Speaker', 109.50, 10),
(4, '234234', 'Raspberry Pi B', 30.00, 2);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE products;
-- +goose StatementEnd
