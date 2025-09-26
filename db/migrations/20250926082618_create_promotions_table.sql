-- +goose Up
-- +goose StatementBegin
CREATE TABLE promotions (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMPTZ NOT NULL DEFAULT (CURRENT_TIMESTAMP),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT (CURRENT_TIMESTAMP),
  deleted_at TIMESTAMPTZ,
  promotion_type VARCHAR(32) NOT NULL,
  start_date TIMESTAMPTZ NOT NULL,
  end_date TIMESTAMPTZ NOT NULL
);

INSERT INTO promotions (id, promotion_type, start_date, end_date) VALUES
(1, 'FREE_ITEM', '2024-01-01', '2026-12-31'),
(2, 'BUY_X_PAY_Y', '2024-01-01', '2026-12-31'),
(3, 'BULK_DISCOUNT', '2024-01-01', '2026-12-31');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE promotions;
-- +goose StatementEnd
