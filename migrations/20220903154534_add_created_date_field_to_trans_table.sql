-- +goose Up
ALTER TABLE transactions 
    ADD created_date VARCHAR(20);


-- +goose Down
ALTER TABLE transactions 
    DROP created_date;
