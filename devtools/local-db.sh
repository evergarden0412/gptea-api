#!/bin/bash

PORT=5432
IMAGENAME=gptea
DBNAME=gptea
[[ -z "$1" ]] || PORT=$1

echo run postgres on port $PORT
[[ -z "$(podman ps -aq -f name=$IMAGENAME)" ]] || podman stop $IMAGENAME |> /dev/null

podman run --name $IMAGENAME --rm -d -e POSTGRES_PASSWORD=password -v $(pwd)/devtools:/devtools  -p $PORT:5432 docker.io/library/postgres:15.1-alpine;

echo waiting for podman database to start
sleep 3
podman exec -it $IMAGENAME psql -U postgres -c "create database $DBNAME";
echo waiting for podman database to start
podman exec -it $IMAGENAME psql -U postgres -d $DBNAME -a -q -f /devtools/db.sql;