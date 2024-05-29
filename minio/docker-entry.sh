#!/bin/sh
set -e

# ----------------------------------------------------------------------------------------------------------------------
chmod +x "/bin/mc";
bg() {
    sleep .5;
    local ALIAS="this";
    mc alias set ${ALIAS} \
        "http://localhost:9000" \
        "$(cat ${MINIO_ROOT_USER_FILE})" \
        "$(cat ${MINIO_ROOT_PASSWORD_FILE})" \
        > /dev/null;

    local BUCKET=$(cat /run/secrets/loki-bucket/loki-bucket)
    local ACCESS=$(cat /run/secrets/loki-access/loki-access)
    local SECRET=$(cat /run/secrets/loki-secret/loki-secret)
    mc mb "${ALIAS}/${BUCKET}"                                   > /dev/null;
    mc admin user add ${ALIAS} "${ACCESS}" "${SECRET}"           > /dev/null;
    mc admin policy attach ${ALIAS} readwrite --user "${ACCESS}" > /dev/null;

    local BUCKET=$(cat /run/secrets/tempo-bucket/tempo-bucket)
    local ACCESS=$(cat /run/secrets/tempo-access/tempo-access)
    local SECRET=$(cat /run/secrets/tempo-secret/tempo-secret)
    mc mb "${ALIAS}/${BUCKET}"                                   > /dev/null;
    mc admin user add ${ALIAS} "${ACCESS}" "${SECRET}"           > /dev/null;
    mc admin policy attach ${ALIAS} readwrite --user "${ACCESS}" > /dev/null;

    local BUCKET=$(cat /run/secrets/mimir-bucket/mimir-bucket)
    local ACCESS=$(cat /run/secrets/mimir-access/mimir-access)
    local SECRET=$(cat /run/secrets/mimir-secret/mimir-secret)
    mc mb "${ALIAS}/${BUCKET}"                                   > /dev/null;
    mc admin user add ${ALIAS} "${ACCESS}" "${SECRET}"           > /dev/null;
    mc admin policy attach ${ALIAS} readwrite --user "${ACCESS}" > /dev/null;

    local BUCKET=$(cat /run/secrets/pyroscope-bucket/pyroscope-bucket)
    local ACCESS=$(cat /run/secrets/pyroscope-access/pyroscope-access)
    local SECRET=$(cat /run/secrets/pyroscope-secret/pyroscope-secret)
    mc mb "${ALIAS}/${BUCKET}"                                   > /dev/null;
    mc admin user add ${ALIAS} "${ACCESS}" "${SECRET}"           > /dev/null;
    mc admin policy attach ${ALIAS} readwrite --user "${ACCESS}" > /dev/null;

    mc admin prometheus generate ${ALIAS};
}; bg &
# ----------------------------------------------------------------------------------------------------------------------
rm -rf /data/minio;
minio server "/data/minio/disk{1...4}" \
    --json --anonymous \
    --address=":9000" \
    --console-address=":9001"
