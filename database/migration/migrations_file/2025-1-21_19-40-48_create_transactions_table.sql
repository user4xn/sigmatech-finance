CREATE TABLE IF NOT EXISTS `transactions` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `consumer_id` bigint(20) unsigned NOT NULL,
  `contract_no` varchar(255) NOT NULL,
  `otr` float DEFAULT 0,
  `admin_fee` float DEFAULT 0,
  `installment_month` tinyint(2) NOT NULL,
  `interest_percentage` float DEFAULT 0,
  `asset_name` varchar(255),
  `status` ENUM('pending', 'process', 'completed', 'denied') DEFAULT 'pending',
  `created_at` timestamp NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  PRIMARY KEY (`id`),
  FOREIGN KEY (`consumer_id`) REFERENCES `consumers`(`id`) ON DELETE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;