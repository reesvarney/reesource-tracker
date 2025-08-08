-- Drop applied_tags table
DROP TABLE IF EXISTS applied_tags;

-- Drop tags table
DROP TABLE IF EXISTS tags;

-- Drop sample_comments table
DROP TABLE IF EXISTS sample_comments;

-- Remove part_number column from products
ALTER TABLE products DROP COLUMN part_number;

-- Remove product_issue column from samples
ALTER TABLE samples DROP COLUMN product_issue;

-- Remove owner_id column from samples
ALTER TABLE samples DROP COLUMN owner_id;

-- Drop users table
DROP TABLE IF EXISTS users;
