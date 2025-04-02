-- migrations/03_create_items_table.up.sql
CREATE TABLE items (
                       id INTEGER PRIMARY KEY AUTOINCREMENT,
                       code TEXT NOT NULL UNIQUE,             -- Уникальный код товара
                       task_id INTEGER NOT NULL,              -- К какому заданию относится
                       container_id INTEGER,                  -- Привязка к контейнеру (может быть NULL)
                       status TEXT NOT NULL,                  -- Статус (отсканирован, агрегирован и т.д.)
                       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                       FOREIGN KEY (container_id) REFERENCES containers(id)
);

-- Индексы для быстрого поиска
CREATE INDEX idx_items_task_id ON items(task_id);
CREATE INDEX idx_items_container_id ON items(container_id);