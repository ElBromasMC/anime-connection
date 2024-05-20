CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pg_trgm";
CREATE EXTENSION IF NOT EXISTS "hstore";

-- Row update management
CREATE OR REPLACE FUNCTION trigger_set_timestamp()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- User administration
CREATE TYPE user_role AS ENUM ('ADMIN', 'NORMAL', 'RECORDER');

CREATE TABLE IF NOT EXISTS users (
    user_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    hashed_password TEXT NOT NULL,
    role user_role NOT NULL DEFAULT 'NORMAL',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS sessions (
    session_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMPTZ NOT NULL DEFAULT NOW() + INTERVAL '1 month',
    FOREIGN KEY (user_id) REFERENCES users(user_id)
);

-- Image administrarion
CREATE TABLE IF NOT EXISTS images (
    id SERIAL PRIMARY KEY,
    filename VARCHAR(25) UNIQUE NOT NULL
);

-- Store administration
CREATE TYPE store_type AS ENUM ('STORE', 'GARANTIA');

CREATE TABLE IF NOT EXISTS store_categories (
    id SERIAL PRIMARY KEY,
    type store_type NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    img_id INT,
    slug VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(type, slug),
    FOREIGN KEY (img_id) REFERENCES images(id) ON DELETE SET NULL
);

CREATE TABLE IF NOT EXISTS store_items (
    id SERIAL PRIMARY KEY,
    category_id INT NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    long_description TEXT NOT NULL DEFAULT '',
    img_id INT,
    largeimg_id INT,
    slug VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(category_id, slug),
    FOREIGN KEY (category_id) REFERENCES store_categories(id) ON DELETE CASCADE,
    FOREIGN KEY (img_id) REFERENCES images(id) ON DELETE SET NULL,
    FOREIGN KEY (largeimg_id) REFERENCES images(id) ON DELETE SET NULL
);
CREATE INDEX idx_items_name ON store_items USING gin (name gin_trgm_ops);

-- Product management
CREATE TABLE IF NOT EXISTS store_products (
    id SERIAL PRIMARY KEY,
    item_id INT NOT NULL,
    name VARCHAR(255) NOT NULL,
    price INT NOT NULL,
    stock INT,
    details HSTORE NOT NULL DEFAULT ''::hstore,
    slug VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(item_id, slug),
    FOREIGN KEY (item_id) REFERENCES store_items(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS product_discount (
    id SERIAL PRIMARY KEY,
    product_id INT NOT NULL,
    discount_value INT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    valid_from TIMESTAMPTZ NOT NULL,
    valid_until TIMESTAMPTZ NOT NULL,
    coupon_code VARCHAR(10),
    minimum_amount INT,
    maximum_amount INT,
    FOREIGN KEY (product_id) REFERENCES store_products(id) ON DELETE CASCADE
);

-- Comment management
CREATE TABLE IF NOT EXISTS item_comments (
    id SERIAL PRIMARY KEY,
    item_id INT NOT NULL,
    commented_by UUID NOT NULL,
    title VARCHAR(255) NOT NULL DEFAULT '',
    message TEXT NOT NULL DEFAULT '',
    rating INT NOT NULL,
    up_votes INT NOT NULL DEFAULT 0,
    down_votes INT NOT NULL DEFAULT 0,
    is_edited BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    edited_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    FOREIGN KEY (item_id) REFERENCES store_items(id) ON DELETE CASCADE,
    FOREIGN KEY (commented_by) REFERENCES users(user_id) ON DELETE CASCADE
);

-- Serial management
CREATE TABLE IF NOT EXISTS store_devices (
    id SERIAL PRIMARY KEY,
    serie VARCHAR(25) UNIQUE NOT NULL,
    valid BOOLEAN NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS store_devices_history (
    id SERIAL PRIMARY KEY,
    device_id INT NOT NULL,
    issued_by VARCHAR(255) NOT NULL,
    issued_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    FOREIGN KEY (device_id) REFERENCES store_devices(id) ON DELETE RESTRICT
);

-- Order administration
CREATE TYPE order_status AS ENUM ('PENDIENTE', 'EN PROCESO', 'POR CONFIRMAR', 'ENTREGADO', 'CANCELADO');

CREATE SEQUENCE purchase_order_seq AS INT START WITH 100000;

CREATE TABLE IF NOT EXISTS store_orders (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    purchase_order INT DEFAULT nextval('purchase_order_seq'),
    email VARCHAR(255) NOT NULL,
    phone_number VARCHAR(25) NOT NULL,
    name VARCHAR(255) NOT NULL,
    dni VARCHAR(25) NOT NULL DEFAULT '',
    address TEXT NOT NULL,
    city VARCHAR(25) NOT NULL,
    postal_code VARCHAR(25) NOT NULL,
    assigned_to UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(purchase_order),
    FOREIGN KEY (assigned_to) REFERENCES users(user_id) ON DELETE SET NULL
);

CREATE TABLE IF NOT EXISTS order_products (
    id SERIAL PRIMARY KEY,
    order_id UUID NOT NULL,
    quantity INT NOT NULL DEFAULT 1,
    details HSTORE NOT NULL DEFAULT ''::hstore,
    product_id INT,
    product_type store_type NOT NULL,
    product_category VARCHAR(255) NOT NULL,
    product_item VARCHAR(255) NOT NULL,
    product_name VARCHAR(255) NOT NULL,
    product_price INT NOT NULL,
    product_details HSTORE NOT NULL DEFAULT ''::hstore,
    status order_status NOT NULL DEFAULT 'PENDIENTE',
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    FOREIGN KEY (order_id) REFERENCES store_orders(id) ON DELETE CASCADE,
    FOREIGN KEY (product_id) REFERENCES store_products(id) ON DELETE SET NULL
);

-- Order triggers
CREATE OR REPLACE TRIGGER set_order_timestamp
BEFORE UPDATE ON order_products
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();

-- Store triggers
CREATE OR REPLACE TRIGGER set_product_timestamp
BEFORE UPDATE ON store_products
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();

CREATE OR REPLACE TRIGGER set_item_timestamp
BEFORE UPDATE ON store_items
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();

CREATE OR REPLACE TRIGGER set_category_timestamp
BEFORE UPDATE ON store_categories
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();

CREATE OR REPLACE TRIGGER set_product_discount_timestamp
BEFORE UPDATE ON product_discount
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();

-- User triggers
CREATE OR REPLACE TRIGGER set_user_timestamp
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();

-- Serial triggers
CREATE OR REPLACE TRIGGER set_device_timestamp
BEFORE UPDATE ON store_devices
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();

-- Comment triggers
-- TODO