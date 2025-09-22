-- Удаляем поле
ALTER TABLE urls
DROP COLUMN correlation_id;

-- Удаляем свойство текущего поля
ALTER TABLE urls DROP CONSTRAINT url_unique_short;