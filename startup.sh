#!/bin/bash

# Run migrations
goose -dir ./migrations sqlite3 "./${DB_URL}" up

# Start the Go application
exec air -c .air.toml