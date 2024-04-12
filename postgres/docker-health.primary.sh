#!/bin/sh
set -e

# ----------------------------------------------------------------------------------------------------------------------
hostname="postgres-primary";
echo "healthcheck starting" | tee /var/log/${hostname}/run.log;
pg_isready -q \
    -h "$(cat ${hostname})" -p "5432" \
    -d "$(cat ${POSTGRES_DB_FILE})" \
    -U "$(cat ${POSTGRES_USER_FILE})" \
;
