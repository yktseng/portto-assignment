# Portto Assignment

## A eth block collector

## Prerequisites

Postgresql

* init docker volume `docker volume create portto_assignment`
* init postgres in docker  `docker run -e POSTGRES_USER=portto -p 5432:5432 -e POSTGRES_PASSWORD=portto -v portto_assignment:/var/lib/postgresql/data -d postgres`


To initialize the db schema,

`psql -p 5432 -U portto -W -h localhost -d portto -a -f ./scripts/create_table.sql`

## To run the program

`go run cmd/collector/collector.go` starts the indexer service

possible input parameters are

*  bWorkerSize: workers to collect blocks, default is 1
*  txWorkerSize: workers to collect transactions, default is 8

`go run cmd/webserver/webserver.go` starts the web server

## What's done

* API server: three apis
* Block indexer service

* Able to adjust block and transaction collector sizes
* Use websocket endpoint to listen to new headers
* It seems that BSC blocks does not contain uncles, so I just blindly see the recent 20 blocks as unconfirmed and fetch the blocks again after 20 blocks

## Perf tuning

Apple M1 Pro 2021, postgresql installed as a docker image

### 2 block collectors 

#### 16 tx collectors

6000 txs per minute

350 blocks

#### 48 tx collectors

16000 txs per minute

2000 blocks


### 4 block collectors 

#### 32 tx collectors

12000 txs per minute

750 blocks

### 8 block collectors 

#### 64 tx collectors

CPU 55%  memory 35mb

22000 txs per minutes

1600 blocks

kicked by bsc nodes...

