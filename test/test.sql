-- Используем тестовую базу
CREATE
DATABASE IF NOT EXISTS testdb;
USE
testdb;

-- Таблица пользователей
CREATE TABLE users
(
    id         INT AUTO_INCREMENT PRIMARY KEY,
    name       VARCHAR(100),
    email      VARCHAR(100) UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO users (name, email)
VALUES ('Иван Иванов', 'ivan@example.com'),
       ('Петр Петров', 'petr@example.com'),
       ('Анна Смирнова', 'anna@example.com');

-- Таблица заказов
CREATE TABLE orders
(
    id         INT AUTO_INCREMENT PRIMARY KEY,
    user_id    INT,
    total      DECIMAL(10, 2),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users (id)
);

INSERT INTO orders (user_id, total)
VALUES (1, 100.50),
       (2, 200.75),
       (3, 150.00);

-- Таблица продуктов
CREATE TABLE products
(
    id    INT AUTO_INCREMENT PRIMARY KEY,
    name  VARCHAR(100),
    price DECIMAL(10, 2),
    stock INT
);

INSERT INTO products (name, price, stock)
VALUES ('Ноутбук', 50000.00, 10),
       ('Смартфон', 30000.00, 20),
       ('Наушники', 5000.00, 50);

-- Таблица платежей
CREATE TABLE payments
(
    id         INT AUTO_INCREMENT PRIMARY KEY,
    order_id   INT,
    amount     DECIMAL(10, 2),
    status     ENUM('pending', 'completed', 'failed') DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (order_id) REFERENCES orders (id)
);

INSERT INTO payments (order_id, amount, status)
VALUES (1, 100.50, 'completed'),
       (2, 200.75, 'pending'),
       (3, 150.00, 'failed');

-- Таблица адресов доставки
CREATE TABLE addresses
(
    id         INT AUTO_INCREMENT PRIMARY KEY,
    user_id    INT,
    address    VARCHAR(255),
    city       VARCHAR(100),
    country    VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users (id)
);

INSERT INTO addresses (user_id, address, city, country)
VALUES (1, 'ул. Ленина, д.10', 'Москва', 'Россия'),
       (2, 'пр. Независимости, д.20', 'Минск', 'Беларусь'),
       (3, 'ул. Шевченко, д.5', 'Киев', 'Украина');

-- Таблица категорий товаров
CREATE TABLE categories
(
    id   INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100)
);

INSERT INTO categories (name)
VALUES ('Электроника'),
       ('Одежда'),
       ('Бытовая техника');

-- Таблица товаров в категориях (многие ко многим)
CREATE TABLE product_categories
(
    product_id  INT,
    category_id INT,
    PRIMARY KEY (product_id, category_id),
    FOREIGN KEY (product_id) REFERENCES products (id),
    FOREIGN KEY (category_id) REFERENCES categories (id)
);

INSERT INTO product_categories (product_id, category_id)
VALUES (1, 1),
       (2, 1),
       (3, 1);

-- Таблица отзывов
CREATE TABLE reviews
(
    id         INT AUTO_INCREMENT PRIMARY KEY,
    product_id INT,
    user_id    INT,
    rating     INT CHECK (rating BETWEEN 1 AND 5),
    comment    TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (product_id) REFERENCES products (id),
    FOREIGN KEY (user_id) REFERENCES users (id)
);

INSERT INTO reviews (product_id, user_id, rating, comment)
VALUES (1, 1, 5, 'Отличный ноутбук!'),
       (2, 2, 4, 'Хороший смартфон, но дорогой.'),
       (3, 3, 3, 'Наушники норм, но звук не топ.');

-- Таблица скидок
CREATE TABLE discounts
(
    id               INT AUTO_INCREMENT PRIMARY KEY,
    product_id       INT,
    discount_percent INT CHECK (discount_percent BETWEEN 1 AND 100),
    active           BOOLEAN DEFAULT TRUE,
    FOREIGN KEY (product_id) REFERENCES products (id)
);

INSERT INTO discounts (product_id, discount_percent, active)
VALUES (1, 10, TRUE),
       (2, 15, FALSE),
       (3, 5, TRUE);

-- Таблица логов действий пользователей
CREATE TABLE user_logs
(
    id        INT AUTO_INCREMENT PRIMARY KEY,
    user_id   INT,
    action    VARCHAR(255),
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users (id)
);

INSERT INTO user_logs (user_id, action)
VALUES (1, 'Авторизация'),
       (2, 'Добавление товара в корзину'),
       (3, 'Оформление заказа');

-- Таблица транзакций
CREATE TABLE transactions
(
    id         INT AUTO_INCREMENT PRIMARY KEY,
    user_id    INT,
    amount     DECIMAL(10, 2),
    type       ENUM('deposit', 'withdrawal'),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users (id)
);

INSERT INTO transactions (user_id, amount, type)
VALUES (1, 1000.00, 'deposit'),
       (2, 500.00, 'withdrawal'),
       (3, 1500.00, 'deposit');

-- Таблица складов
CREATE TABLE warehouses
(
    id       INT AUTO_INCREMENT PRIMARY KEY,
    name     VARCHAR(100),
    location VARCHAR(255)
);

INSERT INTO warehouses (name, location)
VALUES ('Склад 1', 'Москва'),
       ('Склад 2', 'СПб');

-- Таблица логистики (склад → товар)
CREATE TABLE warehouse_products
(
    warehouse_id INT,
    product_id   INT,
    quantity     INT,
    PRIMARY KEY (warehouse_id, product_id),
    FOREIGN KEY (warehouse_id) REFERENCES warehouses (id),
    FOREIGN KEY (product_id) REFERENCES products (id)
);

INSERT INTO warehouse_products (warehouse_id, product_id, quantity)
VALUES (1, 1, 5),
       (2, 2, 10);

-- Таблица промокодов
CREATE TABLE promo_codes
(
    id               INT AUTO_INCREMENT PRIMARY KEY,
    code             VARCHAR(50) UNIQUE,
    discount_percent INT CHECK (discount_percent BETWEEN 1 AND 50),
    valid_until      DATE
);

INSERT INTO promo_codes (code, discount_percent, valid_until)
VALUES ('NEWYEAR2025', 20, '2025-01-01'),
       ('SUMMER2025', 15, '2025-06-01');

-- Таблица подписок пользователей
CREATE TABLE subscriptions
(
    id         INT AUTO_INCREMENT PRIMARY KEY,
    user_id    INT,
    type       ENUM('monthly', 'yearly'),
    expires_at DATE,
    FOREIGN KEY (user_id) REFERENCES users (id)
);

INSERT INTO subscriptions (user_id, type, expires_at)
VALUES (1, 'monthly', '2025-03-01'),
       (2, 'yearly', '2026-01-01');

-- Таблица логов ошибок
CREATE TABLE error_logs
(
    id        INT AUTO_INCREMENT PRIMARY KEY,
    message   TEXT,
    level     ENUM('info', 'warning', 'error'),
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO error_logs (message, level)
VALUES ('Ошибка в системе оплаты', 'error'),
       ('Предупреждение: низкий уровень запаса товара', 'warning');

-- Последняя тестовая таблица
CREATE TABLE test_data
(
    id    INT AUTO_INCREMENT PRIMARY KEY,
    value TEXT
);

INSERT INTO test_data (value)
VALUES ('Тест 1'),
       ('Тест 2');