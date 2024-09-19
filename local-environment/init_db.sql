CREATE TABLE IF NOT EXISTS orders
(
    id                 varchar(255) unique not null,
    order_date         timestamptz,
    order_status       varchar(255),
    vehicle_id         varchar(255),
    price              integer,
    make               varchar(255),
    model              varchar(255),
    version            varchar(255),
    color              varchar(255),
    model_year         integer,
    mileage_kilometers integer
);


create index idx_orders_order_status on orders(order_status);
create index idx_orders_model_year on orders(model_year);
create index idx_orders_color on orders(color);
create index idx_orders_order_date on orders(order_date);


SET datestyle = 'ISO';
COPY orders FROM '/docker-entrypoint-initdb.d/case_order_data.csv' DELIMITER ',' CSV HEADER ;
