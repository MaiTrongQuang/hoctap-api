-- HocTap API Database Setup Script
-- Run this script in your MySQL/MariaDB server

-- Create database
CREATE DATABASE IF NOT EXISTS hoctap_api;

-- Create user (optional - you can use root or existing user)
-- Uncomment and modify the lines below if you want to create a dedicated user
-- CREATE USER IF NOT EXISTS 'hoctap_user'@'localhost' IDENTIFIED BY 'your_secure_password';
-- GRANT ALL PRIVILEGES ON hoctap_api.* TO 'hoctap_user'@'localhost';
-- FLUSH PRIVILEGES;

-- Use the database
USE hoctap_api;

-- Create users table (this will also be created automatically by the Go application)
CREATE TABLE IF NOT EXISTS users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Insert some sample data (optional)
INSERT IGNORE INTO users (name, email) VALUES 
('John Doe', 'john@example.com'),
('Jane Smith', 'jane@example.com'),
('Alice Johnson', 'alice@example.com');

-- Show created tables
SHOW TABLES;

-- Show sample data
SELECT * FROM users;

PRINT 'Database setup completed successfully!';
