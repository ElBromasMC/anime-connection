/* Aspectos generales de la base de datos:
- Se prefiere el valor 'cero' de los tipos antes que NULL
  para mejorar la integración con el servidor web */
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pg_trgm";
CREATE EXTENSION IF NOT EXISTS "hstore";

-- Row update management
/* Función que actualiza la columna 'updated_at' automáticamente
después de un comando UPDATE en una fila */
CREATE OR REPLACE FUNCTION trigger_set_timestamp()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- User administration
/* ADMIN: Representa a los administradores de la página
   SELLER: Representa a los vendedores autenticados
   MODERATOR: Representa a los moderadores de la página
   NORMAL: Representa al usuario interesado en los productos */
CREATE TYPE user_role AS ENUM ('ADMIN', 'SELLER', 'MODERATOR', 'NORMAL');

/* Entidad que representa a los usuarios del sistema, la
contraseña se debe encriptar antes de almacenarla */
CREATE TABLE IF NOT EXISTS users (
    user_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    hashed_password TEXT NOT NULL,
    role user_role NOT NULL DEFAULT 'NORMAL',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

/* Entidad que representa las sesiones de los usuarios en la
página, se eliminan automáticamente después de 1 mes. Por lo que,
los usuarios deben logearse nuevamente */
CREATE TABLE IF NOT EXISTS sessions (
    session_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMPTZ NOT NULL DEFAULT NOW() + INTERVAL '1 month',
    FOREIGN KEY (user_id) REFERENCES users(user_id)
);

-- Image administrarion
/* Entidad que representa las imágenes, se almacenan usando el
sistema de archivos del sistema operativo. La base de datos
solo se almacena el enlace a dicha imagen */
CREATE TABLE IF NOT EXISTS images (
    id SERIAL PRIMARY KEY,
    filename VARCHAR(25) UNIQUE NOT NULL
);

-- Store administration
/* MANGA: Tipo que representa todo lo relacionado a la venta de manga
*/
CREATE TYPE store_type AS ENUM ('MANGA');

/* Entidad que representa las categorías creadas por los administradores
de la página. En las que se pueden clasificar los items de la tienda */
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

/* Entidad que representa a los items de la tienda, aquellos publicados
por los vendedores, pueden incluir más de un producto. Por ejemplo,
pequeñas variaciones o adiciones. Usa un índice GIN para optimizar
las búsquedas por nombre */
CREATE TABLE IF NOT EXISTS store_items (
    id SERIAL PRIMARY KEY,
    category_id INT NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    long_description TEXT NOT NULL DEFAULT '',
    img_id INT,
    largeimg_id INT,
    slug VARCHAR(255) NOT NULL,
    created_by UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(category_id, slug),
    FOREIGN KEY (category_id) REFERENCES store_categories(id) ON DELETE CASCADE,
    FOREIGN KEY (created_by) REFERENCES users(user_id) ON DELETE CASCADE,
    FOREIGN KEY (img_id) REFERENCES images(id) ON DELETE SET NULL,
    FOREIGN KEY (largeimg_id) REFERENCES images(id) ON DELETE SET NULL
);
CREATE INDEX idx_items_name ON store_items USING gin (name gin_trgm_ops);

-- Product management
/* Entidad que representa a los productos asociados a un item. Representan
pequeñas variaciones o adiciones del item que se ofrece. La columna price
almacena el precio base. */
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

/* Descuentos asociados a los productos. Podrán ser reclamados mediante
un cupón. Almacena el periodo en el que es válido y opcionalmente la cantidad
mínima y máxima para el cual tiene efecto */
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
/* Entidad que representa los comentarios y valoraciones asociadas a un item
de la tienda */
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

-- Order administration
/* Representa los estados en los que pueden estar los productos
asociados a una orden de compra */
CREATE TYPE order_status AS ENUM ('PENDIENTE', 'EN PROCESO', 'POR CONFIRMAR', 'ENTREGADO', 'CANCELADO');

CREATE SEQUENCE purchase_order_seq AS INT START WITH 100000;
/* Entidad que representa las órdenes de compra, las que se
generan una vez completado el checkout y el pago. Sirve como
historial y le da información al cliente sobre el estado
de su órden */
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

/* Entidad que representa los productos asociados a las órdenes
de compra, almacena a manera de historial la información del producto
al momento que lo compró */
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

/* Triggers que ejecutan la función 'trigger_set_timestamp' antes de
una instrucción UPDATE */
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
