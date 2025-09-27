-- name: GetProductPromotionsByProductID :many
select p.promotion_type, pp.*
from promotions p join product_promotions pp on pp.promotion_id = p.id
where p.start_date <= CURRENT_DATE
  and p.end_date >= CURRENT_DATE
  and pp.product_id = $1
  and pp.deleted_at is null
  and p.deleted_at is null;