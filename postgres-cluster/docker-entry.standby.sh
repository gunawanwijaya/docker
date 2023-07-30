#!/bin/sh
set -e

# ----------------------------------------------------------------------------------------------------------------------
hostname="postgres-standby";
log_directory="/var/log/${hostname}";
# ----------------------------------------------------------------------------------------------------------------------
rm -rf   /data/${hostname} /tmp/${hostname} /var/log/${hostname};
mkdir -p /data/${hostname} /tmp/${hostname} /var/log/${hostname};
echo -n "" > /var/log/${hostname}/run.log;

# in order to start standby server, we need clone from primary server using `pg_basebackup`
# and `pg_basebackup` need a password prompt if the replication user in primary are set to have a password
# that's why we pipe the password into `pg_basebackup` STDIN. `pg_basebackup` will populate
# identical/mirror of ALL databases from primary that usually contains in ${PGDATA} directory.
# REF - https://www.postgresql.org/docs/current/app-pgbasebackup.html
docker-entrypoint.sh pg_basebackup \
    -D "${PGDATA}" \
    -h "postgres-primary" -p "5432" \
    -X "stream" \
    -c "fast" \
    -U "$(cat /run/secrets/pgunrp/pgunrp)" \
    -S "$(cat /run/secrets/pgslrp/pgslrp)" \
    -PRvW   < /run/secrets/pgpwrp/pgpwrp \
| tee /var/log/${hostname}/run.log;


# once the databases is populated, we start the services, here we're using `postgres`, because the ${PGDATA}
# directory is not empty the files from docker-entrypoint-initdb.d are not executed.
# REF - https://www.postgresql.org/docs/current/app-postgres.html
docker-entrypoint.sh postgres \
    -D "${PGDATA}" \
    -c "log_directory=${log_directory}" \
| tee /var/log/${hostname}/run.log;
