#!/bin/sh
set -e

# ----------------------------------------------------------------------------------------------------------------------
target=${1};
suffix=${2};

s3_host=$(cat /run/secrets/s3host/s3host);
s3_port=$(cat /run/secrets/s3port/s3port);
s3_region=$(cat /run/secrets/s3region/s3region);
s3_bucket=$(cat /run/secrets/s3bucket/s3bucket);
s3_access=$(cat /run/secrets/s3access/s3access);
s3_secret=$(cat /run/secrets/s3secret/s3secret);
# ----------------------------------------------------------------------------------------------------------------------
mkdir -p /data/pyroscope /tmp/pyroscope /var/log/pyroscope /data/pyroscope/pyroscope-${target}${suffix};

query_frontend="query-frontend"
query_scheduler="query-scheduler"
memberlist_join="ingester-0"
if [ ${target} = "all" ]; then
    query_frontend="all"
    query_scheduler="all"
    memberlist_join="all"
fi

echo "ok=${1}${2}"


/usr/bin/pyroscope -config.file="/etc/config/pyroscope.yml" -target="${target}" \
    -storage.s3.endpoint="${s3_host}:${s3_port}" \
    -storage.s3.region="${s3_region}" \
    -storage.s3.bucket-name="${s3_bucket}" \
    -storage.s3.access-key-id="${s3_access}" \
    -storage.s3.secret-access-key="${s3_secret}" \
    -pyroscopedb.data-path="/data/pyroscope/pyroscope-${target}${suffix}" \
    -memberlist.join="pyroscope-${memberlist_join}:7946" \
 2>&1 | tee "/var/log/pyroscope/pyroscope-${target}${suffix}.log";
