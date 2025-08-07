-- Remove readonly field from connections table
ALTER TABLE connections DROP COLUMN is_readonly;
