# Tools and files to support polymesh deployments

## Grafana

The `grafana` directory contains a sample dashboard for monitoring Polymesh nodes. Sample alerts
are included and will be usable out of the box as long as Grafana has a default notification
channel configured.

## Docker

We have included sample `docker-compose` configuration files for operators and sentries.  These
asume that the operators and sentries are deployed to separate hosts, but can be modified to deploy
to a single host or to a pool of hosts (if using Swarm).


