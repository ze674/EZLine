-- migrations/04_add_serial_number_to_container.down.sql
ALTER TABLE containers DROP COLUMN serial_number;