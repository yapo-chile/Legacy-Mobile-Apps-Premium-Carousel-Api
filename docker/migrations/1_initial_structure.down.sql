DROP INDEX IF EXISTS user_product_user_id_idx;
DROP INDEX IF EXISTS user_product_user_email_idx;
DROP INDEX user_product_unique_active_product;
DROP INDEX user_product_unique_active_product_per_email;
DROP TABLE IF EXISTS user_product_param;
DROP TABLE IF EXISTS user_product;
DROP TABLE IF EXISTS purchase;
DROP TYPE enum_purchase_type;
DROP TYPE enum_product_type;
DROP TYPE enum_user_product_status;
DROP TYPE enum_user_product_param_name;
DROP TYPE enum_purchase_status;
