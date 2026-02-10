-- 5. Tabel SERVICE CATEGORIES (Cth: Satuan, Kiloan)
CREATE TABLE `service_categories` (
	`id` BIGINT(19) NOT NULL AUTO_INCREMENT,
	`category_name` VARCHAR(150) NOT NULL COLLATE 'utf8mb4_0900_ai_ci',
	`description` VARCHAR(255) NULL DEFAULT NULL COLLATE 'utf8mb4_0900_ai_ci',
	`is_active` TINYINT(1) NULL DEFAULT '1',
	`created_at` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
	`updated_at` TIMESTAMP NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
	PRIMARY KEY (`id`) USING BTREE,
	UNIQUE INDEX `idx_category_name` (`category_name`) USING BTREE
)
COLLATE='utf8mb4_0900_ai_ci'
ENGINE=InnoDB
AUTO_INCREMENT=3
;

-- 6. Tabel SERVICES (Cth: Cuci Setrika 3 Hari)
CREATE TABLE `services` (
	`id` BIGINT(19) NOT NULL AUTO_INCREMENT,
	`code` VARCHAR(50) NOT NULL COLLATE 'utf8mb4_0900_ai_ci',
	`service_name` VARCHAR(150) NOT NULL COLLATE 'utf8mb4_0900_ai_ci',
	`unit` ENUM('kg','pcs') NOT NULL COLLATE 'utf8mb4_0900_ai_ci',
	`price` DECIMAL(15,2) NOT NULL DEFAULT '0.00',
	`is_active` TINYINT(1) NULL DEFAULT '1',
	`created_at` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
	`updated_at` TIMESTAMP NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
	`category_id` BIGINT(19) NOT NULL,
	`duration_hours` INT(10) NOT NULL DEFAULT '72',
	PRIMARY KEY (`id`) USING BTREE,
	UNIQUE INDEX `code` (`code`) USING BTREE,
	INDEX `fk_services_category` (`category_id`) USING BTREE,
	INDEX `idx_services_name` (`service_name`) USING BTREE,
	CONSTRAINT `fk_services_category` FOREIGN KEY (`category_id`) REFERENCES `service_categories` (`id`) ON UPDATE NO ACTION ON DELETE RESTRICT
)
COLLATE='utf8mb4_0900_ai_ci'
ENGINE=InnoDB
AUTO_INCREMENT=3
;