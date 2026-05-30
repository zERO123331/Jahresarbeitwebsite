CREATE INDEX IF NOT EXISTS shop_title_idx ON shopentry USING GIN (to_tsvector('simple', title));
CREATE INDEX IF NOT EXISTS shop_quantity_idx ON shopentry (quantity);
CREATE INDEX IF NOT EXISTS shop_price_idx ON shopentry (price);