-- +migrate Up
create table accounts (id varchar(255), client_id varchar(255), account_type varchar(255), balance int, created_at date, update_at date);

-- +migrate Down
drop table accounts;