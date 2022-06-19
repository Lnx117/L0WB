CREATE TABLE orders (
    order_uid               varchar(255) PRIMARY KEY,
    track_number            varchar(255),
    entry                   varchar(255),
    locale                  varchar(255),
    internal_signature      varchar(255),
    customer_id             varchar(255),
    delivery_service        varchar(255),
    shardkey                varchar(255),
    sm_id                   int,
    oof_shard               varchar(255),
    date_created            timestamp,
    recording_date          timestamp
);

CREATE TABLE delivery (
    order_uid               varchar(255) references orders (order_uid) on delete cascade not null unique,
    name                    varchar(255),
    phone                   varchar(255),
    zip                     varchar(255),
    city                    varchar(255),
    address                 varchar(255),
    region                  varchar(255),
    email                   varchar(255),
    recording_date          timestamp
);

CREATE TABLE payment (
    order_uid               varchar(255) references orders (order_uid) on delete cascade not null unique,
    transaction             varchar(255),
    request_id              varchar(255),
    currency                varchar(255),
    provider                varchar(255),
    amount                  int,
    payment_dt              int,
    bank                    varchar(255),
    delivery_cost           int,
    goods_total             int,
    custom_fee              int,
    recording_date          timestamp
);

CREATE TABLE items (
    order_uid               varchar(255) references orders (order_uid) on delete cascade not null,
    chrt_id                 int,
    track_number            varchar(255),
    price                   int,
    rid                     varchar(255),
    name                    varchar(255),
    sale                    int,
    size                    varchar(255),
    total_price             int,
    nm_id                   int,
    brand                   varchar(255),
    status                  int,
    recording_date          timestamp
);
