ALTER TABLE shopentry ADD CONSTRAINT price_positive_check CHECK (price > 0);

ALTER TABLE shopentry ADD CONSTRAINT quantity_positive_check CHECK (quantity > 0);

ALTER TABLE shopentry ADD CONSTRAINT image_url_count_check CHECK (ARRAY_LENGTH(image_urls, 1) BETWEEN 1 AND 10);

ALTER TABLE shopentry ADD CONSTRAINT categories_count_check CHECK (ARRAY_LENGTH(categories, 1) BETWEEN 1 AND 5);

ALTER TABLE shopentry ADD CONSTRAINT description_length_check CHECK (LENGTH(description) BETWEEN 50 AND 2000);

ALTER TABLE shopentry ADD CONSTRAINT title_length_check CHECK (LENGTH(title) BETWEEN 5 AND 100);
