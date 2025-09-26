-- +goose Up
-- +goose StatementBegin
CREATE TABLE product_promotions (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMPTZ NOT NULL DEFAULT (CURRENT_TIMESTAMP),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT (CURRENT_TIMESTAMP),
  deleted_at TIMESTAMPTZ,
  promotion_id INT NOT NULL REFERENCES promotions(id) ON DELETE CASCADE,
  product_id INT NOT NULL REFERENCES products(id) ON DELETE CASCADE,
  min_quantity INT NOT NULL DEFAULT 1,
  free_product_id INT NULL REFERENCES products(id),
  discount DECIMAL(3, 2) NULL,
  free_quantity INT NULL,
  pay_y INT NULL
);


INSERT INTO product_promotions (promotion_id, product_id, min_quantity, free_product_id, free_quantity) VALUES
(1, 2, 1, 4, 1);

INSERT INTO product_promotions (promotion_id, product_id, min_quantity, pay_y) VALUES
(2, 1, 3, 2);

INSERT INTO product_promotions (promotion_id, product_id, min_quantity, discount) VALUES
(3, 3, 3, 0.10);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE product_promotions;
-- +goose StatementEnd
