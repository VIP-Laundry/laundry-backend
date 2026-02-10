-- 1. Tabel USERS (Admin & Staff)
CREATE TABLE `users` (
	`id` BIGINT(19) NOT NULL AUTO_INCREMENT,
	`full_name` VARCHAR(150) NOT NULL COLLATE 'utf8mb4_0900_ai_ci',
	`username` VARCHAR(100) NOT NULL COLLATE 'utf8mb4_0900_ai_ci',
	`email` VARCHAR(150) NOT NULL COLLATE 'utf8mb4_0900_ai_ci',
	`password_hash` VARCHAR(255) NOT NULL COLLATE 'utf8mb4_0900_ai_ci',
	`role` ENUM('owner','cashier','staff','courier') NOT NULL COLLATE 'utf8mb4_0900_ai_ci',
	`phone_number` VARCHAR(30) NOT NULL COLLATE 'utf8mb4_0900_ai_ci',
	`is_active` TINYINT(1) NULL DEFAULT '1',
	`last_login_at` DATETIME NULL DEFAULT NULL,
	`created_at` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
	`updated_at` TIMESTAMP NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
	PRIMARY KEY (`id`) USING BTREE,
	UNIQUE INDEX `username` (`username`) USING BTREE,
	UNIQUE INDEX `email` (`email`) USING BTREE,
	UNIQUE INDEX `phone_number` (`phone_number`) USING BTREE,
	INDEX `idx_users_full_name` (`full_name`) USING BTREE
)
COLLATE='utf8mb4_0900_ai_ci'
ENGINE=InnoDB
AUTO_INCREMENT=10
;

-- 2. Tabel CUSTOMERS (Pelanggan)
CREATE TABLE `customers` (
	`id` BIGINT(19) NOT NULL AUTO_INCREMENT,
	`full_name` VARCHAR(150) NOT NULL COLLATE 'utf8mb4_0900_ai_ci',
	`phone_number` VARCHAR(30) NOT NULL COLLATE 'utf8mb4_0900_ai_ci',
	`address` TEXT NULL DEFAULT NULL COLLATE 'utf8mb4_0900_ai_ci',
	`is_active` TINYINT(1) NULL DEFAULT '1',
	`created_at` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
	`updated_at` TIMESTAMP NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
	PRIMARY KEY (`id`) USING BTREE,
	UNIQUE INDEX `phone` (`phone_number`) USING BTREE,
	INDEX `idx_customers_name` (`full_name`) USING BTREE
)
COLLATE='utf8mb4_0900_ai_ci'
ENGINE=InnoDB
AUTO_INCREMENT=2
;

-- 3. Tabel REFRESH TOKENS (Sesi Login)
CREATE TABLE `refresh_tokens` (
	`id` BIGINT(19) NOT NULL AUTO_INCREMENT,
	`user_id` BIGINT(19) NOT NULL,
	`token` VARCHAR(255) NOT NULL COLLATE 'utf8mb4_0900_ai_ci',
	`expires_at` DATETIME NOT NULL,
	`created_at` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY (`id`) USING BTREE,
	INDEX `idx_refresh_tokens_user` (`user_id`) USING BTREE,
	INDEX `idx_token_value` (`token`) USING BTREE,
	CONSTRAINT `fk_refresh_tokens_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON UPDATE NO ACTION ON DELETE CASCADE
)
COLLATE='utf8mb4_0900_ai_ci'
ENGINE=InnoDB
AUTO_INCREMENT=62
;

-- 4. Tabel TOKEN BLACKLIST (Logout)
CREATE TABLE `token_blacklist` (
	`id` BIGINT(19) NOT NULL AUTO_INCREMENT,
	`jti` VARCHAR(255) NOT NULL COLLATE 'utf8mb4_0900_ai_ci',
	`expires_at` DATETIME NOT NULL,
	`created_at` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY (`id`) USING BTREE,
	INDEX `idx_jti` (`jti`) USING BTREE
)
COLLATE='utf8mb4_0900_ai_ci'
ENGINE=InnoDB
AUTO_INCREMENT=12
;
