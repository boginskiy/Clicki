-- Добавляем новую таблицу с users
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT,
    password TEXT,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP,
    last_login_at TIMESTAMP,
    is_active BOOLEAN DEFAULT TRUE,
    roles TEXT[]
);

-- Обновляем таблицу urls. Добавим столбец и связь с таблицей users
ALTER TABLE urls ADD COLUMN user_id INT;
ALTER TABLE urls ADD FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
