-- +goose Up
-- create enum type "product_status"
CREATE TYPE "product_status" AS ENUM ('active', 'inactive', 'draft', 'archived');
-- create "products" table
CREATE TABLE "products" (
  "id" uuid NOT NULL DEFAULT uuidv7(),
  "name" text NOT NULL,
  "slug" text NOT NULL,
  "description" text NULL,
  "short_description" text NULL,
  "display_image_url" text NULL,
  "price_cents" bigint NOT NULL DEFAULT 0,
  "currency" character(3) NOT NULL DEFAULT 'USD',
  "status" "product_status" NOT NULL DEFAULT 'active',
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  "deleted_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "products_currency_check" CHECK (currency ~ '^[A-Z]{3}$'::text),
  CONSTRAINT "products_display_image_url_check" CHECK (display_image_url ~ '^https?://.*$'::text),
  CONSTRAINT "products_name_check" CHECK ((length(TRIM(BOTH FROM name)) >= 1) AND (length(TRIM(BOTH FROM name)) <= 128)),
  CONSTRAINT "products_price_cents_check" CHECK (price_cents >= 0),
  CONSTRAINT "products_short_description_check" CHECK (length(TRIM(BOTH FROM COALESCE(short_description, ''::text))) <= 256),
  CONSTRAINT "products_slug_check" CHECK (((length(TRIM(BOTH FROM slug)) >= 1) AND (length(TRIM(BOTH FROM slug)) <= 128)) AND (slug ~ '^[a-z0-9]([a-z0-9-]*[a-z0-9])?$'::text))
);
-- create index "idx__products__price_cents" to table: "products"
CREATE INDEX "idx__products__price_cents" ON "products" ("price_cents") WHERE ((deleted_at IS NULL) AND (status = 'active'::product_status));
-- create index "idx__products__status__created_at" to table: "products"
CREATE INDEX "idx__products__status__created_at" ON "products" ("status", "created_at" DESC) WHERE (deleted_at IS NULL);
-- create index "idx__products__status__price_cents__created_at" to table: "products"
CREATE INDEX "idx__products__status__price_cents__created_at" ON "products" ("status", "price_cents", "created_at" DESC) WHERE (deleted_at IS NULL);
-- create index "idx__products__uniq__slug" to table: "products"
CREATE UNIQUE INDEX "idx__products__uniq__slug" ON "products" ("slug") WHERE (deleted_at IS NULL);

-- +goose Down
-- reverse: create index "idx__products__uniq__slug" to table: "products"
DROP INDEX "idx__products__uniq__slug";
-- reverse: create index "idx__products__status__price_cents__created_at" to table: "products"
DROP INDEX "idx__products__status__price_cents__created_at";
-- reverse: create index "idx__products__status__created_at" to table: "products"
DROP INDEX "idx__products__status__created_at";
-- reverse: create index "idx__products__price_cents" to table: "products"
DROP INDEX "idx__products__price_cents";
-- reverse: create "products" table
DROP TABLE "products";
-- reverse: create enum type "product_status"
DROP TYPE "product_status";
