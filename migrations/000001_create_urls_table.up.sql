-- migrations/000001_create_urls_table.up.sql
-- Создание таблицы для хранения URL
CREATE TABLE urls (
    id SERIAL PRIMARY KEY,
    original_url TEXT NOT NULL,
    short_url VARCHAR(8) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Базовый индекс для поиска по code
CREATE INDEX idx_urls_short_url ON urls(short_url);
