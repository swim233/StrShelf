create table shelf_item_v1
(
    id           bigint generated always as identity
        primary key,
    title        text                    not null,
    link         text                    not null,
    comment      text,
    gmt_created  timestamp default now() not null,
    gmt_modified timestamp default now() not null,
    deleted      boolean                 not null,
    gmt_deleted  timestamp default now() not null
);

create table shelf_user_v1
(
    id          bigint generated always as identity
        primary key,
    username    text                    not null,
    password    text                    not null,
    gmt_created timestamp default now() not null
);
