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
    'custom_query',
    'exclude',
    'price_range',
    'gaps_with_random'
);

CREATE TABLE IF NOT EXISTS user_product(
	id              SERIAL PRIMARY KEY,
	product_type    enum_product_type NOT NULL,
	user_id         INTEGER NOT NULL,
	user_email      TEXT NOT NULL,
	status          enum_user_product_status NOT NULL,
	expired_at      TIMESTAMP NOT NULL,
	created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	comment         TEXT
);

-- create index to allow only one active product type per user
create unique index user_product_unique_active_product on user_product(product_type, user_id, status)
    where status = 'ACTIVE';

CREATE TABLE IF NOT EXISTS user_product_history(
	user_product_id INTEGER REFERENCES user_product(id),
	status         enum_user_product_status NOT NULL,
	timestamp      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS user_product_config(
	user_product_id INTEGER REFERENCES user_product(id),
	name            enum_user_product_config_name NOT NULL,
	value           TEXT,
	unique (user_product_id, name)

);

CREATE TABLE IF NOT EXISTS carousel_report(
	user_product_id INTEGER REFERENCES user_product(id),
	list_id         INTEGER NOT NULL,
	views_counter   NUMERIC NOT NULL DEFAULT 0
);

CREATE INDEX user_product_user_email_idx ON user_product(user_email);
CREATE INDEX user_product_user_id_idx ON user_product(user_id);
