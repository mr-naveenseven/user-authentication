-- This file is used to revert the changes made in the corresponding up migration file.

-- Revert the creation of the user_accounts table
DROP TABLE IF EXISTS user_accounts;