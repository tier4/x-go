create table "users"
(
    "id"    serial primary key,
    "email" varchar(255) not null unique
);
