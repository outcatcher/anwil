# API Reference

Current API version is `v1`.
All endpoints should be prefixed with `/api/v1`.

## Authentication

Most of available endpoints requires user to be authenticated.

Currently, standard JWT auth is used:
`Authorization` request header should have value `Bearer <token>`.

## Endpoints

### Debug endpoints

#### `GET /api/v1/echo`

Responds with "OK" on each request.

```shell
$ curl http://localhost:8010/api/v1/echo
OK
```

#### `GET /api/v1/auth-echo`

Responds with "OK" on each request with valid JWT.

```shell
$ curl http://localhost:8010/api/v1/auth-echo -H "Authorization: Bearer $TOKEN"
OK
```

### User endpoints

#### `POST /api/v1/wisher`

Create new user.

#### `POST /api/v1/login`

Authorize user.

For details see [user API reference](../users/handlers/README.md).
