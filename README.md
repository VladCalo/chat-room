# Chat Room

TCP chat room server written in Go. Clients connect via netcat or telnet, join rooms, and chat with others in the same room.

## How it works

The server listens on a TCP port. When a client connects, they enter their name and land in the lobby. From there they can join rooms, send messages, and see who else is online. Messages only go to people in the same room.

## Running the server

```bash
go run ./cmd/server
```

By default it listens on `127.0.0.1:9000`.

## Connecting as a client

```bash
nc 127.0.0.1 9000
```

## Commands

- `/join <room>` - Join a room (creates it if it doesnot exist)
- `/exit` - Leave the current room
- `/list` - Show all rooms
- `/list <room>` - Show members in a room
- `/whereAmI` - Show which room you are in
- `/help` - Show available commands

Any text that doesn't start with `/` gets sent to everyone in your current room.

## Notes

- Graceful shutdown with Ctrl+C
- Thread-safe with mutexes for concurrent access
- Clients are removed from rooms on disconnect
