## Install postgres, golang, golang/dep

## Download dataset

- visit https://data.cityofnewyork.us/Housing-Development/Building-Footprints/nqwf-w8eh
- Export using CSV format, about 500MB
- save it to `./data/building.csv`

## Setup postgres db.

- **NOTE: Must** create non-empty password for database user.
  Application default using:

```
DB_USERNAME=postgres
DB_PASSWORD=postgres
DB_NAME=nycbuilding
DB_HOST=localhost
DB_PORT=5432
DB_SSL_MODE=disable
```

## Install depencency library

- Run command:

```
dep ensure
```

## Build etl tool and server application

- Run commands:

```
./scripts/buildAll.sh
ls -l etl server
```

## ETL

- Modify the postgres connection config(DB\_\* variable) accordingly in `./data/etl.csv` if needed.
- Run with sample `etl.sh ./data/sample_buildings.csv` or full dataset `./scripts/etl.sh ./data/building.csv`.
  - Create db and tables
  - Run etl tool `./etl <data.csv>` extract and transform data, and save to `./data/etl.csv`
  - Load `./data/etl.csv` to db

```
./scripts/etl.sh ./data/sample_buildings.csv
CREATE DATABASE
CREATE TABLE
CREATE INDEX
CREATE INDEX
Processing begin
Process  499  records in 6.270295ms
COPY 499

ls -l ./data/etl.csv
```

## Start server

- Modify the postgres connection config(DB\_\* variable) and server port(PORT) in `./scripts/start-server.sh` if needed

```
DB_USER=postgres
DB_HOST=localhost
DB_PORT=5432
DB_NAME=nycbuilding
DB_SSL_MODE=disable

PORT=8080
```

- **Must** set environment variable DB_PASSWORD with database password, should be same to `etl.sh` db password (default: "postgres")

```
export DB_PASSWORD=postgres
./scripts/start-server.sh
```

OR

```
DB_PASSWORD=postgres ./scripts/start-server.sh
```

## Test queries

- Modify environment variable `SERVER` in `./scripts/query.sh` accordingly, default: `localhost:8080`
- Run

```
./scripts/query.sh
```
