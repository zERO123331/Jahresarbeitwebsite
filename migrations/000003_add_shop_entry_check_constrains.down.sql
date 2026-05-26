ALTER TABLE shopentry DROP CONSTRAINT IF EXISTS price_positive_check;

ALTER TABLE shopentry DROP CONSTRAINT IF EXISTS quantity_positive_check;

ALTER TABLE shopentry DROP CONSTRAINT IF EXISTS image_url_count_check;

ALTER TABLE shopentry DROP CONSTRAINT IF EXISTS categories_count_check;

ALTER TABLE shopentry DROP CONSTRAINT IF EXISTS description_length_check;

ALTER TABLE shopentry DROP CONSTRAINT IF EXISTS title_length_check;


