#!/bin/sh
set -e
panic(){
    echo "make sure in root directory, currently in ${PWD}";
    exit 1;
}
rand(){
    LEN=$1;
    cat /dev/urandom | tr -dc '[:alpha:]' | fold -w${LEN} | head -n 1;
}
# ----------------------------------------------------------------------------------------------------------------------
[ -d "./grafana-loki-cluster" ] || panic;
DIR="./grafana-loki-cluster/.secret";
mkdir -p "${DIR}";
test ! -f "${DIR}/.s3host" && echo "minio"               > "${DIR}/.s3host"     && echo "${DIR}/.s3host created";
test ! -f "${DIR}/.s3port" && echo "9000"                > "${DIR}/.s3port"     && echo "${DIR}/.s3port created";
test ! -f "${DIR}/.s3region" && echo ""                  > "${DIR}/.s3region"   && echo "${DIR}/.s3region created";
test ! -f "${DIR}/.s3bucket" && echo "gr-loki-bucket"    > "${DIR}/.s3bucket"   && echo "${DIR}/.s3bucket created";
test ! -f "${DIR}/.s3access" && echo "gr-loki"           > "${DIR}/.s3access"   && echo "${DIR}/.s3access created";
test ! -f "${DIR}/.s3secret" && echo $(rand 40)          > "${DIR}/.s3secret"   && echo "${DIR}/.s3secret created";
# ----------------------------------------------------------------------------------------------------------------------
[ -d "./grafana-mimir-cluster" ] || panic;
DIR="./grafana-mimir-cluster/.secret";
mkdir -p "${DIR}";
test ! -f "${DIR}/.s3host" && echo "minio"               > "${DIR}/.s3host"     && echo "${DIR}/.s3host created";
test ! -f "${DIR}/.s3port" && echo "9000"                > "${DIR}/.s3port"     && echo "${DIR}/.s3port created";
test ! -f "${DIR}/.s3region" && echo ""                  > "${DIR}/.s3region"   && echo "${DIR}/.s3region created";
test ! -f "${DIR}/.s3bucket" && echo "gr-mimir-bucket"   > "${DIR}/.s3bucket"   && echo "${DIR}/.s3bucket created";
test ! -f "${DIR}/.s3access" && echo "gr-mimir"          > "${DIR}/.s3access"   && echo "${DIR}/.s3access created";
test ! -f "${DIR}/.s3secret" && echo $(rand 40)          > "${DIR}/.s3secret"   && echo "${DIR}/.s3secret created";
# ----------------------------------------------------------------------------------------------------------------------
[ -d "./grafana-tempo-cluster" ] || panic;
DIR="./grafana-tempo-cluster/.secret";
mkdir -p "${DIR}";
test ! -f "${DIR}/.s3host" && echo "minio"               > "${DIR}/.s3host"     && echo "${DIR}/.s3host created";
test ! -f "${DIR}/.s3port" && echo "9000"                > "${DIR}/.s3port"     && echo "${DIR}/.s3port created";
test ! -f "${DIR}/.s3region" && echo ""                  > "${DIR}/.s3region"   && echo "${DIR}/.s3region created";
test ! -f "${DIR}/.s3bucket" && echo "gr-tempo-bucket"   > "${DIR}/.s3bucket"   && echo "${DIR}/.s3bucket created";
test ! -f "${DIR}/.s3access" && echo "gr-tempo"          > "${DIR}/.s3access"   && echo "${DIR}/.s3access created";
test ! -f "${DIR}/.s3secret" && echo $(rand 40)          > "${DIR}/.s3secret"   && echo "${DIR}/.s3secret created";
# ----------------------------------------------------------------------------------------------------------------------
[ -d "./minio" ] || panic;
DIR="./minio/.secret";
mkdir -p "${DIR}";
test ! -f "${DIR}/.uname" && echo $(rand 40)             > "${DIR}/.uname"  && echo "${DIR}/.uname created";
test ! -f "${DIR}/.paswd" && echo $(rand 40)             > "${DIR}/.paswd"  && echo "${DIR}/.paswd created";
# ----------------------------------------------------------------------------------------------------------------------
exit 0;
