-- migrations/02_create_containers_table.up.sql
CREATE TABLE containers (
                            id INTEGER PRIMARY KEY AUTOINCREMENT,
                            code TEXT NOT NULL UNIQUE,             -- Уникальный код контейнера
                            task_id INTEGER NOT NULL,              -- К какому заданию относится
                            status TEXT NOT NULL,                  -- Статус (создан, отправлен и т.д.)
                            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                            updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Индекс для быстрого поиска по заданию
CREATE INDEX idx_containers_task_id ON containers(task_id);