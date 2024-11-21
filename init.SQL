CREATE DATABASE IF NOT EXISTS ecommerce;
use ecommerce;

CREATE TABLE IF NOT EXISTS users (
    id VARCHAR(36) PRIMARY KEY,  
    name VARCHAR(255) NOT NULL,  
    lastname VARCHAR(255) NOT NULL,  
    email VARCHAR(255) UNIQUE NOT NULL, 
    password VARCHAR(255) NOT NULL, 
    created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    role VARCHAR(255) NOT NULL 
);
CREATE TABLE IF NOT EXISTS accessories(
    id VARCHAR(36) PRIMARY KEY,  
    name VARCHAR(255) NOT NULL,
    description VARCHAR(255) NOT NULL,   
    quantity INT NOT NULL,
    price FLOAT NOT NULL,
    category VARCHAR(255),
    image VARCHAR(255)
);
CREATE TABLE IF NOT EXISTS bicycles (
    id VARCHAR(36) PRIMARY KEY,  
    name VARCHAR(255) NOT NULL,
    model VARCHAR(255) NOT NULL,
    description VARCHAR(255) NOT NULL,
    type VARCHAR(255) NOT NULL,
    size VARCHAR(255) NOT NULL,   
    material VARCHAR(255) NOT NULL,
    quantity INT NOT NULL,
    price FLOAT NOT NULL,
    image VARCHAR(255),
    color VARCHAR(50) NOT NULL,
    weight FLOAT,
    release_year INT,
    brake_system VARCHAR(100),
    gears INT,
    brand VARCHAR(255) NOT NULL,
    suspension VARCHAR(100),
    wheel_size INT,
    frame_size VARCHAR(50)
);


CREATE TABLE IF NOT EXISTS cart(
	user_id VARCHAR(255) NOT NULL,
    product_id VARCHAR(255) NOT NULL,
    quantity VARCHAR(255) NOT NULL,
    created TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS orders(
    order_id VARCHAR(255) NOT NULL,
	user_id VARCHAR(255) NOT NULL,
    status VARCHAR(255),
    total_price INT NOT NULL,
    created TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS orders_products(
    user_id VARCHAR(255) NOT NULL,
    order_id VARCHAR(255) NOT NULL,
    product_id VARCHAR(255) NOT NULL,
    quantity int NOT NULL,
    total_price FLOAT NOT NULL
);

CREATE TABLE IF NOT EXISTS  order_details (
    user_id VARCHAR(255) NOT NULL,
    order_id VARCHAR(255) NOT NULL,
    first_name   VARCHAR(255) NOT NULL,
	last_name    VARCHAR(255) NOT NULL,
	email       VARCHAR(255) NOT NULL,
	phone_number VARCHAR(255) NOT NULL,
	country     VARCHAR(255) NOT NULL,
	city        VARCHAR(255) NOT NULL,
	street      VARCHAR(255) NOT NULL,
	house       VARCHAR(255) NOT NULL
    );

CREATE TABLE IF NOT EXISTS reviews();
