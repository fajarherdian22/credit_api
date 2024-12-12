DROP TABLE IF EXISTS `sessions`;

ALTER TABLE payment_details DROP FOREIGN KEY payment_details_ibfk_1;

DROP TABLE IF EXISTS payment_details;