-- Добавляем новое поле
ALTER TABLE urls
ADD COLUMN user_id INT NOT NULL;