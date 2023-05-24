-- 20201125101852_create_users_table.up.sql
CREATE TABLE users
(
    acct VARCHAR(255) PRIMARY KEY NOT NULL,
    pwd VARCHAR(255) NOT NULL,
    fullname VARCHAR(255) NOT NULL,
    created_at timestamp,
    updated_at timestamp
)
