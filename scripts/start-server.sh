#!/bin/bash

DB_USER=postgres
DB_HOST=localhost
DB_PORT=5432
DB_NAME=nycbuilding
DB_SSL_MODE=disable

PORT=8080

if [ -z "$DB_PASSWORD" ]; then
    echo "should 'export DB_PASSWORD=<password>' or 'DB_PASSWORD=<password> ${BASH_SOURCE[0]}'"
    echo "current using default password: postgres"
    DB_PASSWORD=postgres ./server
else
    ./server
fi


