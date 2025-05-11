-- +goose Up
-- +goose StatementBegin
create index allData on Users (Name, Surname, Patronymic, Age, Nation, Gender);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop index allData;
-- +goose StatementEnd
