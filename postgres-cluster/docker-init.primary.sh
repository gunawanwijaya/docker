#!/bin/sh
set -e

# ----------------------------------------------------------------------------------------------------------------------
POSTGRES_REPLICA_SLOT="$(cat /run/secrets/pgslrp/pgslrp)";
POSTGRES_REPLICA_USER="$(cat /run/secrets/pgunrp/pgunrp)";
POSTGRES_REPLICA_PASSWORD="$(cat /run/secrets/pgpwrp/pgpwrp)";

# Worth to note that `${POSTGRES_REPLICA_SLOT}` and `${POSTGRES_REPLICA_USER}` should be in lowercase.
# when providing `${POSTGRES_REPLICA_SLOT}` with non compliant format (not in lowercase alphanumeric & underscore) will
# return the error message, but on CREATE USER `${POSTGRES_REPLICA_USER}` there is no error message and instead it will
# transform the value to the lowercase format.
psql -v "ON_ERROR_STOP=1" -U "$POSTGRES_USER" -d "$POSTGRES_DB" <<-EOSQL
	CREATE USER ${POSTGRES_REPLICA_USER} WITH REPLICATION
		ENCRYPTED PASSWORD '${POSTGRES_REPLICA_PASSWORD}';
	SELECT * FROM pg_create_physical_replication_slot('${POSTGRES_REPLICA_SLOT}');
	\du
EOSQL

LINE="host replication ${POSTGRES_REPLICA_USER} all trust";
grep -qxF "${LINE}" "${PGDATA}/pg_hba.conf" || echo "${LINE}" >> "${PGDATA}/pg_hba.conf";
# ----------------------------------------------------------------------------------------------------------------------
