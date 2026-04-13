CREATE TYPE product_status AS ENUM ('active', 'inactive', 'draft', 'archived');

CREATE TABLE products (
   id UUID PRIMARY KEY DEFAULT uuidv7(),

   name TEXT NOT NULL CHECK (length(trim(name)) BETWEEN 1 AND 128),
   slug TEXT NOT NULL CHECK (length(trim(slug)) BETWEEN 1 AND 128 AND slug ~ '^[a-z0-9]([a-z0-9-]*[a-z0-9])?$'),

   description       TEXT,
   short_description TEXT CHECK (length(trim(coalesce(short_description, ''))) <= 256),
   display_image_url TEXT CHECK (display_image_url ~ '^https?://.*$'),

   price_cents BIGINT  NOT NULL DEFAULT 0     CHECK (price_cents >= 0),
   currency    CHAR(3) NOT NULL DEFAULT 'USD' CHECK (currency ~ '^[A-Z]{3}$'),

   status product_status NOT NULL DEFAULT 'active',

   created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
   updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
   deleted_at TIMESTAMPTZ  -- soft delete
);

CREATE UNIQUE INDEX idx__products__uniq__slug                      ON products (slug)                                 WHERE deleted_at IS NULL;
CREATE        INDEX idx__products__status__created_at              ON products (status, created_at DESC)              WHERE deleted_at IS NULL;
CREATE        INDEX idx__products__status__price_cents__created_at ON products (status, price_cents, created_at DESC) WHERE deleted_at IS NULL;
CREATE        INDEX idx__products__price_cents                     ON products (price_cents)                          WHERE deleted_at IS NULL AND status = 'active';

CREATE OR REPLACE FUNCTION trigger_set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
   IF NEW.updated_at IS DISTINCT FROM OLD.updated_at THEN
      RETURN NEW;
   END IF;
   NEW.updated_at := NOW();
   RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg__products__set__updated_at
   BEFORE UPDATE ON products
   FOR EACH ROW
   EXECUTE FUNCTION trigger_set_updated_at();