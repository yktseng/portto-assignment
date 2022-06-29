# Portto Assignment

## A eth block collector

## Prerequisites

Postgresql

* init docker volume `docker volume create portto_assignment`
* init postgres in docker  `docker run -e POSTGRES_USER=portto -p 5432:5432 -e POSTGRES_PASSWORD=portto -v portto_assignment:/var/lib/postgresql/data -d postgres`


To initialize the db schema,

`psql -p 5432 -U portto -W -h localhost -d portto -a -f ./scripts/create_table.sql`

## Perf tuning

2 block collectors and 16 tx collectors

6000 txs per minute
350 blocks

4 block collectors and 32 tx collectors

12000 txs per minute
750 blocks

8 block collectors and 64 tx collectors

CPU 55%
22000 txs per minutes
1600 blocks

