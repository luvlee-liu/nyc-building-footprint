#!/bin/bash
set -e

# create database and table
SCRIPT_DIR=$(dirname "${BASH_SOURCE[0]}")
PROJECT_DIR=$(realpath "$SCRIPT_DIR/..")
SQL_DIR="$PROJECT_DIR/sql" 
DATA_DIR="$PROJECT_DIR/data"

DB_USERNAME=postgres
DB_PASSWORD=postgres
DB_NAME=nycbuilding
DB_HOST=localhost
DB_PORT=5432
DB_SSL_MODE=disable
DB_CONFIG="-U ${DB_USERNAME} -h ${DB_HOST} -p ${DB_PORT}"

# create db
PGPASSWORD=${DB_PASSWORD} psql ${DB_CONFIG} -c "CREATE DATABASE ${DB_NAME};"

# create table
PGPASSWORD=${DB_PASSWORD} psql ${DB_CONFIG} -d ${DB_NAME} -f ${SQL_DIR}/create.sql

# extract column and transform
DATA_RAW_CSV="$1"
DATA_ETL_CSV="${DATA_DIR}/etl.csv"
if [ -z "$DATA_RAW_CSV" ]; then
    DATA_RAW_CSV="${DATA_DIR}/building.csv"
    echo "default use data from ${DATA_RAW_CSV}"
fi

${PROJECT_DIR}/etl -input ${DATA_RAW_CSV} -output ${DATA_ETL_CSV}

# copy csv to db
BULDING_CSV=\'$(realpath ${DATA_ETL_CSV})\'
PGPASSWORD=${DB_PASSWORD} psql ${DB_CONFIG} -d ${DB_NAME} -f "${SQL_DIR}/copy.sql" -v building_csv=${BULDING_CSV}

exit 0
