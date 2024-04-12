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

s3_url="http://${s3_access}:${s3_secret}@${s3_host}.${s3_region}:${s3_port}/${s3_bucket}";
# ----------------------------------------------------------------------------------------------------------------------
tmp="/tmp/loki/loki-${target}${suffix}";
mkdir -p /data/loki /tmp/loki /var/log/loki ${tmp};

querier="querier-0"
query_frontend="query-frontend"
query_scheduler="query-scheduler"
memberlist_join="ingester-0"
compactor="compactor"
if [ ${target} = "all" ]; then
    querier="all"
    query_frontend="all"
    query_scheduler="all"
    memberlist_join="all"
    compactor="all"
fi

echo "ok=${1}${2}"
# /usr/bin/loki -config.file="/etc/config/loki.yml" -list-targets;
# /usr/bin/loki -config.file="/etc/config/loki.yml" -version;
/usr/bin/loki -config.file="/etc/config/loki.yml" -target="${target}" \
    -common.storage.s3.url="${s3_url}" \
    -ruler.rule-path="${tmp}/rules" \
    -ruler.storage.local.directory="${tmp}/ruler_storage" \
    -ruler.wal.dir="${tmp}/ruler_wal" \
    -ingester.wal-dir="${tmp}/ingester_wal" \
    -boltdb.dir="${tmp}/boltdb" \
    -boltdb.shipper.active-index-directory="${tmp}/boltdb_index" \
    -boltdb.shipper.cache-location="${tmp}/boltdb_cache" \
    -boltdb.shipper.compactor.working-directory="${tmp}/boltdb_compactor" \
    -tsdb.shipper.active-index-directory="${tmp}/tsdb_index" \
    -tsdb.shipper.cache-location="${tmp}/tsdb_cache" \
    -local.chunk-directory="${tmp}/local_chunk" \
    -frontend.tail-proxy-url="http://loki-${querier}:3100" \
    -querier.frontend-address="loki-${query_frontend}:3101" \
    -common.compactor-address="http://loki-${compactor}:3100" \
    -common.compactor-grpc-address="loki-${compactor}:3101" \
    -memberlist.join="loki-${memberlist_join}:7946" \
 2>&1 | tee "/var/log/loki/loki-${target}${suffix}.log";

# -querier.scheduler-address="loki-${query_scheduler}:3101" \
# -query-scheduler.ring.tokens-file-path \
# -ingester.tokens-file-path \
# -index-gateway.ring.tokens-file-path \
# -boltdb.shipper.compactor.ring.tokens-file-path \
# -common.storage.ring.tokens-file-path \
