all:
	@cat ./Makefile;

init-secret:
	@sh ./.script/init-secret.sh

compose-minio:
	@bash -c "trap '\
		docker compose -f ./minio/docker-compose.yml \
		down --remove-orphans --volumes;' EXIT;\
		docker compose -f ./minio/docker-compose.yml \
		up --remove-orphans --build;";
compose-loki-cluster:
	@bash -c "trap '\
		docker compose -f ./grafana-loki-cluster/docker-compose.yml \
		down --remove-orphans --volumes;' EXIT;\
		docker compose -f ./grafana-loki-cluster/docker-compose.yml \
		up --remove-orphans --build;";
compose-tempo-cluster:
	@bash -c "trap '\
		docker compose -f ./grafana-tempo-cluster/docker-compose.yml \
		down --remove-orphans --volumes;' EXIT;\
		docker compose -f ./grafana-tempo-cluster/docker-compose.yml \
		up --remove-orphans --build;";
compose-mimir-cluster:
	@bash -c "trap '\
		docker compose -f ./grafana-mimir-cluster/docker-compose.yml \
		down --remove-orphans --volumes;' EXIT;\
		docker compose -f ./grafana-mimir-cluster/docker-compose.yml \
		up --remove-orphans --build;";
compose-postgres-cluster:
	@bash -c "trap '\
		docker compose -f ./postgres-cluster/docker-compose.yml \
		down --remove-orphans --volumes;' EXIT;\
		docker compose -f ./postgres-cluster/docker-compose.yml \
		up --remove-orphans --build;";
