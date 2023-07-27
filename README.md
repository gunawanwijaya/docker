# docker
Docker related collection

## Grafana Loki Cluster Mode

Read this [config][gr_loki_1] to familiarize with all of the Loki components, after that read this [Deployment modes][gr_loki_2] for more details what approach that we want to choose.

### Targets & [Components][gr_loki_3] [v2.8.x]
- `target=all` meaning that those instance are running as monolithic, doing all read & write operations this is really simple to use and suitable for development
- `target=read` & `target=write` combine multiple components but separating the read & write operations, combined with load balancer like nginx, we can have multiple read & multiple write components
- `target=compactor`, `target=distributor`, `target=ingester`, `target=querier`, `target=query-scheduler`, `target=ingester-querier`, `target=query-frontend`, `target=index-gateway`, `target=ruler`, `target=table-manager` these are individual components that can be tuned for more scalability.


<!-- LINKS -->
[gr_loki_1]: https://grafana.com/docs/loki/latest/configuration/#supported-contents-and-default-values-of-lokiyaml
[gr_loki_2]: https://grafana.com/docs/loki/latest/fundamentals/architecture/deployment-modes
[gr_loki_3]: https://grafana.com/docs/loki/latest/fundamentals/architecture/components
