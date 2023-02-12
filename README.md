# core banking

a minimal implementation of a banking service, with support for:
- [X] creating users
- [X] depositing and withdrawing money
- [X] consistent transfer of funds from one user to another
- [X] transaction history (date, amount, balance)
- [X] transaction history filters (just deposits, withdrawal, date)
- [ ] idempotent requests
- [ ] introduce account tiers

### Setup
- clone the contents of `.env.sample` into `.env` and `.env.test` files
- build and run using `docker-compose up` and `docker-compose up --build` if you need to make any changes
- navigate into the `/migrations` dir and run goose migrations with `goose postgres "name=postgres password=password host=localhost:6543 dbname=corebanking sslmode=disable" up`, changing the dbname to corebanking_test for the test db
- hack away!


### API Docs (Endpoints)

- Create Account:

Sample `POST` request to `/createAccount` with `Content-Type`: `application/json`

```json
{
    "name": "James Worthy",
    "email": "JarryDeesNut5@gmail.com"
}
```
Sample response
```json
{
    "data": {
        "account_number": "eb155b3d-5a77-4ab5-969d-19955abc1f1b",
        "balance": "$0",
        "email": "JarryDeesNut5@gmail.com",
        "name": "James Worthy",
        "created_at": "2022-09-24 16:07:42"
    },
    "message": "user successfully created"
}
```

- Deposit and Withdraw from Account:
  
Sample `POST` request to `/deposit` or `/withdraw` with `Content-Type`: `application/json`

```json
{
    "account_number": "eb155b3d-5a77-4ab5-969d-19955abc1f1b",
    "amount": "$1,993.92"
}
```
Sample response
```json
{
    "data": {
        "account_number": "eb155b3d-5a77-4ab5-969d-19955abc1f1b",
        "balance": "$22,005.92",
        "name": "James Worthy",
        "updated_at": "2022-09-24 16:09:34"
    },
    "message": "deposit successful"
}
```

- Transfer funds:

Sample `POST` request to `/createAccount` with `Content-Type`: `application/json`

```json
{
    "account_number": "eb155b3d-5a77-4ab5-969d-19955abc1f1b",
    "recipient": "934ec732-6191-43f8-96b0-716b8e142346",
    "amount": "$500.00"
}
```
Sample response
```json
{
    "data": {
        "account_number": "eb155b3d-5a77-4ab5-969d-19955abc1f1b",
        "balance": "$11505.00",
        "name": "James Worthy",
        "recipient_name": "Dibri Nsofor",
        "updated_at": "2022-09-24 16:10:38"
    },
    "message": "transfer successful"
}
```

- View Transaction History

Sample `GET` request to `/transHistory` with `Content-Type`: `application/json` without query params

```json
{
    "account_number": "eb155b3d-5a77-4ab5-969d-19955abc1f1b"
}   
```
Sample response (**cut short to keep this pithy**)
```json
{
    "data": [
        {
            "id": "97cd4420-a6c5-4c64-ab69-e6a9678d392d",
            "account_number": "eb155b3d-5a77-4ab5-969d-19955abc1f1b",
            "action_performed": "DEPOSIT",
            "balance": "$19.00",
            "created_at": "2022-09-24T16:08:31.434155Z",
            "created_date": "2022-09-24"
        },
        {
            "id": "3eceb330-1ae7-4183-a893-22275144eb76",
            "account_number": "eb155b3d-5a77-4ab5-969d-19955abc1f1b",
            "action_performed": "TRANSFER",
            "recipient": "934ec732-6191-43f8-96b0-716b8e142346",
            "balance": "$11,505.00",
            "created_at": "2022-09-24T16:10:38.787012Z",
            "created_date": "2022-09-24"
        }
    ],
    "message": "successfully retrieved 10 most recent transactions"
}
```

- View Transaction History (with query params)

Sample `POST` request to `/transHistory?date=&action=` with `Content-Type`: `application/json`

```json
{
    "account_number": "eb155b3d-5a77-4ab5-969d-19955abc1f1b"
}   
```
Sample response (**cut short for readability sake**)
```json
{
    "data": [
        {
            "id": "3eceb330-1ae7-4183-a893-22275144eb76",
            "account_number": "eb155b3d-5a77-4ab5-969d-19955abc1f1b",
            "action_performed": "TRANSFER",
            "recipient": "934ec732-6191-43f8-96b0-716b8e142346",
            "balance": "$11,505.00",
            "created_at": "2022-09-24T16:10:38.787012Z",
            "created_date": "2022-09-24"
        }
    ],
    "message": "successfully retrieved transactions"
}
```