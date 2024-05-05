#!/bin/bash

# Run migrations
goose -dir ./migrations sqlite3 "./${DB_URL}" up

# TODO: Ask for app port and create .env file

# Dev
# exec air -c .air.toml

# Prod
/ssh-client