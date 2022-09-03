# core banking

a minimal implementation of a banking service, with support for:
- [X] creating users
- [X] depositing and withdrawing money
- [X] consistent transfer of funds from one user to another
- [ ] transaction history (date, amount, balance)
- [ ] transaction history filters (just deposits, withdrawal, date)
- [ ] idempotent transactions

### Setup

### API Docs (Endpoints)

- Create Account:

Sample `POST` request to `/createAccount` with `Content-Type`: `application/json`

```json
{
    "name": "Dibri Nsofor",
    "email": "jamesdeez@gmail.com"
}
```
Sample response
```json
{
    "data": {
        "account_number": "6eb55ee8-cdcb-4819-9301-0ab1c3a5cb21",
        "name": "Dibri Nsofor",
        "email": "dibrinsofor@gmail.com",
        "balance": 0,
        "created_at": "2022-09-03T10:41:29.5850223+01:00"
    },
    "message": "user successfully created"
}
```

- Deposit and Withdraw from Account:
  
Sample `POST` request to `/deposit` or `/withdraw` with `Content-Type`: `application/json`

```json
{
    "account_number": "6eb55ee8-cdcb-4819-9301-0ab1c3a5cb21",
    "amount": 12000
}
```
Sample response
```json
{
    "data": {
        "account_number": "6eb55ee8-cdcb-4819-9301-0ab1c3a5cb21",
        "name": "Dibri Nsofor",
        "email": "dibrinsofor@gmail.com",
        "balance": 150000,
        "created_at": "2022-09-03T10:41:29.585022Z"
    },
    "message": "user withdrawal successful"
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
        "balance": 50000,
        "created_at": "2022-09-03T10:41:29.585022Z"
    },
    "message": "transfer successful"
}
```
