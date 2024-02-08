-- +migrate Up
create table transactions (id varchar(255), transaction_type varchar(255), account_id_from varchar(255), account_id_to varchar(255), amount int, created_at date);

-- +migrate Down
drop table transactions;