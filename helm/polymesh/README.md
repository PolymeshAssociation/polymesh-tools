# Polymesh node helm chart

This helm chart deploys Polymesh nodes (operators, sentries, etc) in Kubernetes.


## Installation

* Add the Polymath chart repository:

```
helm repo add polymath https://charts.polymesh.live/charts/
helm repo up
```

## Examples

### Operator with arbitrary session keys

**NOTE: The method of generating session keys described here will be the default in future releases.**

* Create the chart configuration file `operator-values.yaml`

```
---

image:
  pullPolicy: Always

persistence:
  size: 100Gi

replicaCount: 1

operatorKeys:
  requireSecret: false

polymesh:
  args:
    - --base-path
    -   /var/lib/polymesh
    - --operator
    - --prometheus-external
    - --telemetry-url
    -   "'wss://telemetry.polymesh.live:443/submit 0'"
    - --db-cache
    -   '"4096"'
    - --pruning
    -   archive
    - --wasm-execution
    -   compiled
```

* Install the operator

```
helm install --namespace my-namespace polymesh-operator polymath/polymesh -f operator-values.yaml
```

This will create an operator node without session keys. To generate a new set of session keys run the `/usr/local/bin/rotate`
binary on the operator container.  This will call the `session_rotateKeys` method on the node and print out the string
containing the public portion of the session keys.

`kubectl exec --namespace my-namespace polymesh-operator-0 -- /usr/local/bin/rotate`

### Operator with pregenerated session keys


**NOTE: The method of providing the session keys described here is deprecated and will be removed in future releases.**

* Create the operator keys out of band and load them into a kubernetes secret

```
kubectl create secret generic --namespace my-namespace operator-keys --from-file=path/to-keys/
```

* Create the chart configuration file `operator-values.yaml`

```
---

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
    - --telemetry-url
    -   "'wss://telemetry.polymesh.live:443/submit 0'"
    - --db-cache
    -   '"4096"'
    - --pruning
    -   archive
    - --wasm-execution
    -   compiled
```

* Install the operator

```
helm install --namespace my-namespace polymesh-operator polymath/polymesh -f operator-values.yaml
```

This will create an operator node and assign it the keys in the kubernetes secret `operator-keys`.  The secret's
keys should have the filename as the key and the file contents as the value. They will be copied into an ephemeral
volume on the pod when it is running. Note that this setup will **not** persist any keys generated during runtime
(e.g. by calling the `rotateKeys` RPC method).


