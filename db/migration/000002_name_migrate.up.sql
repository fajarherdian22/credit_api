CREATE TABLE `sessions` (
    `id` varchar(255) PRIMARY KEY,
    `email` varchar(255) NOT NULL,
    `customer_id` varchar(255) NOT NULL,
    `refresh_token` text NOT NULL,
    `user_agent` varchar(255) NOT NULL,
    `client_ip` varchar(255) NOT NULL,
    `is_blocked` boolean NOT NULL DEFAULT false,
    `expires_at` datetime NOT NULL,
    `created_at` datetime NOT NULL DEFAULT (now())
);

CREATE TABLE `payment_details` (
  `id` varchar(255) PRIMARY KEY NOT NULL,
  `transaction_id` varchar(255) NOT NULL,
  `amount` DOUBLE NOT NULL,
  `due_date` datetime NOT NULL,
  `is_paid` boolean NOT NULL DEFAULT false
);

ALTER TABLE `sessions` ADD FOREIGN KEY (`customer_id`) REFERENCES `customers` (`id`);
ALTER TABLE `payment_details` ADD FOREIGN KEY (`transaction_id`) REFERENCES `transaction` (`id`);