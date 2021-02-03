# Polymesh node helm chart

This helm chart deploys Polymesh nodes (operators, sentries, etc) in Kubernetes.


## Installation

* Add the Polymath chart repository:

```
helm repo add polymath https://charts.polymesh.live/charts/
helm repo up
```

## Examples

### Sentries

* Create the chart configuration file `sentries-values.yaml`

```
---

image:
  pullPolicy: Always

persistence:
  size: 100Gi

replicaCount: 2

polymesh:
  args:
    - --base-path
    -   /var/lib/polymesh
    - --prometheus-external
    - --telemetry-url
    -   "'wss://telemetry.polymesh.live:443/submit 0'"
    - --db-cache
    -   '"4096"'
    - --pruning
    -   archive
```

* Install the sentries

```
helm install --namespace my-namespace polymesh-sentry polymath/polymesh -f sentries-values.yaml
```

This setup will create two sentries (number controlled by the `replicaCount` value) with
an arbitrary `peerId`.  Each replica will also have a separate `LoadBalancer` service so
that you can deterministically pair the `peerId` and external address of the replicas.

You will need both of them to configure the operator nodes that connect to them.


### Operator

* Create the operator keys out of band and load them into a kubernetes secret

```
kubectl create secret generic --namespace my-namespace operator-keys --from-file=path/to-keys/
```

* Create the chart configuration file `operator-values.yaml`

```
---

reservedPeers:
  - "/dns4/sentry-1.my-domain.com/tcp/30333/p2p/12D..."
  - "/ip4/X.X.X.X/tcp/30333/p2p/12D..."

image:
  pullPolicy: Always

persistence:
  size: 100Gi

replicaCount: 1

operatorKeys:
  existingSecret: operator-keys

polymesh:
  args:
    - --base-path
    -   /var/lib/polymesh
    - --operator
    - --prometheus-external
    - --reserved-only
    - --telemetry-url
    -   "'wss://telemetry.polymesh.live:443/submit 0'"
    - --db-cache
    -   '"4096"'
    - --pruning
    -   archive
```

* Install the operator

```
helm install --namespace my-namespace polymesh-operator polymath/polymesh -f operator-values.yaml
```

This will create an operator node and assign it the keys in the kubernetes secret `operator-keys`.  The secret's
keys should have the filename as the key and the file contents as the value. They will be copied into an ephemeral
volume on the pod when it is running. Note that this setup will **not** persist any keys generated during runtime
(e.g. by calling the `rotateKeys` RPC method).

The `reservedPeers` list will be used to create a configmap.  The operator starts with no peers, but a sidecar container
will read the `reservedPeers` list and use the appropriate RPC methods to add the peers after startup.  If you'd rather
not use the sidecar and instead set static reserved peers you can set `peerUpdaterSideCar: false` and adjust the
polymesh args accordingly.

