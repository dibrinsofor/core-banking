-- +goose Up
CREATE TABLE IF NOT EXISTS transactions
(
	id 				UUID 		PRIMARY KEY DEFAULT gen_random_uuid(),	
	account_number 	VARCHAR NOT NULL,
	action_performed			VARCHAR(20) 	NOT NULL,
    recipient VARCHAR,
    balance bigint NOT NULL,
	created_at 		TIMESTAMP 	NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS transactions ;