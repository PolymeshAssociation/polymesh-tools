# Notes

* Our docker repository is specified in the `image` section.  We also have a
  `distroless` tag that has no shell.  The specified health checks will not work
  with that image, but it's more lightweight and safer.
* The `docker-compose` files will intentionally not work out of the box because
  some parameters need to be set (see below)
* The operators and sentries need to know each other's peerId.  The simplest way
  of doing this is by starting the node and getting the information from the
  logs or by querying the `system_localPeerId` method on the jsonrpc or
  websocket rpc port
* The operator and sentry names can be arbitrary, but there is a character
  limit.
* If you wish to scrape the prometheus exporter or query the RPC endpoints from
  non-localhost sources then the appropriate options need to be uncommented
* We recommend running the container as a non-root user.

