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

s3_url="http://${s3_access}:${s3_secret}@${s3_host}.${s3_region}:${s3_port}/${s3_bucket}";
# ----------------------------------------------------------------------------------------------------------------------
tmp="/tmp/loki/loki-${target}${suffix}";
mkdir -p        /data/loki /tmp/loki /var/log/loki ${tmp};
# chown -R loki   /data/loki /tmp/loki /var/log/loki ${tmp};
# sed -i "s|/home/loki:/sbin/nologin|/home/loki:/bin/sh|" /etc/passwd;
# echo "A ------------"
# timeout 1 sh -c 'cat < /dev/null > /dev/tcp/127.0.0.1/3100'; echo $?;
# echo "B ------------"
# su -c "timeout 1 sh -c 'cat < /dev/null > /dev/tcp/127.0.0.1/3100'; echo $?" loki;
# echo "C ------------"
# exit 1

/usr/bin/loki -target="${target}" \
    -config.file="/etc/config/loki.yml" \
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
 2>&1 | tee "/var/log/loki/loki-${target}${suffix}.log";
# -query-scheduler.ring.tokens-file-path \
# -ingester.tokens-file-path \
# -index-gateway.ring.tokens-file-path \
# -boltdb.shipper.compactor.ring.tokens-file-path \
# -common.storage.ring.tokens-file-path \
