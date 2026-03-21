# Chirpy

A Twitter-like REST API server built in Go with PostgreSQL. Supports user registration, JWT authentication, chirp (post) management, and premium upgrades via webhooks.

## Table of Contents

- [Setup](#setup)
- [Environment Variables](#environment-variables)
- [Database Migrations](#database-migrations)
- [API Reference](#api-reference)
  - [Health](#health)
  - [Users](#users)
  - [Authentication](#authentication)
  - [Chirps](#chirps)
  - [Webhooks](#webhooks)
  - [Admin](#admin)
- [Authentication Flow](#authentication-flow)
- [Error Responses](#error-responses)

---

## Setup

**Requirements:** Go 1.21+, PostgreSQL

```bash
git clone <repo>
cd webserver
cp .env.example .env   # fill in your values
go build ./...
go run .
```

Server starts on port `8080`.

---

## Environment Variables

Create a `.env` file in the project root:

```env
DB_URL=postgres://user:password@localhost:5432/chirpy?sslmode=disable
PLATFORM=dev
JWT_SECRET=your-secret-key-here
POLKA_KEY=your-polka-api-key-here
```

| Variable     | Required | Description                                              |
|--------------|----------|----------------------------------------------------------|
| `DB_URL`     | Yes      | PostgreSQL connection string                             |
| `JWT_SECRET` | Yes      | Secret key for signing JWTs (HS256)                      |
| `POLKA_KEY`  | Yes      | API key for validating Polka webhook requests            |
| `PLATFORM`   | No       | Set to `dev` to enable the `POST /admin/reset` endpoint  |

---

## Database Migrations

Migrations are plain SQL files in `sql/schema/`. Run them in order against your PostgreSQL database:

```bash
psql $DB_URL -f sql/schema/001_users.sql
psql $DB_URL -f sql/schema/002_chirps.sql
psql $DB_URL -f sql/schema/003_users.sql
psql $DB_URL -f sql/schema/004_refresh_tokens.sql
psql $DB_URL -f sql/schema/005_users.sql
```

---

## API Reference

### Health

#### `GET /api/healthz`

Returns server status.

**Response `200 OK`**
```
OK
```

---

### Users

#### `POST /api/users`

Register a new user.

**Request Body**
```json
{
  "email": "user@example.com",
  "password": "secret123"
}
```

**Response `201 Created`**
```json
{
  "id": "uuid",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z",
  "email": "user@example.com",
  "is_chirpy_red": false
}
```

---

#### `PUT /api/users`

Update the authenticated user's email and password.

**Headers**
```
Authorization: Bearer <access_token>
```

**Request Body**
```json
{
  "email": "new@example.com",
  "password": "newpassword"
}
```

**Response `200 OK`**
```json
{
  "id": "uuid",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z",
  "email": "new@example.com",
  "is_chirpy_red": false
}
```

---

### Authentication

#### `POST /api/login`

Login and receive an access token + refresh token.

**Request Body**
```json
{
  "email": "user@example.com",
  "password": "secret123"
}
```

**Response `200 OK`**
```json
{
  "id": "uuid",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z",
  "email": "user@example.com",
  "is_chirpy_red": false,
  "token": "<jwt_access_token>",
  "refresh_token": "<refresh_token>"
}
```

- `token`: JWT valid for **1 hour**
- `refresh_token`: Opaque token valid for **60 days**

---

#### `POST /api/refresh`

Exchange a valid refresh token for a new access token.

**Headers**
```
Authorization: Bearer <refresh_token>
```

**Response `200 OK`**
```json
{
  "token": "<new_jwt_access_token>"
}
```

---

#### `POST /api/revoke`

Revoke a refresh token.

**Headers**
```
Authorization: Bearer <refresh_token>
```

**Response `204 No Content`**

---

### Chirps

#### `POST /api/chirps`

Create a new chirp. Maximum 140 characters. Profane words (`kerfuffle`, `sharbert`, `fornax`) are replaced with `****`.

**Headers**
```
Authorization: Bearer <access_token>
```

**Request Body**
```json
{
  "body": "Hello world!"
}
```

**Response `201 Created`**
```json
{
  "id": "uuid",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z",
  "body": "Hello world!",
  "user_id": "uuid"
}
```

---

#### `GET /api/chirps`

List all chirps. Optionally filter by author and sort order.

**Query Parameters**

| Parameter   | Type   | Description                                 |
|-------------|--------|---------------------------------------------|
| `author_id` | UUID   | Filter chirps by this user's ID             |
| `sort`      | string | `desc` for newest first (default: ascending)|

**Response `200 OK`**
```json
[
  {
    "id": "uuid",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z",
    "body": "Hello world!",
    "user_id": "uuid"
  }
]
```

---

#### `GET /api/chirps/{chirpID}`

Get a single chirp by ID.

**Response `200 OK`**
```json
{
  "id": "uuid",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z",
  "body": "Hello world!",
  "user_id": "uuid"
}
```

---

#### `DELETE /api/chirps/{chirpID}`

Delete a chirp. Only the chirp's author can delete it.

**Headers**
```
Authorization: Bearer <access_token>
```

**Response `204 No Content`**

---

### Webhooks

#### `POST /api/polka/webhooks`

Webhook endpoint for Polka payment events. Upgrades a user to Chirpy Red on `user.upgraded` events.

**Headers**
```
Authorization: ApiKey <polka_key>
```

**Request Body**
```json
{
  "event": "user.upgraded",
  "data": {
    "user_id": "uuid"
  }
}
```

**Response `204 No Content`**

Other event types are acknowledged with `204` and ignored.

---

### Admin

#### `GET /admin/metrics`

Returns an HTML page showing the number of times the fileserver has been accessed.

#### `POST /admin/reset`

Resets the fileserver hit counter and deletes all users. **Only available when `PLATFORM=dev`.**

**Response `200 OK`** (dev) / **`403 Forbidden`** (non-dev)

---

## Authentication Flow

```
1. POST /api/users          → create account
2. POST /api/login          → get access_token (1h) + refresh_token (60d)
3. Use access_token         → Bearer token on protected endpoints
4. POST /api/refresh        → get new access_token when expired
5. POST /api/revoke         → invalidate refresh_token on logout
```

---

## Error Responses

All errors return JSON:

```json
{
  "error": "description of the error"
}
```

| Status | Meaning                                      |
|--------|----------------------------------------------|
| `400`  | Bad request (invalid input, missing fields)  |
| `401`  | Unauthorized (missing or invalid token)      |
| `403`  | Forbidden (wrong ownership or API key)       |
| `404`  | Resource not found                           |
| `500`  | Internal server error                        |
