-- Добавляем новое поле
ALTER TABLE urls
ADD COLUMN deleted_flag BOOLEAN NOT NULL DEFAULT FALSE;
