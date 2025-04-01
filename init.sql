-- 1. ENUM Types
CREATE TYPE order_status AS ENUM ('pending', 'preparing', 'completed', 'canceled');
CREATE TYPE payment_method AS ENUM ('cash', 'card', 'online');
CREATE TYPE staff_role AS ENUM ('barista', 'cashier', 'manager');
CREATE TYPE item_size AS ENUM ('small', 'medium', 'large');
CREATE TYPE unit_type AS ENUM ('grams', 'ml', 'pcs');

-- 2. Customers Table
CREATE TABLE customers (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    preferences JSONB
);

-- 3. Orders Table
CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    customer_id INTEGER REFERENCES customers(id),
    status order_status DEFAULT 'pending',
    special_instructions JSONB,
    total_amount NUMERIC(10,2) DEFAULT 0,
    order_date TIMESTAMPTZ DEFAULT NOW()
);

-- 4. Order Status History
CREATE TABLE order_status_history (
    id SERIAL PRIMARY KEY,
    order_id INTEGER REFERENCES orders(id) ON DELETE CASCADE,
    status order_status,
    changed_at TIMESTAMPTZ DEFAULT NOW()
);

-- 5. Menu Items
CREATE TABLE menu_items (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT,
    price NUMERIC(10,2) NOT NULL,
    category TEXT[],
    allergens TEXT[],
    customization_options JSONB,
    size item_size,
    metadata JSONB,
    UNIQUE (name, description, price, size)
);


-- 6. Order Items
CREATE TABLE order_items (
    id SERIAL PRIMARY KEY,
    order_id INTEGER REFERENCES orders(id) ON DELETE CASCADE,
    menu_item_id INTEGER REFERENCES menu_items(id),
    quantity INTEGER NOT NULL,
    price_at_order_time NUMERIC(10,2) NOT NULL,
    customization JSONB
);

-- 7. Inventory
CREATE TABLE inventory (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    quantity INTEGER NOT NULL,
    unit unit_type,
    price_per_unit NUMERIC(10,2),
    last_updated TIMESTAMPTZ DEFAULT NOW()
);

-- 8. Menu Item Ingredients (Junction)
CREATE TABLE menu_item_ingredients (
    id SERIAL PRIMARY KEY,
    menu_item_id INTEGER REFERENCES menu_items(id) ON DELETE CASCADE,
    ingredient_id INTEGER REFERENCES inventory(id),
    quantity_required INTEGER NOT NULL
);

-- 9. Price History
CREATE TABLE price_history (
    id SERIAL PRIMARY KEY,
    menu_item_id INTEGER REFERENCES menu_items(id) ON DELETE CASCADE,
    price NUMERIC(10,2) NOT NULL,
    changed_at TIMESTAMPTZ DEFAULT NOW()
);

-- 10. Inventory Transactions
CREATE TABLE inventory_transactions (
    id SERIAL PRIMARY KEY,
    inventory_id INTEGER REFERENCES inventory(id) ON DELETE CASCADE,
    change_amount INTEGER NOT NULL,
    transaction_date TIMESTAMPTZ DEFAULT NOW(),
    reason TEXT
);

-- 11. Indexes
CREATE INDEX idx_orders_customer_id ON orders(customer_id);
CREATE INDEX idx_order_items_order_id ON order_items(order_id);
CREATE INDEX idx_menu_items_search ON menu_items USING gin (to_tsvector('english', name || ' ' || description));
CREATE INDEX idx_inventory_name ON inventory(name);

CREATE OR REPLACE FUNCTION update_order_total()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE orders
    SET total_amount = (
        SELECT COALESCE(SUM(quantity * price_at_order_time), 0)
        FROM order_items
        WHERE order_id = NEW.order_id
    )
    WHERE id = NEW.order_id;

    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

-- Срабатывает при вставке
CREATE TRIGGER recalculate_total_after_insert
AFTER INSERT ON order_items
FOR EACH ROW
EXECUTE FUNCTION update_order_total();

-- Срабатывает при обновлении
CREATE TRIGGER recalculate_total_after_update
AFTER UPDATE ON order_items
FOR EACH ROW
EXECUTE FUNCTION update_order_total();

-- Срабатывает при удалении
CREATE TRIGGER recalculate_total_after_delete
AFTER DELETE ON order_items
FOR EACH ROW
EXECUTE FUNCTION update_order_total();


-- 12. Mock Data

-- Customers
-- 1. ENUM Types
CREATE TYPE order_status AS ENUM ('pending', 'preparing', 'completed', 'canceled');
CREATE TYPE payment_method AS ENUM ('cash', 'card', 'online');
CREATE TYPE staff_role AS ENUM ('barista', 'cashier', 'manager');
CREATE TYPE item_size AS ENUM ('small', 'medium', 'large');
CREATE TYPE unit_type AS ENUM ('grams', 'ml', 'pcs');

-- 2. Customers Table
CREATE TABLE customers (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    preferences JSONB
);

-- 3. Orders Table
CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    customer_id INTEGER REFERENCES customers(id),
    status order_status DEFAULT 'pending',
    special_instructions JSONB,
    total_amount NUMERIC(10,2) DEFAULT 0,
    order_date TIMESTAMPTZ DEFAULT NOW()
);

-- 4. Order Status History
CREATE TABLE order_status_history (
    id SERIAL PRIMARY KEY,
    order_id INTEGER REFERENCES orders(id) ON DELETE CASCADE,
    status order_status,
    changed_at TIMESTAMPTZ DEFAULT NOW()
);

-- 5. Menu Items
CREATE TABLE menu_items (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT,
    price NUMERIC(10,2) NOT NULL,
    category TEXT[],
    allergens TEXT[],
    customization_options JSONB,
    size item_size,
    metadata JSONB,
    UNIQUE (name, description, price, size)
);


-- 6. Order Items
CREATE TABLE order_items (
    id SERIAL PRIMARY KEY,
    order_id INTEGER REFERENCES orders(id) ON DELETE CASCADE,
    menu_item_id INTEGER REFERENCES menu_items(id),
    quantity INTEGER NOT NULL,
    price_at_order_time NUMERIC(10,2) NOT NULL,
    customization JSONB
);

-- 7. Inventory
CREATE TABLE inventory (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    quantity INTEGER NOT NULL,
    unit unit_type,
    price_per_unit NUMERIC(10,2),
    last_updated TIMESTAMPTZ DEFAULT NOW()
);

-- 8. Menu Item Ingredients (Junction)
CREATE TABLE menu_item_ingredients (
    id SERIAL PRIMARY KEY,
    menu_item_id INTEGER REFERENCES menu_items(id) ON DELETE CASCADE,
    ingredient_id INTEGER REFERENCES inventory(id),
    quantity_required INTEGER NOT NULL
);

-- 9. Price History
CREATE TABLE price_history (
    id SERIAL PRIMARY KEY,
    menu_item_id INTEGER REFERENCES menu_items(id) ON DELETE CASCADE,
    price NUMERIC(10,2) NOT NULL,
    changed_at TIMESTAMPTZ DEFAULT NOW()
);

-- 10. Inventory Transactions
CREATE TABLE inventory_transactions (
    id SERIAL PRIMARY KEY,
    inventory_id INTEGER REFERENCES inventory(id) ON DELETE CASCADE,
    change_amount INTEGER NOT NULL,
    transaction_date TIMESTAMPTZ DEFAULT NOW(),
    reason TEXT
);

-- 11. Indexes
CREATE INDEX idx_orders_customer_id ON orders(customer_id);
CREATE INDEX idx_order_items_order_id ON order_items(order_id);
CREATE INDEX idx_menu_items_search ON menu_items USING gin (to_tsvector('english', name || ' ' || description));
CREATE INDEX idx_inventory_name ON inventory(name);

CREATE OR REPLACE FUNCTION update_order_total()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE orders
    SET total_amount = (
        SELECT COALESCE(SUM(quantity * price_at_order_time), 0)
        FROM order_items
        WHERE order_id = NEW.order_id
    )
    WHERE id = NEW.order_id;

    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

-- Срабатывает при вставке
CREATE TRIGGER recalculate_total_after_insert
AFTER INSERT ON order_items
FOR EACH ROW
EXECUTE FUNCTION update_order_total();

-- Срабатывает при обновлении
CREATE TRIGGER recalculate_total_after_update
AFTER UPDATE ON order_items
FOR EACH ROW
EXECUTE FUNCTION update_order_total();

-- Срабатывает при удалении
CREATE TRIGGER recalculate_total_after_delete
AFTER DELETE ON order_items
FOR EACH ROW
EXECUTE FUNCTION update_order_total();


-- 12. Mock Data

-- Customers
 INSERT INTO customers (name, preferences) VALUES
('Alice Brown', '{"favorite_drink": "Latte", "no_sugar": true}'),
('Bob Smith', '{"allergy": "nuts"}'),
('Charlie Green', '{}'),
('Charlie Put','{}'),
('John Lenon','{}'),
('Michael Joardan','{}'),
('Michael Jackson','{}'),
('Conor McGregor','{}'),
('Justin Gatche','{}'),
('Muhammad Ali','{}'),
('Islam Mackhachev','{}'),
('Ibragim Betaev','{}'),
('Danil Stenkin','{}');

-- Menu Items
INSERT INTO menu_items (name, description, price, category, allergens, customization_options, size, metadata) VALUES
('Latte', 'Classic milk coffee', 4.50, ARRAY['coffee', 'hot'], ARRAY['milk'], '{"syrup": "vanilla"}', 'medium', '{"season": "winter"}'),
('Espresso', 'Strong black coffee', 3.00, ARRAY['coffee'], ARRAY[]::TEXT[], '{}', 'small', '{}'),
('Muffin', 'Chocolate muffin', 2.00, ARRAY['dessert'], ARRAY['gluten', 'eggs'], '{}', NULL, '{}'),
('Cappuccino', 'Espresso with steamed milk and foam', 4.00, ARRAY['coffee', 'hot'], ARRAY['milk'], '{"cocoa_powder": "optional"}', 'medium', '{"season": "all"}'),
('Iced Americano', 'Chilled espresso with water', 3.50, ARRAY['coffee', 'cold'], ARRAY[]::TEXT[], '{}', 'large', '{"season": "summer"}'),
('Green Tea', 'Hot green tea', 2.50, ARRAY['tea', 'hot'], ARRAY[]::TEXT[], '{}', 'medium', '{}'),
('Croissant', 'Buttery flaky pastry', 2.20, ARRAY['pastry', 'breakfast'], ARRAY['gluten', 'eggs', 'milk'], '{}', NULL, '{"origin": "France"}'),
('Orange Juice', 'Freshly squeezed orange juice', 3.00, ARRAY['drink', 'cold'], ARRAY[]::TEXT[], '{}', 'medium', '{"vitaminC": "high"}'),
('Bagel with Cream Cheese', 'Toasted bagel with cream cheese', 3.20, ARRAY['breakfast', 'pastry'], ARRAY['gluten', 'milk'], '{"toasting_level": "light"}', NULL, '{}'),
('Mocha', 'Chocolate-flavored coffee drink', 4.80, ARRAY['coffee', 'hot'], ARRAY['milk'], '{"whipped_cream": "optional"}', 'medium', '{"season": "winter"}'),
('Chai Latte', 'Spiced tea with steamed milk', 4.30, ARRAY['tea', 'hot'], ARRAY['milk'], '{"sweetness": "medium"}', 'medium', '{}'),
('Vegan Brownie', 'Rich chocolate brownie', 2.80, ARRAY['dessert', 'vegan'], ARRAY['gluten'], '{}', NULL, '{"calories": "250"}'),
('Matcha Latte', 'Green tea latte', 4.70, ARRAY['tea', 'hot'], ARRAY['milk'], '{"milk_type": "oat"}', 'medium', '{"trend": "popular"}');

-- Inventory
INSERT INTO inventory (name, quantity, unit, price_per_unit) VALUES
('Coffee Beans', 10000, 'grams', 0.05),
('Milk', 5000, 'ml', 0.03),
('Chocolate', 2000, 'grams', 0.10),
('Flour', 3000, 'grams', 0.02),
('Eggs', 200, 'pcs', 0.15),
('Sugar', 4000, 'grams', 0.01),
('Butter', 1500, 'grams', 0.08),
('Vanilla Syrup', 1000, 'ml', 0.12),
('Cocoa Powder', 800, 'grams', 0.09),
('Tea Leaves', 1200, 'grams', 0.04),
('Orange Juice', 3000, 'ml', 0.05),
('Cream Cheese', 1000, 'grams', 0.11),
('Bagels', 100, 'pcs', 0.40),
('Green Tea Powder (Matcha)', 500, 'grams', 0.20),
('Oat Milk', 2000, 'ml', 0.07),
('Whipped Cream', 700, 'grams', 0.10),
('Baking Powder', 300, 'grams', 0.03),
('Yeast', 200, 'grams', 0.02),
('Salt', 1000, 'grams', 0.005),
('Lemon Juice', 500, 'ml', 0.06);

-- Menu Item Ingredients
INSERT INTO menu_item_ingredients (menu_item_id, ingredient_id, quantity_required) VALUES
(1, 1, 100),
(1, 2, 200),
(2, 1, 80),
(3, 4, 150),
(3, 5, 2),
(3, 3, 50),
(4, 1, 80),
(4, 2, 150),
(4, 9, 10),

(5, 1, 80),
(5, 2, 50),
(5, 6, 100),

(6, 5, 0),
(6, 10, 5),

(7, 4, 100),
(7, 5, 1),
(7, 6, 20),
(7, 7, 50),

(8, 11, 200),

(9, 8, 1),
(9, 12, 50),

(10, 1, 90),
(10, 2, 150),
(10, 3, 30),
(10, 13, 20),

(11, 10, 100),
(11, 2, 150),

(12, 3, 50),
(12, 4, 120),
(12, 6, 40),
(12, 7, 30),

(13, 14, 5),
(13, 15, 150);

-- Orders
INSERT INTO orders (customer_id, status, special_instructions, total_amount, order_date) VALUES
(1, 'completed', '{"extra_shot": true}', 9.50, NOW() - INTERVAL '2 days'),
(2, 'preparing', '{}', 5.00, NOW()),
(3, 'pending', '{"no_milk": true}', 3.00, NOW()),
(4, 'completed', '{"add_sugar": true}', 7.20, NOW() - INTERVAL '3 days'),
(2, 'completed', '{}', 11.50, NOW() - INTERVAL '1 day'),
(5, 'cancelled', '{"reason": "customer_request"}', 4.00, NOW() - INTERVAL '5 days'),
(1, 'preparing', '{"less_foam": true}', 6.75, NOW()),
(3, 'completed', '{"no_chocolate": true}', 8.10, NOW() - INTERVAL '4 days'),
(6, 'pending', '{}', 5.50, NOW()),
(4, 'completed', '{"oat_milk": true}', 9.00, NOW() - INTERVAL '7 days'),
(7, 'preparing', '{"no_sugar": true}', 4.80, NOW());

-- Order Items
INSERT INTO order_items (order_id, menu_item_id, quantity, price_at_order_time, customization) VALUES
(1, 1, 2, 4.50, '{"syrup": "caramel"}'),
(2, 2, 1, 3.00, '{}'),
(3, 2, 1, 3.00, '{}'),
(4, 4, 1, 4.00, '{"cocoa_powder": "optional"}'),
(4, 7, 1, 2.20, '{}'),
(5, 3, 1, 2.00, '{}'),
(6, 1, 1, 4.50, '{"syrup": "vanilla"}'),
(6, 10, 1, 4.80, '{"whipped_cream": "yes"}'),
(7, 2, 1, 3.00, '{}'),
(8, 5, 1, 3.50, '{}'),
(9, 13, 1, 4.70, '{"milk_type": "oat"}'),
(10, 11, 1, 3.00, '{}');

-- Order Status History
INSERT INTO order_status_history (order_id, status, changed_at) VALUES
(1, 'pending', NOW() - INTERVAL '3 days'),
(1, 'completed', NOW() - INTERVAL '2 days'),
(2, 'pending', NOW() - INTERVAL '1 day'),
(2, 'preparing', NOW()),
(3, 'pending', NOW() - INTERVAL '12 hours'),
(4, 'pending', NOW() - INTERVAL '3 days'),
(4, 'completed', NOW() - INTERVAL '3 days' + INTERVAL '2 hours'),
(5, 'pending', NOW() - INTERVAL '5 days'),
(5, 'cancelled', NOW() - INTERVAL '5 days' + INTERVAL '1 hour'),
(6, 'pending', NOW() - INTERVAL '6 hours'),
(7, 'pending', NOW() - INTERVAL '8 days'),
(7, 'completed', NOW() - INTERVAL '7 days'),
(8, 'pending', NOW() - INTERVAL '1 day'),
(8, 'preparing', NOW()),
(9, 'pending', NOW() - INTERVAL '4 days'),
(9, 'completed', NOW() - INTERVAL '4 days' + INTERVAL '3 hours'),
(10, 'pending', NOW() - INTERVAL '3 hours'),
(10, 'preparing', NOW());

-- Price History
INSERT INTO price_history (menu_item_id, price, changed_at) VALUES
(1, 4.00, NOW() - INTERVAL '6 months'),
(1, 4.50, NOW() - INTERVAL '1 month'),
(2, 3.00, NOW() - INTERVAL '3 months'),
(3, 1.80, NOW() - INTERVAL '4 months'),
(3, 2.00, NOW() - INTERVAL '1 month'),
(4, 3.80, NOW() - INTERVAL '5 months'),
(4, 4.00, NOW() - INTERVAL '2 months'),
(5, 3.20, NOW() - INTERVAL '6 months'),
(5, 3.50, NOW() - INTERVAL '1 month'),
(6, 2.30, NOW() - INTERVAL '3 months'),
(6, 2.50, NOW() - INTERVAL '2 weeks'),
(7, 2.00, NOW() - INTERVAL '6 months'),
(7, 2.20, NOW() - INTERVAL '1 month'),
(10, 4.50, NOW() - INTERVAL '4 months'),
(10, 4.80, NOW() - INTERVAL '2 weeks');

-- Inventory Transactions
INSERT INTO inventory_transactions (inventory_id, change_amount, transaction_date, reason) VALUES
(1, -200, NOW() - INTERVAL '1 day', 'Order #1'),
(2, -200, NOW() - INTERVAL '1 day', 'Order #1'),
(4, -150, NOW() - INTERVAL '2 days', 'Order #3'),
(1, -80, NOW() - INTERVAL '3 days', 'Order #4'),
(2, -150, NOW() - INTERVAL '3 days', 'Order #4'),
(9, -10, NOW() - INTERVAL '3 days', 'Order #4'),

(3, -50, NOW() - INTERVAL '2 days', 'Order #3'),

(8, -1, NOW() - INTERVAL '1 day', 'Order #9'),
(12, -50, NOW() - INTERVAL '1 day', 'Order #9'),

(14, -5, NOW() - INTERVAL '4 hours', 'Order #13'),
(15, -150, NOW() - INTERVAL '4 hours', 'Order #13'),

(10, -100, NOW() - INTERVAL '2 days', 'Order #11'),
(2, -150, NOW() - INTERVAL '2 days', 'Order #11');