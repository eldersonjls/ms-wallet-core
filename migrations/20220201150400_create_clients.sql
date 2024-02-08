-- +migrate Up
create table clients (id varchar(255), name varchar(255), email varchar(255), bank_id varchar(255), date_of_birth date, created_at date, update_at date);

-- +migrate Down
drop table clients;