#!/bin/sh
GOBLET_BINARY=/app/goblet

# Checking if logstash monitoring required
if [ -n "${LOGSTASH_ENABLED+set}" ]; then
    if [[ "$LOGSTASH_ENABLED" == "true" ]]; then
        echo "Waiting for Logstash/Elasticsearch to start..."
        while ! nc -z ${LOGSTASH_HOST} ${LOGSTASH_PORT}; do
          sleep 0.1
        done
        echo "Logstash online, continuing boot."
    fi
fi

# Waiting for Postgres to start
while ! nc -z ${DB_HOST} ${DB_PORT}; do
  sleep 0.1
done
echo "Postgres online, continuing boot."

#echo "Creating default database"
${GOBLET_BINARY} dbcreate

#echo "Running database migrations"
${GOBLET_BINARY} dbmigrate

#echo "Seeding database with initial data"
${GOBLET_BINARY} dbseed

#echo "Launching Goblet server"
${GOBLET_BINARY} serve
