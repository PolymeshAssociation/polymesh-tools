# Notes

* Our docker repository is specified in the `image` section.  We also have a
  `distroless` tag that has no shell.  The specified health checks will not work
  with that image, but it's more lightweight and safer.
* The `docker-compose` files will intentionally not work out of the box because
  some parameters need to be set (see below)
* The operator name can be arbitrary, but there is a character limit.
* If you wish to scrape the prometheus exporter or query the RPC endpoints from
  non-localhost sources then the appropriate options need to be uncommented
* We recommend running the container as a non-root user.
* Due to how libp2p peer addresses work you cannot use more than one replica per
  docker-compose deployment. You can use multiple deployments as long as you adjust
  the port mapping.


