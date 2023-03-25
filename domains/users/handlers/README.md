# User service handlers

## POST `/wisher`

Creates new user with unique username.

### Request attributes

---

**username** `string`

*Required*

Username should be a unique and contains 2 to 20 UTF-8 characters.

--- 

**password** `string`

*Required*

Password should have no less than 50 entropy bits,
see [description](https://github.com/wagslane/go-password-validator#what-entropy-value-should-i-use).

Valid password can be something like `Let_Me_In_Please!`, but not `aN!9PP\/F-`, and certainly not `qwerty`.

---

**full_name** `string`

*Optional*

A full name or an alias of the user to be shown in UI.

---

### Example

```shell
curl -X POST http://localhost:8010/api/v1/wisher -d @new_user.json -H "content-type: application/json" 
```

**new_user.json**
```json
{
  "username": "unique",
  "password": "Z3XKLtqeyoYQSwwK",
  "full_name": "John Doe"
}
```

### Response

Statuses:

- `201`: User successfully created
- `400`: Input attributes restrictions not met
- `409`: User with given username already exists

## POST `/login`

Authorizes user returning new token.

### Request attributes

---

**username** `string`

*Required*

---

**password** `string`

*Required*

---

### Example

```shell
$ curl -X POST http://localhost:8010/api/v1/login -d @login_user.json -H "content-type: application/json"

{"token":"eyJhbGciOiJFZERTQSIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJhbndpbCIsImV4cCI6MTY4MjM3NjExOSwiaWF0IjoxNjc5Nzg0MTE5LCJ1c2VybmFtZSI6InVuaXF1ZSJ9.jzmRmvw9IuueTogPbzXHKprKZv-TAyGJIPHSYom8HTV6tqUKSp3q8Z6WHrz-l35jQB7bX35Ncm6x35QBuJf8Bg"}
```

**login_user.json**

```json
{
  "username": "unique",
  "password": "Z3XKLtqeyoYQSwwK"
}
```

### Response

Statuses:

- `200`: Token successfully created
- `400`: Request body invalid
- `401`: Password invalid
- `404`: User doesn't exist

### Response Body

---

**token** `string`

Base64-encoded JWT.

---

