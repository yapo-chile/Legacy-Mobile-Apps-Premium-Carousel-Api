-- create index to allow only one active product type per user_id
CREATE unique index user_product_unique_active_product on user_product(product_type, user_id, status)
    where status = 'ACTIVE';
-- create index to allow only one active product type per user_email
CREATE unique index user_product_unique_active_product_per_email on user_product(product_type, user_email, status)
    where status = 'ACTIVE';
