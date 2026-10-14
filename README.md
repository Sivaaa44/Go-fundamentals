# mini-redis

A small Redis-compatible in-memory key-value server written in Go, built from scratch using only the standard library. It speaks the Redis wire protocol (RESP), handles many clients concurrently, supports key expiry, and persists writes to an append-only file so data survives restarts.

I built this to get myself familiar with GO and also to learn how systems work under the hood.

## Features

- **RESP protocol**: hand-written parser and reply encoder (no libraries)
- **Concurrent clients**: one goroutine per connection, sharing a mutex-protected store
- **Key expiry**: `SET key value EX seconds`, with both lazy expiry (checked on read) and active expiry (a background sweeper)
- **Persistence**: every successful write is appended to `appendonly.aof` and replayed on startup
- **Graceful Ctrl+C**: the log file is closed cleanly on shutdown
- **Tests**: unit tests for the store, parser, and command handler

## Supported commands

| Command | Description | Reply |
|---|---|---|
| `PING` | Health check | `+PONG` |
| `SET key value [EX seconds]` | Store a string, optionally with a TTL | `+OK` |
| `GET key` | Fetch a value | bulk string, or null (`$-1`) if missing/expired |
| `DEL key` | Delete a key | `:1` if deleted, `:0` if it didn't exist |
| `EXISTS key` | Check whether a key exists | `:1` or `:0` |
| `TTL key` | Seconds until expiry | seconds, `-1` if no expiry, `-2` if missing |

Command names are case-insensitive. Wrong argument counts and unknown commands return a RESP error (`-ERR ...`).

## Getting started

Requires Go 1.21 or newer.

```bash
git clone https://github.com/<your-username>/mini-redis.git
cd mini-redis
go run .
```

The server listens on port `6379`. On startup it replays `appendonly.aof` if it exists.


## Trying it out

Any client that speaks RESP can connect. With `redis-cli` (via WSL, Docker, or a local install):

```bash
redis-cli -p 6379
127.0.0.1:6379> SET name siva EX 30
OK
127.0.0.1:6379> GET name
"siva"
127.0.0.1:6379> TTL name
(integer) 29
```

Or talk to it with raw bytes. This sends `SET name siva`:

```
*3\r\n$3\r\nSET\r\n$4\r\nname\r\n$4\r\nsiva\r\n
```

## How it works

```
client ──TCP──► main.go ──► parser.go ──► handler.go ──► store.go
                                              │
                                              └──► aof.go  (append writes to disk)
```

| File | Responsibility |
|---|---|
| `main.go` | Startup sequence, TCP listener, per-connection goroutine loop |
| `parser.go` | Turns RESP bytes into a `[]string` command. Reads exact bulk-string lengths with `io.ReadFull`, so values may contain spaces or newlines |
| `handler.go` | Maps a parsed command to a store operation and builds the RESP reply |
| `store.go` | The key-value map, guarded by a `sync.Mutex`, plus expiry logic and the background sweeper |
| `aof.go` | Appends write commands to disk and replays them on startup |


### Design notes

- **Concurrency model**: real Redis executes commands on a single thread. This project uses a goroutine per client with a mutex around the shared map instead, which is the idiomatic Go approach.
- **Expiry**: each entry stores an optional deadline. `Get` deletes expired keys it touches (lazy), and a ticker-driven goroutine removes the rest every second (active), so unread keys don't pile up.
- **Persistence**: the AOF stores commands in the same RESP format clients send, so replay reuses the exact same parser and handler. Replay runs before the log is opened for appending, so replayed commands are not logged twice. A truncated final command (for example after a crash) stops replay cleanly instead of failing startup.
- **Layering**: the parser knows nothing about the store, and the store knows nothing about the network. That is what makes each layer testable with plain strings.


## Known limitations

This is a learning project, not a production database.

- **Strings only**: no lists, sets, hashes, or sorted sets
- **TTLs restart on replay**: the AOF logs `EX 10` as a relative duration, so a replayed key gets a fresh 10 seconds. Real Redis logs an absolute deadline
- **No fsync**: a crash or power loss can lose the most recent writes
- **AOF grows forever**: there is no log compaction or rewrite
- **No authentication, TLS, replication, transactions, or pub/sub**
- **Global store**: the store and AOF are package-level variables, which keeps the code simple but makes tests share state
- **Inline commands** (plain `PING\r\n` typed into telnet) are not supported; only RESP arrays are



## License

MIT