-- +goose Up
-- +goose StatementBegin
CREATE INDEX idx_product_promotions_product ON product_promotions (product_id);
CREATE INDEX idx_product_promotions_promotion ON product_promotions (promotion_id);
CREATE INDEX idx_promotions_date ON promotions (start_date, end_date);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_product_promotions_product;
DROP INDEX IF EXISTS idx_product_promotions_promotion;
DROP INDEX IF EXISTS idx_promotions_date;
-- +goose StatementEnd
