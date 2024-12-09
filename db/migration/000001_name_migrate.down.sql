-- Drop foreign keys
ALTER TABLE loan_limit DROP FOREIGN KEY loan_limit_ibfk_1;

ALTER TABLE `transaction` DROP FOREIGN KEY transaction_ibfk_1;

ALTER TABLE payment_details DROP FOREIGN KEY payment_details_ibfk_1;

DROP TABLE IF EXISTS customers;

DROP TABLE IF EXISTS loan_limit;

DROP TABLE IF EXISTS `transaction`;

DROP TABLE IF EXISTS payment_details;