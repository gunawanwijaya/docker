all:
	@cat ./Makefile;

init-secret:
	@sh ./.script/init-secret.sh

compose-minio:
	@UID=${UID} GID=${GID} bash -c "trap '\
		docker compose -f ./minio/docker-compose.yml \
		down --remove-orphans --volumes;' EXIT;\
		docker compose -f ./minio/docker-compose.yml \
		up --remove-orphans --build --force-recreate;";
compose-postgres-single:
	@UID=${UID} GID=${GID} bash -c "trap '\
		docker compose -f ./postgres/docker-compose.yml \
		--profile single down --remove-orphans --volumes;' EXIT;\
		docker compose -f ./postgres/docker-compose.yml \
		--profile single up --remove-orphans --build --force-recreate;";
compose-loki-all:
	@UID=${UID} GID=${GID} bash -c "trap '\
		docker compose -f ./grafana-loki/docker-compose.yml \
		--profile all down --remove-orphans --volumes;' EXIT;\
		docker compose -f ./grafana-loki/docker-compose.yml \
		--profile all up --remove-orphans --build --force-recreate;";
compose-tempo-all:
	@UID=${UID} GID=${GID} bash -c "trap '\
		docker compose -f ./grafana-tempo/docker-compose.yml \
		--profile all down --remove-orphans --volumes;' EXIT;\
		docker compose -f ./grafana-tempo/docker-compose.yml \
		--profile all up --remove-orphans --build --force-recreate;";
compose-mimir-all:
	@UID=${UID} GID=${GID} bash -c "trap '\
		docker compose -f ./grafana-mimir/docker-compose.yml \
		--profile all down --remove-orphans --volumes;' EXIT;\
		docker compose -f ./grafana-mimir/docker-compose.yml \
		--profile all up --remove-orphans --build --force-recreate;";
compose-postgres-ha:
	@UID=${UID} GID=${GID} bash -c "trap '\
		docker compose -f ./postgres/docker-compose.yml \
		--profile ha down --remove-orphans --volumes;' EXIT;\
		docker compose -f ./postgres/docker-compose.yml \
		--profile ha up --remove-orphans --build --force-recreate;";
compose-loki-cluster:
	@UID=${UID} GID=${GID} bash -c "trap '\
		docker compose -f ./grafana-loki/docker-compose.yml \
		--profile cluster down --remove-orphans --volumes;' EXIT;\
		docker compose -f ./grafana-loki/docker-compose.yml \
		--profile cluster up --remove-orphans --build --force-recreate;";
compose-tempo-cluster:
	@UID=${UID} GID=${GID} bash -c "trap '\
		docker compose -f ./grafana-tempo/docker-compose.yml \
		--profile cluster down --remove-orphans --volumes;' EXIT;\
		docker compose -f ./grafana-tempo/docker-compose.yml \
		--profile cluster up --remove-orphans --build --force-recreate;";
compose-mimir-cluster:
	@UID=${UID} GID=${GID} bash -c "trap '\
		docker compose -f ./grafana-mimir/docker-compose.yml \
		--profile cluster down --remove-orphans --volumes;' EXIT;\
		docker compose -f ./grafana-mimir/docker-compose.yml \
		--profile cluster up --remove-orphans --build --force-recreate;";
compose-grafana:
	@UID=${UID} GID=${GID} bash -c "trap '\
		docker compose -f ./grafana/docker-compose.yml \
		down --remove-orphans --volumes;' EXIT;\
		docker compose -f ./grafana/docker-compose.yml \
		up --remove-orphans --build --force-recreate;";
compose-otel-collector:
	@UID=${UID} GID=${GID} bash -c "trap '\
		docker compose -f ./otel-collector/docker-compose.yml \
		down --remove-orphans --volumes;' EXIT;\
		docker compose -f ./otel-collector/docker-compose.yml \
		up --remove-orphans --build --force-recreate;";
