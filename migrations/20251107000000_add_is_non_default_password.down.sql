-- Rollback: Remove is_non_default_password column

ALTER TABLE {{ index .Options "Namespace" }}.users
DROP COLUMN IF EXISTS is_non_default_password;
