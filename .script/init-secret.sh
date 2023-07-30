#!/bin/sh
set -e

# ----------------------------------------------------------------------------------------------------------------------
exitNotInDir(){
    echo "make sure in root directory and call me using $ .script/init-secret.sh;";
    echo "currently in ${PWD}";
    exit 1;
}
msgChangeDir(){
    echo "adjusting directory to ${PWD}";
}
rand(){
    LEN=$1;
    cat /dev/urandom | tr -dc '[:alpha:]' | fold -w${LEN} | head -n 1;
}
# ----------------------------------------------------------------------------------------------------------------------
this="init-secret.sh";
trap 'echo "${this} is done"' EXIT;
echo "${this} is starting";
# ----------------------------------------------------------------------------------------------------------------------
# directory adjustment
[ $(dirname $(dirname $0)) != "." ] && cd $(dirname $(dirname $0)) && msgChangeDir;
[ -x "${this}" ] && [ "$(cat "${this}")" = "$(cat $0)" ] && cd .. && msgChangeDir;
# ----------------------------------------------------------------------------------------------------------------------
[ -d "./grafana-loki-cluster" ] || exitNotInDir;
DIR="./grafana-loki-cluster/.secret";
mkdir -p "${DIR}";
[ ! -f "${DIR}/.s3host" ] && echo -n "minio"               > "${DIR}/.s3host"     && echo "${DIR}/.s3host created";
[ ! -f "${DIR}/.s3port" ] && echo -n "9000"                > "${DIR}/.s3port"     && echo "${DIR}/.s3port created";
[ ! -f "${DIR}/.s3region" ] && echo -n ""                  > "${DIR}/.s3region"   && echo "${DIR}/.s3region created";
[ ! -f "${DIR}/.s3bucket" ] && echo -n "gr-loki-bucket"    > "${DIR}/.s3bucket"   && echo "${DIR}/.s3bucket created";
[ ! -f "${DIR}/.s3access" ] && echo -n "gr-loki"           > "${DIR}/.s3access"   && echo "${DIR}/.s3access created";
[ ! -f "${DIR}/.s3secret" ] && echo -n "$(rand 40)"        > "${DIR}/.s3secret"   && echo "${DIR}/.s3secret created";
# ----------------------------------------------------------------------------------------------------------------------
[ -d "./grafana-mimir-cluster" ] || exitNotInDir;
DIR="./grafana-mimir-cluster/.secret";
mkdir -p "${DIR}";
[ ! -f "${DIR}/.s3host" ] && echo -n "minio"               > "${DIR}/.s3host"     && echo "${DIR}/.s3host created";
[ ! -f "${DIR}/.s3port" ] && echo -n "9000"                > "${DIR}/.s3port"     && echo "${DIR}/.s3port created";
[ ! -f "${DIR}/.s3region" ] && echo -n ""                  > "${DIR}/.s3region"   && echo "${DIR}/.s3region created";
[ ! -f "${DIR}/.s3bucket" ] && echo -n "gr-mimir-bucket"   > "${DIR}/.s3bucket"   && echo "${DIR}/.s3bucket created";
[ ! -f "${DIR}/.s3access" ] && echo -n "gr-mimir"          > "${DIR}/.s3access"   && echo "${DIR}/.s3access created";
[ ! -f "${DIR}/.s3secret" ] && echo -n "$(rand 40)"        > "${DIR}/.s3secret"   && echo "${DIR}/.s3secret created";
# ----------------------------------------------------------------------------------------------------------------------
[ -d "./grafana-tempo-cluster" ] || exitNotInDir;
DIR="./grafana-tempo-cluster/.secret";
mkdir -p "${DIR}";
[ ! -f "${DIR}/.s3host" ] && echo -n "minio"               > "${DIR}/.s3host"     && echo "${DIR}/.s3host created";
[ ! -f "${DIR}/.s3port" ] && echo -n "9000"                > "${DIR}/.s3port"     && echo "${DIR}/.s3port created";
[ ! -f "${DIR}/.s3region" ] && echo -n ""                  > "${DIR}/.s3region"   && echo "${DIR}/.s3region created";
[ ! -f "${DIR}/.s3bucket" ] && echo -n "gr-tempo-bucket"   > "${DIR}/.s3bucket"   && echo "${DIR}/.s3bucket created";
[ ! -f "${DIR}/.s3access" ] && echo -n "gr-tempo"          > "${DIR}/.s3access"   && echo "${DIR}/.s3access created";
[ ! -f "${DIR}/.s3secret" ] && echo -n "$(rand 40)"        > "${DIR}/.s3secret"   && echo "${DIR}/.s3secret created";
# ----------------------------------------------------------------------------------------------------------------------
[ -d "./minio" ] || exitNotInDir;
DIR="./minio/.secret";
mkdir -p "${DIR}";
[ ! -f "${DIR}/.uname" ] && echo -n "$(rand 40)"           > "${DIR}/.uname"  && echo "${DIR}/.uname created";
[ ! -f "${DIR}/.paswd" ] && echo -n "$(rand 40)"           > "${DIR}/.paswd"  && echo "${DIR}/.paswd created";
# ----------------------------------------------------------------------------------------------------------------------
[ -d "./postgres-cluster" ] || exitNotInDir;
DIR="./postgres-cluster/.secret";
mkdir -p "${DIR}";
[ ! -f "${DIR}/.pgreplication-slot" ] && echo -n "$(rand 40)"      > "${DIR}/.pgreplication-slot"  && echo "${DIR}/.uname created";
[ ! -f "${DIR}/.pgreplication-username" ] && echo -n "$(rand 40)"  > "${DIR}/.pgreplication-username"  && echo "${DIR}/.uname created";
[ ! -f "${DIR}/.pgreplication-password" ] && echo -n "$(rand 40)"  > "${DIR}/.pgreplication-password"  && echo "${DIR}/.uname created";
[ ! -f "${DIR}/.pg-password" ] && echo -n "$(rand 40)"             > "${DIR}/.pg-password"  && echo "${DIR}/.paswd created";
# ----------------------------------------------------------------------------------------------------------------------
exit 0;
