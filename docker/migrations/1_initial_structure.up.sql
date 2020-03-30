CREATE TYPE enum_product_type AS ENUM (
    'PREMIUM_CAROUSEL'
);

CREATE TYPE enum_user_product_status AS ENUM (
    'INACTIVE',
    'ACTIVE',
    'EXPIRED'
);

CREATE TYPE enum_user_product_config_name AS ENUM (
    'categories',
    'limit',
    'keywords',
    'exclude',
    'price_range',
    'fill_random',
    'comment'
);

CREATE TYPE enum_purchase_type AS ENUM (
    'ADMIN'
);

CREATE TYPE enum_purchase_status AS ENUM (
    'PENDING',
    'ACCEPTED',
    'REJECTED'
);

CREATE TABLE IF NOT EXISTS purchase(
    id              SERIAL PRIMARY KEY,
    purchase_number INTEGER NOT NULL DEFAULT 0,
    purchase_type   enum_purchase_type DEFAULT 'ADMIN',
    purchase_status enum_purchase_status DEFAULT 'PENDING',
    price           INTEGER NOT NULL DEFAULT 0,
    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS user_product(
    id              SERIAL PRIMARY KEY,
    product_type    enum_product_type NOT NULL,
    user_id         INTEGER NOT NULL,
    user_email      TEXT NOT NULL,
    purchase_id     INTEGER REFERENCES purchase(id),
    status          enum_user_product_status NOT NULL,
    start_at        TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expired_at      TIMESTAMP NOT NULL,
    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- create index to allow only one active product type per user_id
CREATE unique index user_product_unique_active_product on user_product(product_type, user_id, status)
    where status = 'ACTIVE';
-- create index to allow only one active product type per user_email
CREATE unique index user_product_unique_active_product_per_email on user_product(product_type, user_email, status)
    where status = 'ACTIVE';

CREATE TABLE IF NOT EXISTS user_product_config(
    user_product_id INTEGER REFERENCES user_product(id),
    name            enum_user_product_config_name NOT NULL,
    value           TEXT,
    unique (user_product_id, name)
);

CREATE INDEX user_product_user_email_idx ON user_product(user_email);
CREATE INDEX user_product_user_id_idx ON user_product(user_id);
