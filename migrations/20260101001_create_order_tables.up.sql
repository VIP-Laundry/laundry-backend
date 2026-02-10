-- 7. Tabel ORDERS (Nota Induk)
CREATE TABLE `orders` (
	`id` BIGINT(19) NOT NULL AUTO_INCREMENT,
	`invoice_number` VARCHAR(50) NOT NULL COLLATE 'utf8mb4_0900_ai_ci',
	`customer_id` BIGINT(19) NULL DEFAULT NULL,
	`customer_name` VARCHAR(150) NULL DEFAULT NULL COLLATE 'utf8mb4_0900_ai_ci',
	`customer_phone` VARCHAR(30) NULL DEFAULT NULL COLLATE 'utf8mb4_0900_ai_ci',
	`customer_address` TEXT NULL DEFAULT NULL COLLATE 'utf8mb4_0900_ai_ci',
	`is_delivery` TINYINT(1) NULL DEFAULT '0',
	`total_price` DECIMAL(15,2) NULL DEFAULT '0.00',
	`payment_status` ENUM('unpaid','paid','cod_pending') NULL DEFAULT 'unpaid' COLLATE 'utf8mb4_0900_ai_ci',
	`status_internal` ENUM('pending','in-progress','ready-pickup','ready-delivery','being-delivered','finished-delivery','picked-up','cancelled') NULL DEFAULT 'pending' COLLATE 'utf8mb4_0900_ai_ci',
	`estimated_ready_at` TIMESTAMP NULL DEFAULT NULL,
	`notes` TEXT NULL DEFAULT NULL COLLATE 'utf8mb4_0900_ai_ci',
	`created_by` BIGINT(19) NULL DEFAULT NULL,
	`created_at` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
	`updated_at` TIMESTAMP NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
	PRIMARY KEY (`id`) USING BTREE,
	UNIQUE INDEX `unique_invoice_number` (`invoice_number`) USING BTREE,
	INDEX `idx_orders_status` (`status_internal`) USING BTREE,
	INDEX `created_by` (`created_by`) USING BTREE,
	INDEX `fk_orders_customer` (`customer_id`) USING BTREE,
	INDEX `idx_orders_created_at` (`created_at`) USING BTREE,
	CONSTRAINT `fk_orders_customer` FOREIGN KEY (`customer_id`) REFERENCES `customers` (`id`) ON UPDATE NO ACTION ON DELETE SET NULL,
	CONSTRAINT `orders_ibfk_1` FOREIGN KEY (`created_by`) REFERENCES `users` (`id`) ON UPDATE NO ACTION ON DELETE SET NULL
)
COLLATE='utf8mb4_0900_ai_ci'
ENGINE=InnoDB
AUTO_INCREMENT=6
;

-- 8. Tabel ORDER ITEMS (Rincian Cucian)
CREATE TABLE `order_items` (
	`id` BIGINT(19) NOT NULL AUTO_INCREMENT,
	`order_id` BIGINT(19) NOT NULL,
	`service_id` BIGINT(19) NULL DEFAULT NULL,
	`item_notes` VARCHAR(255) NULL DEFAULT NULL COLLATE 'utf8mb4_0900_ai_ci',
	`quantity` INT(10) NULL DEFAULT NULL,
	`qty_pieces` INT(10) NULL DEFAULT NULL,
	`weight_kg` DECIMAL(8,2) NULL DEFAULT NULL,
	`unit_price` DECIMAL(15,2) NOT NULL,
	`subtotal` DECIMAL(15,2) NOT NULL,
	PRIMARY KEY (`id`) USING BTREE,
	INDEX `order_id` (`order_id`) USING BTREE,
	INDEX `order_items_ibfk_2` (`service_id`) USING BTREE,
	CONSTRAINT `order_items_ibfk_1` FOREIGN KEY (`order_id`) REFERENCES `orders` (`id`) ON UPDATE NO ACTION ON DELETE CASCADE,
	CONSTRAINT `order_items_ibfk_2` FOREIGN KEY (`service_id`) REFERENCES `services` (`id`) ON UPDATE NO ACTION ON DELETE RESTRICT
)
COLLATE='utf8mb4_0900_ai_ci'
ENGINE=InnoDB
AUTO_INCREMENT=5
;

-- 9. Tabel PAYMENTS (Riwayat Bayar)
CREATE TABLE `payments` (
	`id` BIGINT(19) NOT NULL AUTO_INCREMENT,
	`order_id` BIGINT(19) NOT NULL,
	`method` ENUM('cash','transfer','qris','ewallet') NULL DEFAULT NULL COLLATE 'utf8mb4_0900_ai_ci',
	`amount` DECIMAL(15,2) NOT NULL,
	`amount_received` DECIMAL(15,2) NOT NULL DEFAULT '0.00',
	`amount_change` DECIMAL(15,2) NOT NULL DEFAULT '0.00',
	`reference_no` VARCHAR(100) NULL DEFAULT NULL COLLATE 'utf8mb4_0900_ai_ci',
	`status` ENUM('pending','confirmed','void') NOT NULL DEFAULT 'pending' COLLATE 'utf8mb4_0900_ai_ci',
	`created_by` BIGINT(19) NOT NULL,
	`collected_by` BIGINT(19) NULL DEFAULT NULL,
	`collected_at` TIMESTAMP NULL DEFAULT NULL,
	`created_at` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
	`updated_at` TIMESTAMP NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
	PRIMARY KEY (`id`) USING BTREE,
	INDEX `collected_by` (`collected_by`) USING BTREE,
	INDEX `payments_ibfk_1` (`order_id`) USING BTREE,
	INDEX `fk_payments_creator` (`created_by`) USING BTREE,
	INDEX `idx_payments_created_at` (`created_at`) USING BTREE,
	CONSTRAINT `fk_payments_creator` FOREIGN KEY (`created_by`) REFERENCES `users` (`id`) ON UPDATE NO ACTION ON DELETE RESTRICT,
	CONSTRAINT `payments_ibfk_1` FOREIGN KEY (`order_id`) REFERENCES `orders` (`id`) ON UPDATE NO ACTION ON DELETE RESTRICT,
	CONSTRAINT `payments_ibfk_2` FOREIGN KEY (`collected_by`) REFERENCES `users` (`id`) ON UPDATE NO ACTION ON DELETE SET NULL
)
COLLATE='utf8mb4_0900_ai_ci'
ENGINE=InnoDB
AUTO_INCREMENT=5
;

-- 10. Tabel DELIVERIES (Antar Jemput)
CREATE TABLE `deliveries` (
	`id` BIGINT(19) NOT NULL AUTO_INCREMENT,
	`order_id` BIGINT(19) NOT NULL,
	`delivery_status` ENUM('ready-delivery','being-delivered','finished-delivery') NULL DEFAULT NULL COLLATE 'utf8mb4_0900_ai_ci',
	`shipping_cost` DECIMAL(15,2) NOT NULL DEFAULT '0.00',
	`courier_id` BIGINT(19) NULL DEFAULT NULL,
	`courier_departed_at` TIMESTAMP NULL DEFAULT NULL,
	`courier_arrived_at` TIMESTAMP NULL DEFAULT NULL,
	`receiver_name` VARCHAR(100) NULL DEFAULT NULL COLLATE 'utf8mb4_0900_ai_ci',
	`cod_collected_amount` DECIMAL(15,2) NULL DEFAULT '0.00',
	`created_at` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
	`updated_at` TIMESTAMP NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
	PRIMARY KEY (`id`) USING BTREE,
	UNIQUE INDEX `order_id` (`order_id`) USING BTREE,
	INDEX `courier_id` (`courier_id`) USING BTREE,
	INDEX `delivery_status_idx` (`delivery_status`) USING BTREE,
	CONSTRAINT `deliveries_ibfk_1` FOREIGN KEY (`order_id`) REFERENCES `orders` (`id`) ON UPDATE NO ACTION ON DELETE CASCADE,
	CONSTRAINT `deliveries_ibfk_2` FOREIGN KEY (`courier_id`) REFERENCES `users` (`id`) ON UPDATE NO ACTION ON DELETE SET NULL
)
COLLATE='utf8mb4_0900_ai_ci'
ENGINE=InnoDB
AUTO_INCREMENT=4
;

-- 11. Tabel STATUS HISTORY (Log Perubahan Status Cucian)
CREATE TABLE `status_history` (
	`id` BIGINT(19) NOT NULL AUTO_INCREMENT,
	`order_id` BIGINT(19) NOT NULL,
	`previous_status` VARCHAR(50) NULL DEFAULT NULL COLLATE 'utf8mb4_0900_ai_ci',
	`new_status` VARCHAR(50) NOT NULL COLLATE 'utf8mb4_0900_ai_ci',
	`actor_id` BIGINT(19) NULL DEFAULT NULL,
	`actor_role` VARCHAR(50) NULL DEFAULT NULL COLLATE 'utf8mb4_0900_ai_ci',
	`notes` TEXT NULL DEFAULT NULL COLLATE 'utf8mb4_0900_ai_ci',
	`created_at` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY (`id`) USING BTREE,
	INDEX `order_id` (`order_id`) USING BTREE,
	INDEX `actor_id` (`actor_id`) USING BTREE,
	INDEX `idx_status_history_created_at` (`created_at`) USING BTREE,
	CONSTRAINT `status_history_ibfk_1` FOREIGN KEY (`order_id`) REFERENCES `orders` (`id`) ON UPDATE NO ACTION ON DELETE CASCADE,
	CONSTRAINT `status_history_ibfk_2` FOREIGN KEY (`actor_id`) REFERENCES `users` (`id`) ON UPDATE NO ACTION ON DELETE SET NULL
)
COLLATE='utf8mb4_0900_ai_ci'
ENGINE=InnoDB
AUTO_INCREMENT=15
;
