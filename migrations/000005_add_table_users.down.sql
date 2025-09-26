-- Удаление столбца
ALTER TABLE urls DROP COLUMN user_id CASCADE;
-- Удаление новой таблицы
DROP TABLE IF EXISTS users; 