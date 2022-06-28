# Portto Assignment

## A eth block collector

## Prerequisites

Postgresql

* init docker volume `docker volume create portto_assignment`
* init postgres in docker  `docker run -e POSTGRES_USER=portto -p 5432:5432 -e POSTGRES_PASSWORD=portto -v portto_assignment:/var/lib/postgresql/data -d postgres`


To initialize the db schema,

`psql -p 5432 -U portto -W -h localhost -d portto -a -f ./scripts/create_table.sql`