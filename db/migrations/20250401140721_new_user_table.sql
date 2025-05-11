-- +goose Up
-- +goose StatementBegin
create table if not exists Users (
ID serial primary key,
Name varchar(255) not null,
Surname varchar(255) not null,
Patronymic varchar(255),
Nation varchar(255) not null,
Gender varchar(255) not null,
Age integer
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists Users;
-- +goose StatementEnd
