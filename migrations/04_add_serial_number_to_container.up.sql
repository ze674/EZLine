-- Миграция для добавления поля serial_number как INTEGER
ALTER TABLE containers ADD COLUMN serial_number INTEGER NOT NULL DEFAULT 0;