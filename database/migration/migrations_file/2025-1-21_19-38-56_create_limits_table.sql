CREATE TABLE IF NOT EXISTS `limits` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `consumer_id` bigint(20) unsigned NOT NULL,
  `amount` float DEFAULT 0,
  `tenor` tinyint(2) NOT NULL,
  `created_at` timestamp NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  PRIMARY KEY (`id`),
  FOREIGN KEY (`consumer_id`) REFERENCES `consumers`(`id`) ON DELETE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;