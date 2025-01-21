ALTER TABLE `users`
ADD FOREIGN KEY (`consumer_id`) REFERENCES `consumers`(`id`)
ON DELETE SET NULL;