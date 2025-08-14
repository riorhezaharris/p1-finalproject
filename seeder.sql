-- 1. Populate `users` table (5 users)
-- Passwords should be hashed in a real application. Using plain text for example purposes.
INSERT INTO `users` (`email`, `password`) VALUES
('alice.johnson@example.com', 'password1'),
('bob.williams@example.com', 'password2'),
('charlie.brown@example.com', 'password3'),
('diana.miller@example.com', 'password4'),
('edward.jones@example.com', 'password5');


-- 2. Populate `sizes` table (5 sizes)
INSERT INTO `sizes` (`name`, `width`, `length`) VALUES
('S', 46, 66),
('M', 51, 71),
('L', 56, 74),
('XL', 61, 76),
('XXL', 66, 81);


-- 3. Populate `suppliers` table (5 suppliers)
INSERT INTO `suppliers` (`name`, `location`) VALUES
('Global Textiles Inc.', 'Bandung, Indonesia'),
('Apex Apparel', 'Surabaya, Indonesia'),
('Vertex Garments', 'Ho Chi Minh City, Vietnam'),
('ProStitch Manufacturing', 'Bangkok, Thailand'),
('Quality Threads Co.', 'Dhaka, Bangladesh');


-- 4. Populate `products` table (30 products)
-- Products are randomly assigned a size and supplier.
INSERT INTO `products` (`name`, `price`, `size_id`, `supplier_id`) VALUES
('Classic Crew Neck T-Shirt', 15.99, 2, 1),
('V-Neck T-Shirt', 16.50, 3, 2),
('Long Sleeve Henley', 25.00, 4, 3),
('Graphic Print Tee', 22.99, 1, 4),
('Pocket T-Shirt', 18.00, 5, 5),
('Performance Polo Shirt', 35.50, 2, 1),
('Slim Fit Chinos', 49.99, 3, 2),
('Denim Jeans', 65.00, 4, 3),
('Cargo Shorts', 39.99, 1, 4),
('Fleece Hoodie', 55.00, 5, 5),
('Zip-Up Sweatshirt', 52.50, 2, 1),
('Lightweight Bomber Jacket', 75.00, 3, 2),
('Quilted Puffer Vest', 68.99, 4, 3),
('Wool Blend Overcoat', 150.00, 1, 4),
('Linen Button-Down Shirt', 45.00, 5, 5),
('Flannel Plaid Shirt', 48.50, 2, 1),
('Knit Sweater', 60.00, 3, 2),
('Turtleneck Pullover', 58.00, 4, 3),
('Jogger Sweatpants', 42.99, 1, 4),
('Athletic Running Shorts', 29.99, 5, 5),
('Cotton Boxer Briefs (3-Pack)', 28.00, 2, 1),
('Merino Wool Socks', 12.50, 3, 2),
('Leather Belt', 38.00, 4, 3),
('Canvas Web Belt', 20.00, 1, 4),
('Baseball Cap', 24.99, 5, 5),
('Beanie Hat', 19.99, 2, 1),
('Scarf', 26.00, 3, 2),
('Leather Gloves', 45.00, 4, 3),
('Waterproof Rain Jacket', 89.99, 1, 4),
('Tailored Suit Jacket', 250.00, 5, 5);