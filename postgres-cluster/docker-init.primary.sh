#!/bin/sh
set -e

# ----------------------------------------------------------------------------------------------------------------------
POSTGRES_REPLICA_SLOT="$(cat /run/secrets/pgslrp/pgslrp)";
POSTGRES_REPLICA_USER="$(cat /run/secrets/pgunrp/pgunrp)";
POSTGRES_REPLICA_PASSWORD="$(cat /run/secrets/pgpwrp/pgpwrp)";
POSTGRES_READONLY_USER="$(cat /run/secrets/pgunro/pgunro)";
POSTGRES_READONLY_PASSWORD="$(cat /run/secrets/pgpwro/pgpwro)";

# Worth to note that `${POSTGRES_REPLICA_SLOT}` and `${POSTGRES_REPLICA_USER}` should be in lowercase.
# but when running CREATE USER `${POSTGRES_REPLICA_USER}` there are no errors, instead it will
# transform the value to the lowercase format.
psql -v "ON_ERROR_STOP=1" -U "$POSTGRES_USER" -d "$POSTGRES_DB" <<-EOSQL
	-----------------------------------------------------------------------------------------------
	-- REF: https://www.keyvanfatehi.com/2021/07/14/how-to-create-read-only-user-in-postgresql/
	CREATE ROLE __readonly;
	GRANT USAGE ON SCHEMA public TO __readonly;
	GRANT SELECT ON ALL TABLES IN SCHEMA public TO __readonly;
	ALTER DEFAULT PRIVILEGES IN SCHEMA public
		GRANT SELECT ON TABLES TO __readonly;
	CREATE USER ${POSTGRES_READONLY_USER} WITH
		ENCRYPTED PASSWORD '${POSTGRES_READONLY_PASSWORD}';
	GRANT __readonly TO ${POSTGRES_READONLY_USER};
	-----------------------------------------------------------------------------------------------
	CREATE USER ${POSTGRES_REPLICA_USER} WITH REPLICATION
		ENCRYPTED PASSWORD '${POSTGRES_REPLICA_PASSWORD}';
	SELECT * FROM pg_create_physical_replication_slot('${POSTGRES_REPLICA_SLOT}');
	-----------------------------------------------------------------------------------------------
	CREATE USER _grafana WITH
		ENCRYPTED PASSWORD '_grafana';
	CREATE DATABASE _grafana OWNER _grafana;
	-----------------------------------------------------------------------------------------------
	CREATE USER _sonar WITH
		ENCRYPTED PASSWORD '_sonar';
	CREATE DATABASE _sonar OWNER _sonar;
	-----------------------------------------------------------------------------------------------
	\du
EOSQL

LINE="host replication ${POSTGRES_REPLICA_USER} all trust";
grep -qxF "${LINE}" "${PGDATA}/pg_hba.conf" || echo "${LINE}" >> "${PGDATA}/pg_hba.conf";
# ----------------------------------------------------------------------------------------------------------------------
