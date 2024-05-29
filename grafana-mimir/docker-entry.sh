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
tmp="/tmp/mimir/mimir-${target}${suffix}";
mkdir -p /data/mimir/ /tmp/mimir /var/log/mimir;
mkdir -p ${tmp}/rules ${tmp}/ruler ${tmp}/alertmanager;

query_frontend="query-frontend"
query_scheduler="query-scheduler"
memberlist_join="ingester-0"
if [ "${target}" = "all" ]; then
    query_frontend="all"
    query_scheduler="all"
    memberlist_join="all"
fi

echo "ok=${1}${2}"
# mimir -config.file="/etc/config/mimir.yml" -version;
mimir -config.file="/etc/config/mimir.yml" -target="${target}" \
    -common.storage.s3.endpoint="${s3_host}:${s3_port}" \
    -common.storage.s3.bucket-name="${s3_bucket}" \
    -common.storage.s3.access-key-id="${s3_access}" \
    -common.storage.s3.secret-access-key="${s3_secret}" \
    -ruler.rule-path="${tmp}/rules" \
    -ruler-storage.local.directory="${tmp}/ruler" \
    -alertmanager.storage.path="${tmp}/alertmanager" \
    -blocks-storage.bucket-store.sync-dir="${tmp}/tsdb_sync" \
    -blocks-storage.tsdb.dir="${tmp}/tsdb" \
    -compactor.data-dir="${tmp}/compactor" \
    -activity-tracker.filepath="${tmp}.activity-tracker.log" \
    -querier.frontend-address="mimir-${query_frontend}:3301" \
    -memberlist.join="mimir-${memberlist_join}:7946" \
 2>&1 | tee "/var/log/mimir/mimir-${target}${suffix}.log";

# -querier.scheduler-address="mimir-${query_scheduler}:3301" \
# -ingester.ring.tokens-file-path \
# -alertmanager-storage.local.path \
# -store-gateway.sharding-ring.tokens-file-path
