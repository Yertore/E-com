-- 000001_init_schema.up.sql
-- Catalog Service: categories, products, inventory

CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE categories (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name        VARCHAR(255) NOT NULL,
    slug        VARCHAR(255) NOT NULL UNIQUE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE products (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    category_id  UUID NOT NULL REFERENCES categories(id) ON DELETE RESTRICT,
    sku          VARCHAR(64) NOT NULL UNIQUE,
    name         VARCHAR(255) NOT NULL,
    description  TEXT,
    price        NUMERIC(12,2) NOT NULL CHECK (price >= 0),
    currency     CHAR(3) NOT NULL DEFAULT 'KZT',
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Внешний ключ category_id будет постоянно использоваться в JOIN и
-- WHERE category_id = ... при листинге товаров по категории — индекс обязателен,
-- иначе Postgres будет делать full scan на каждый запрос каталога.
CREATE INDEX idx_products_category_id ON products(category_id);

-- SKU уже уникален (UNIQUE создаёт индекс автоматически), но если в реальности
-- по нему ищут чаще, чем по id, это уже покрыто.

CREATE TABLE inventory (
    product_id  UUID PRIMARY KEY REFERENCES products(id) ON DELETE CASCADE,
    quantity    INTEGER NOT NULL DEFAULT 0 CHECK (quantity >= 0),
    reserved    INTEGER NOT NULL DEFAULT 0 CHECK (reserved >= 0),
    version     INTEGER NOT NULL DEFAULT 0,
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    CHECK (reserved <= quantity)
);
