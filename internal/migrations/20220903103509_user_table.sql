-- +goose Up
CREATE TABLE IF NOT EXISTS users 
(
	account_number 				UUID 		PRIMARY KEY DEFAULT gen_random_uuid(),	
	name 	VARCHAR NOT NULL,
	email 			VARCHAR(50) 	NOT NULL,
    balance number NOT NULL,
	created_at 		TIMESTAMP 	NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS users ;

