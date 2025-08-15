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
('Classic Crew Neck T-Shirt', 160000, 2, 1),
('V-Neck T-Shirt', 165000, 3, 2),
('Long Sleeve Henley', 250000, 4, 3),
('Graphic Print Tee', 230000, 1, 4),
('Pocket T-Shirt', 180000, 5, 5),
('Performance Polo Shirt', 355000, 2, 1),
('Slim Fit Chinos', 499000, 3, 2),
('Denim Jeans', 650000, 4, 3),
('Cargo Shorts', 399000, 1, 4),
('Fleece Hoodie', 550000, 5, 5),
('Zip-Up Sweatshirt', 525000, 2, 1),
('Lightweight Bomber Jacket', 750000, 3, 2),
('Quilted Puffer Vest', 689900, 4, 3),
('Wool Blend Overcoat', 150000, 1, 4),
('Linen Button-Down Shirt', 450000, 5, 5),
('Flannel Plaid Shirt', 485000, 2, 1),
('Knit Sweater', 600000, 3, 2),
('Turtleneck Pullover', 580000, 4, 3),
('Jogger Sweatpants', 429900, 1, 4),
('Athletic Running Shorts', 299900, 5, 5),
('Cotton Boxer Briefs (3-Pack)', 280000, 2, 1),
('Merino Wool Socks', 125000, 3, 2),
('Leather Belt', 380000, 4, 3),
('Canvas Web Belt', 200000, 1, 4),
('Baseball Cap', 249900, 5, 5),
('Beanie Hat', 199900, 2, 1),
('Scarf', 260000, 3, 2),
('Leather Gloves', 450000, 4, 3),
('Waterproof Rain Jacket', 899900, 1, 4),
('Tailored Suit Jacket', 250000, 5, 5);