-- name: GetProductByID :one
select * from products where id = $1 and deleted_at is null and quantity > 0;

-- name: UpdateProductQuantity :one
update products 
set quantity = quantity - $1, updated_at = CURRENT_TIMESTAMP 
where 
    id = $2 and quantity >= $1 and deleted_at is null
returning *;