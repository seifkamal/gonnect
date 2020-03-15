# Gonnect

A simple online game matchmaking server.

## Features

You can run any command or subcommand with the `help` flag to see descriptions, examples, and acceptable flags.

Here's a summary of the available features:

### Server

A server can be started by running the `serve` command, and specifying a handler. Currently the supported handlers are:
- `player`
- `match`

**Example:**
```shell script
> gonnect serve match --port :8080 -u admin -p honeyisgood
```

Some endpoints require basic authorisation; See the `help` print for this command for information on how to change
the default credentials.

#### `player`

This will expose a `GET player/match` endpoint; Requests to this endpoint will be upgraded to a WebSocket
connection. The server will attempt to find a match for the player as long as the connection is maintained.
Once found, the match data will be returned in a JSON response and the connection will be closed.

#### `match`

This will expose the following endpoints:
- `GET match/all?state=ready` - Returns all matches found matching the specified `state`
- `GET match/{matchId}` - Returns a match with the specified `matchId`
- `POST match/{matchId}/end` - Sets the state of the match with the specified `matchId` to `ended`

### Worker

A matchmaking worker can be started by running the `match` command. A match size (ie. player count) can be set
via the `batch` flag; The default is 10.

**Example:**
```shell script
> gonnect match --batch 5
```

This will run an ongoing process that will create a new match whenever enough players are searching.
