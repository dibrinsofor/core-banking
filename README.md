# core banking

a minimal implementation of a banking service, with support for:
- [X] creating users
- [X] depositing and withdrawing money
- [X] consistent transfer of funds from one user to another
- [X] transaction history (date, amount, balance)
- [X] transaction history filters (just deposits, withdrawal, date)
- [X] idempotent requests
- [X] request timeouts
- [ ] introduce account tiers

### Setup
- clone the contents of `.env.sample` into `.env` and `.env.test` files
- build and run using `docker-compose up` and `docker-compose up --build` if you need to make any changes
- navigate into the `/migrations` dir and run goose migrations with `goose postgres "name=postgres password=password host=localhost:6543 dbname=corebanking sslmode=disable" up`, changing the dbname to corebanking_test for the test db
- hack away!


### API Docs (Endpoints)

The client is responsible for sending `Idempotency-Key` Headers (preferably UUID) along with `POST` requests. [Read more here](https://stripe.com/blog/idempotency)

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
    "account_number": "6eb55ee8-cdcb-4819-9301-0ab1c3a5cb21",
    "recipient": "017d7b89-8d38-488a-8e0a-8289dbbb427e",
    "amount": 20000
}
```
Sample response
```json
{
    "data": {
        "account_number": "6eb55ee8-cdcb-4819-9301-0ab1c3a5cb21",
        "name": "Dibri Nsofor",
        "email": "dibrinsofor@gmail.com",
        "balance": "$22,005.92",
        "created_at": "2022-09-03T10:41:29.585022Z"
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
            "id": "416ef053-1f5b-43d2-9478-10d6a9c98fce",
            "account_number": "9cbf9d63-8510-4f20-928a-80a75818ebb1",
            "action_performed": "Deposit",
            "recipient": "",
            "balance": "$12,000.00",
            "created_at": "2022-09-03T13:29:42.94933Z"
        },
        {
            "id": "0e2c72b8-2b42-4fff-9d5e-71f968bdf3af",
            "account_number": "9cbf9d63-8510-4f20-928a-80a75818ebb1",
            "action_performed": "Deposit",
            "recipient": "",
            "balance": "$24,000.00",
            "created_at": "2022-09-03T13:29:50.259873Z"
        },
        {
            "id": "6d20b744-b4ec-4f33-95b5-c517bbc775dc",
            "account_number": "9cbf9d63-8510-4f20-928a-80a75818ebb1",
            "action_performed": "WITHDRAW",
            "recipient": "",
            "balance": "$134,005.92",
            "created_at": "2022-09-03T14:05:28.708591Z"
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
            "id": "2b69db60-d99a-46b6-8468-bbc2745a4a5f",
            "account_number": "9cbf9d63-8510-4f20-928a-80a75818ebb1",
            "action_performed": "DEPOSIT",
            "recipient": "",
            "balance": "$329,492,080.92",
            "created_at": "2022-09-03T16:06:36.625368Z",
            "created_date": "2022-09-03"
        },
        {
            "id": "0b335e05-8645-43e5-8ba8-b5f4f446dac2",
            "account_number": "9cbf9d63-8510-4f20-928a-80a75818ebb1",
            "action_performed": "DEPOSIT",
            "recipient": "",
            "balance": "$389,385,560.00",
            "created_at": "2022-09-03T16:07:10.709543Z",
            "created_date": "2022-09-03"
        },
        {
            "id": "a6182b27-3dd1-4d35-a281-2adb4eb93c10",
            "account_number": "9cbf9d63-8510-4f20-928a-80a75818ebb1",
            "action_performed": "DEPOSIT",
            "recipient": "",
            "balance": 389385572,
            "created_at": "2022-09-03T16:07:12.816867Z",
            "created_date": "2022-09-03"
        }
    ],
    "message": "successfully retrieved transactions"
}
```
