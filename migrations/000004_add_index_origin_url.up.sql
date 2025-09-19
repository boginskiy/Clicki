-- 000004_add_index_origin_url.up.sql

-- Обновляем свойство текущего поля
ALTER TABLE urls
ADD CONSTRAINT url_unique_original UNIQUE(original_url);
