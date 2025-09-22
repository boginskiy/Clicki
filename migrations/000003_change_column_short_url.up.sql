-- 000003_change_column_short_url.up.sql
-- Обновляем свойство текущего поля
ALTER TABLE urls
ALTER COLUMN short_url TYPE TEXT;