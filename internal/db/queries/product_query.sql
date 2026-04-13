
-- name: CreateProduct :one
INSERT INTO products (
  name, slug,
  description, short_description, display_image_url,
  price_cents, currency,
  status
) VALUES (
  @name, @slug,
  @description, @short_description, @display_image_url,
  @price_cents, @currency,
  @status
)
RETURNING *;

-- name: UpdateProduct :one
UPDATE products
SET
  name = $2,
  slug = $3,
  description = $4,
  short_description = $5,
  display_image_url = $6,
  price_cents = $7,
  currency = $8
WHERE id = $1
  AND deleted_at IS NULL
RETURNING *;

-- name: UpdateStatusProduct :one
UPDATE products
SET
  status = $2
WHERE id = $1
  AND deleted_at IS NULL
RETURNING *;

-- name: DeleteProduct :execrows
UPDATE products
SET 
  deleted_at = NOW()
WHERE id = $1
  AND deleted_at IS NULL;

-- name: GetProductByID :one
SELECT * FROM products
WHERE id = $1
  AND deleted_at IS NULL;

-- name: GetProductBySlug :one
SELECT * FROM products
WHERE slug = $1
  AND deleted_at IS NULL;

-- name: GetListActiveProducts :many
SELECT * FROM products
WHERE status = 'active'
  AND deleted_at IS NULL
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;