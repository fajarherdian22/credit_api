CREATE TABLE customers (
  `id` varchar(255) PRIMARY KEY NOT NULL,
  `nik` VARCHAR(16) UNIQUE NOT NULL NOT NULL,
  `hashed_password` varchar(255) NOT NULL,
  `email` varchar(255) UNIQUE NOT NULL,
  `full_name` varchar(255)  NOT NULL,
  `legal_name` varchar(255) NOT NULL,
  `tempat_lahir` varchar(255) NOT NULL,
  `tanggal_lahir` date NOT NULL,
  `gaji` DOUBLE NOT NULL,
  `photo_ktp` text NOT NULL,
  `foto_selfie` text NOT NULL
);

CREATE TABLE loan_limit (
  `id` varchar(255) PRIMARY KEY NOT NULL,
  `customer_id` varchar(255) NOT NULL,
  `tenor` int NOT NULL,
  `limit` DOUBLE NOT NULL
);

CREATE TABLE `transaction` (
  `id` varchar(255) PRIMARY KEY NOT NULL,
  `customer_id` varchar(255) NOT NULL,
  `product_name` varchar(255) NOT NULL,
  `price` DOUBLE NOT NULL,
  `bunga` DOUBLE NOT NULL,
  `jumlah_cicilan` DOUBLE NOT NULL,
  `tenor` int NOT NULL,
  `admin_fee` DOUBLE NOT NULL,
  `created_at` datetime NOT NULL DEFAULT (now())
);

ALTER TABLE `loan_limit` ADD FOREIGN KEY (`customer_id`) REFERENCES `customers` (`id`);
ALTER TABLE `transaction` ADD FOREIGN KEY (`customer_id`) REFERENCES `customers` (`id`);
ALTER TABLE `loan_limit` ADD CONSTRAINT unique_customer_id_tenor UNIQUE (`customer_id`, `tenor`);
