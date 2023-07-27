#!/bin/sh
set -e

# ----------------------------------------------------------------------------------------------------------------------
target="$1";
suffix="$2";

s3_host=$(cat /run/secrets/s3host/s3host);
s3_port=$(cat /run/secrets/s3port/s3port);
s3_region=$(cat /run/secrets/s3region/s3region);
s3_bucket=$(cat /run/secrets/s3bucket/s3bucket);
s3_access=$(cat /run/secrets/s3access/s3access);
s3_secret=$(cat /run/secrets/s3secret/s3secret);
# ----------------------------------------------------------------------------------------------------------------------
tmp="/tmp/tempo/tempo-${target}${suffix}";
cp /etc/config/tempo.yml /etc/config/tempo.TMP.yml;
sed -i "s|/tmp/tempo/metrics_generator.storage.path|${tmp}/metrics|"                        /etc/config/tempo.TMP.yml;
sed -i "s|grpc: { endpoint: \"\" }|grpc: { endpoint: \"tempo-${target}${suffix}:4317\" }|"  /etc/config/tempo.TMP.yml;
mkdir -p /data/tempo /tmp/tempo /var/log/tempo;
/tempo -target="${target}" \
    -config.file="/etc/config/tempo.TMP.yml" \
    -storage.trace.s3.endpoint="${s3_host}:${s3_port}" \
    -storage.trace.s3.bucket="${s3_bucket}" \
    -storage.trace.s3.access_key="${s3_access}" \
    -storage.trace.s3.secret_key="${s3_secret}" \
    -storage.trace.local.path="${tmp}/traces" \
    -storage.trace.wal.path="${tmp}/wal" \
 2>&1 | tee "/var/log/tempo/tempo-${target}${suffix}.log"
