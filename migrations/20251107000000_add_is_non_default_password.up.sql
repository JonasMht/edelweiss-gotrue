-- Add is_non_default_password column to track user-set vs default passwords
-- Merged from appflowy branch

ALTER TABLE {{ index .Options "Namespace" }}.users
ADD COLUMN IF NOT EXISTS is_non_default_password boolean NOT NULL DEFAULT false;
