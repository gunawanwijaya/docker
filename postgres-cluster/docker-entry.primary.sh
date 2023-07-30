#!/bin/sh
set -e

# ----------------------------------------------------------------------------------------------------------------------
hostname="postgres-primary";
log_directory="/var/log/${hostname}";
config_file="/etc/config/postgres.conf";
# ----------------------------------------------------------------------------------------------------------------------
rm -rf                  /data/${hostname} /tmp/${hostname} /var/log/${hostname};
mkdir -p                /data/${hostname} /tmp/${hostname} /var/log/${hostname};
chown postgres:postgres /data/${hostname} /tmp/${hostname} /var/log/${hostname};
echo -n "" > /var/log/${hostname}/run.log;


# once the databases is populated, we start the services, here we're using `postgres`, because the $${PGDATA}
# directory is not empty the files from docker-entrypoint-initdb.d are not executed.
# REF - https://www.postgresql.org/docs/current/app-postgres.html
docker-entrypoint.sh postgres \
    -D "${PGDATA}" \
    -c "config_file=${config_file}" \
    -c "log_directory=${log_directory}" \
    -c "log_timezone=${TZ}" \
    -c "timezone=${TZ}" \
| tee /var/log/${hostname}/run.log;
