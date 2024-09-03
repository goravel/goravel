CREATE TABLE users (
   id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
   name VARCHAR(255) NOT NULL,
   email VARCHAR(255) NOT NULL UNIQUE,
   password VARCHAR(255) NOT NULL,
   email_verified_at TIMESTAMP NULL DEFAULT NULL,
   created_at DATETIME(3) NOT NULL,
   updated_at DATETIME(3) NOT NULL,
   deleted_at DATETIME(3) NULL DEFAULT NULL,
   KEY idx_users_created_at (created_at),
   KEY idx_users_updated_at (updated_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;