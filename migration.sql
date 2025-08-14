CREATE TABLE IF NOT EXISTS `users` (
  `id` int PRIMARY KEY NOT NULL AUTO_INCREMENT,
  `email` varchar(255),
  `password` varchar(255)
);

CREATE TABLE IF NOT EXISTS `sizes` (
  `id` int PRIMARY KEY NOT NULL AUTO_INCREMENT,
  `name` varchar(255),
  `width` int,
  `length` int
);

CREATE TABLE IF NOT EXISTS `suppliers` (
  `id` int PRIMARY KEY NOT NULL AUTO_INCREMENT,
  `name` varchar(255),
  `location` varchar(255)
);

CREATE TABLE IF NOT EXISTS `products` (
  `id` int PRIMARY KEY NOT NULL AUTO_INCREMENT,
  `name` varchar(255),
  `price` float,
  `size_id` int,
  `supplier_id` int
);

CREATE TABLE IF NOT EXISTS `carts` (
  `id` int PRIMARY KEY NOT NULL AUTO_INCREMENT,
  `user_id` int,
  `product_id` int,
  `quantity` int
);

CREATE TABLE IF NOT EXISTS `orders` (
  `id` int PRIMARY KEY NOT NULL AUTO_INCREMENT,
  `user_id` int,
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  `total_price` float,
  `status` varchar(255)
);

CREATE TABLE IF NOT EXISTS `order_details` (
  `id` int PRIMARY KEY NOT NULL AUTO_INCREMENT,
  `product_id` int,
  `quantity` int,
  `price` float,
  `order_id` int
);

CREATE TABLE IF NOT EXISTS `payments` (
  `id` int PRIMARY KEY NOT NULL AUTO_INCREMENT,
  `order_id` int,
  `total_payment` int,
  `status` varchar(255)
);

ALTER TABLE `products` ADD FOREIGN KEY (`size_id`) REFERENCES `sizes` (`id`);

ALTER TABLE `products` ADD FOREIGN KEY (`supplier_id`) REFERENCES `suppliers` (`id`);

ALTER TABLE `carts` ADD FOREIGN KEY (`product_id`) REFERENCES `products` (`id`);

ALTER TABLE `carts` ADD FOREIGN KEY (`user_id`) REFERENCES `users` (`id`);

ALTER TABLE `orders` ADD FOREIGN KEY (`user_id`) REFERENCES `users` (`id`);

ALTER TABLE `order_details` ADD FOREIGN KEY (`order_id`) REFERENCES `orders` (`id`);

ALTER TABLE `order_details` ADD FOREIGN KEY (`product_id`) REFERENCES `products` (`id`);

ALTER TABLE `payments` ADD FOREIGN KEY (`order_id`) REFERENCES `orders` (`id`);
