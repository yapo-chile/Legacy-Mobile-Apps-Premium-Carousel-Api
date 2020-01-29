CREATE TYPE enum_product_type AS ENUM (
	'PREMIUM_CAROUSEL',
);

CREATE TABLE IF NOT EXISTS user_product(
	id              SERIAL PRIMARY KEY,
	product_type    enum_product_type NOT NULL,
	user_email      VARCHAR(50) NOT NULL,
	comment         text,
	expired_at      TIMESTAMP NOT NULL,
	created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS carousel_report(
	id              SERIAL PRIMARY KEY,
	user_product_id INTEGER REFERENCES user_product(id),
	list_id      	INTEGER NOT NULL,
	views_counter   BIGINT NOT NULL DEFAULT 0
);

CREATE INDEX user_product_user_email_idx ON user_product(user_email);
