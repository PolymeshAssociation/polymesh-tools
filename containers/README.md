# polymesh containers

The Dockerfiles and source code listed here can be used to reproduce the
container images that we use in the Helm chart published in this repository.

## The `polymesh` container

Our `polymesh` container includes both the `polymesh` binary as well as an
utility binary used by the docker and kubernetes probes to probe the health
of the container.

## The `peer-updater` container

This is a sidecar container that allows dynamic configuration of the
reserved peers of a `polymesh` node.


