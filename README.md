# Tools and files to support polymesh deployments

## Documentation

The operator runbook in `docs/operator` will guide you in setting up your operator node.

## Grafana

The `grafana` directory contains a sample dashboard for monitoring Polymesh nodes. Sample alerts
are included and will be usable out of the box as long as Grafana has a default notification
channel configured.

## Containers

We provide the appropriate source code and Dockerfiles in the `containers` directory to allow you to
reproduce our official images.

## Docker

We have included sample `docker-compose` configuration files for operators.  These
asume that the operators are deployed to separate hosts, but can be modified to deploy
to a single host or to a pool of hosts (if using Swarm).

## Kubernetes

The Helm chart published under `helm/polymesh` can be used to easily deploy operator or sentry nodes on a
Kubernetes cluster.

