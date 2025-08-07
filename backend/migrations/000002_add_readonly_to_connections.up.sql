-- Add readonly field to connections table
ALTER TABLE connections ADD COLUMN is_readonly BOOLEAN DEFAULT FALSE;
