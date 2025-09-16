-- migrations/000002_add_correlation_id_column.up.sql
-- Добавляем новое поле
ALTER TABLE urls
ADD COLUMN correlation_id TEXT UNIQUE;

-- Обновляем свойство текущего поля
ALTER TABLE urls
ADD CONSTRAINT url_unique_short UNIQUE(short_url);
