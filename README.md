# userService

User and auth service for my wallet project. It handles sign-up and login with JWT, stores which Telegram chat belongs to which user, and on sign-up asks [newWalletService](https://github.com/irenakk/newWalletService) over gRPC to create a wallet with USD, EUR and RUB accounts.

## How it fits together

- userService (HTTP :8082, gRPC :9090) handles sign-up and login and issues JWTs. On sign-up it calls `CreateWallet` and `CreateAccount` on newWalletService.
- [newWalletService](https://github.com/irenakk/newWalletService) (HTTP :8081, gRPC :9091) accepts the same JWT, looks up user ids here via the `Find` RPC, and writes a Kafka message after every top-up or transfer.
- [telegramMessenger](https://github.com/irenakk/telegramMessenger) is a Telegram bot. It calls `/link-telegram` and `/get-chatid` here to link a chat to a user and to find where to send a notification.

## Stack

Go 1.23, Gin, gRPC/protobuf, PostgreSQL (`database/sql` + `lib/pq`, plain SQL), JWT HS256 (`golang-jwt/jwt/v5`), bcrypt for passwords.

## How to run

You need PostgreSQL. For sign-up, newWalletService has to be running on `localhost:9091` too.

1. Create the database `users_walletservice` with a `users` table: `id`, `username`, `password`, `telegram_nickname`, `telegram_chat_id`.
2. Set the config, from env variables or a `.env` file in the working directory:

   | Variable | Default |
   |---|---|
   | `SERVER_HOST` / `SERVER_PORT` | `0.0.0.0` / `8082` |
   | `DB_HOST` / `DB_PORT` | `localhost` / `5432` |
   | `DB_USER` | `postgres` |
   | `DB_PASSWORD` | set your own |
   | `DB_NAME` | `users_walletservice` |
   | `DB_SSLMODE` | `disable` |
   | `JWT_SECRET` | use the same value in newWalletService |
   | `ENV` | `development` (`production` puts Gin in release mode) |

3. Start it:

   ```
   go run .
   ```

HTTP listens on `:8082`, gRPC on `:9090`.

## HTTP API

All routes are under `/api/v1`.

| Method | Path | Auth | Body | Notes |
|---|---|---|---|---|
| POST | `/register` | - | `username`, `telegram_nickname`, `password` | password min 6 chars; also creates the wallet and 3 accounts |
| POST | `/login` | - | `username`, `password` | returns a JWT valid for 24h |
| GET | `/hello` | Bearer | - | returns `Hello <username>`, handy to test a token |
| POST | `/link-telegram` | - | `username`, `chat_id`, `tg_nickname` | called by the bot |
| POST | `/get-chatid` | - | `username` | returns `{"chat_id": ...}`, called by the bot |

Example:

```
curl -X POST http://localhost:8082/api/v1/register \
  -H "Content-Type: application/json" \
  -d '{"username": "alice", "telegram_nickname": "alice_tg", "password": "secret123"}'

curl -X POST http://localhost:8082/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{"username": "alice", "password": "secret123"}'
# {"expires_in": 86400, "token": "...", "token_type": "Bearer"}
```

## gRPC

`user.UserService` on `:9090`, defined in `proto/server/user.proto`. It has one method, `Find(username)`, which returns `user_id` and `username`. newWalletService uses it to resolve usernames from the JWT.

Server reflection is on:

```
grpcurl -plaintext -d '{"username": "alice"}' localhost:9090 user.UserService/Find
```

`proto/client/wallet.proto` is the newWalletService contract used for creating wallets. The generated `.pb.go` files are committed.

## Project structure

```
main.go                wiring: HTTP server, gRPC server, gRPC client to the wallet service
config/                env config, Postgres connection
internal/handler/      HTTP handlers
internal/middleware/   JWT check for protected routes
internal/usecase/      sign-up (user + wallet + accounts), login, JWT, Telegram linking
internal/repository/   SQL for the users table
internal/rpctransfer/  gRPC handler for Find
internal/service/      gRPC client for WalletService
internal/utils/        bcrypt helpers
proto/                 .proto files and generated code
```
