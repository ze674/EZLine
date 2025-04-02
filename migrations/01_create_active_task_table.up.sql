-- migrations/01_create_active_task_table.up.sql
CREATE TABLE active_task (
                             id INTEGER PRIMARY KEY AUTOINCREMENT,
                             task_id INTEGER NOT NULL,
                             created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);