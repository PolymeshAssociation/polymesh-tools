# Polymesh Operator - Detailed Guide

### **Strictly Confidential**

## About
Copyright © 2020 Polymath Inc. All Rights Reserved.

No part of this manual, including the products and software described in it, may be reproduced,
transmitted or transcribed to a third-party, or translated into any language in any form or by any
means without the express written permission of Polymath Inc. (“Polymath”).

THIS MANUAL IS PROVIDED “AS-IS” WITHOUT WARRANTY OF ANY KIND, EITHER EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE IMPLIED WARRANTIES OR CONDITIONS OF
COMPLETENESS, ACCURACY, MERCHANTABILITY OR FITNESS FOR A PARTICULAR PURPOSE. IN
NO EVENT SHALL POLYMATH, ITS AFFILIATES OR ANY OF THEIR DIRECTORS, OFFICERS,
EMPLOYEES OR AGENTS BE LIABLE FOR ANY INDIRECT, SPECIAL, INCIDENTAL OR
CONSEQUENTIAL DAMAGES (INCLUDING DAMAGES FOR LOSS OF PROFITS, LOSS OF BUSINESS,
LOSS OF USE OR DATA, INTERRUPTION OF BUSINESS AND THE LIKE), EVEN IF POLYMATH HAS
BEEN ADVISED OF THE POSSIBILITY OF SUCH DAMAGES ARISING FROM ANY DEFECT OR ERROR
IN THIS MANUAL.

Specifications and information contained in this manual are furnished for informational use only,
and are subject to change without notice, and should not be construed as advice by Polymath.
Recipient must obtain their own professional or specialist advice before taking, or refraining from,
any action on the basis of the information contained in this manual. Polymath assumes no
responsibility or liability for any errors or inaccuracies that may appear in this manual.

Polymath assumes no responsibility or liability for any errors or inaccuracies that may appear in
this manual and gives no undertaking, and is under no obligation, to update this document if any
errors or inaccuracies become apparent. The information in this document has not been
independently verified.

## Overview

This guide will show how a Polymesh operator node can be set up. Please ensure that you have read
the [Polymesh Operator - Overview](https://github.com/PolymathNetwork/polymesh-tools/tree/main/docs/operator/overview/README.md) section first.

## Getting the Polymesh Node Software

Both sentry and operator nodes use the same binary and only differ in the parameters used to
run them.

There are a number of ways to get and deploy the node binary:

* Fetch the prebuilt container image from the [Polymath Docker Hub repository](https://hub.docker.com/r/polymathnet/polymesh).
  There are two flavours available:  `debian` and `distroless`.  The latter has no shell and thus provides a reduced attack
  surface, whereas the former's shell helps with debugging during the initial setup.  The images are tagged with `<flavour>`
  and `<flavour>-<version>`.  We recommend using the latter for deterministic versioning, but the former can be used if you
  set your image pull policy to always pull.  We have also published [sample docker-compose files](https://github.com/PolymathNetwork/polymesh-tools/tree/main/docker-compose).
  The two release flavours (`debian` and `distroless`) are interchangeable in terms of operation - a setup running the
  `debian` flavour can be changed to use the `distroless` flavour by only changing the container tag and vice-versa.
* Fetch the precompiled binary from our [GitHub releases page](https://github.com/PolymathNetwork/Polymesh/releases).  In addition
  to the release source code files we publish four files:  The Polymesh binary and its checksum (identified by the `-linux-amd64`
  suffix indicating the CPU platform it is compiled for) and an archive of Polymesh runtimes and its checksum.  You do not need
  the runtime files as they are already included in the binary.
* Build your own binary from the [release branch of our source code](https://github.com/PolymathNetwork/Polymesh/tree/alcyone)

## Auto Restarting Nodes

All your nodes should automatically restart in the case of an intermittent failure.

For container-based nodes use your container runtime's features: `restart_policy.condition: any`
for `docker-compose`, `restartPolicy: Always` for `kubernetes`, etc.

If running the node as a binary we recommend using a supervisor process to ensure that the
binary is restarted if terminated abnormally.  Most contemporary Linux distributions use
`systemd` for this purpose, so we will focus on that, but you are not limited to using it
if your infrastructure uses a different supervisor process.

To get started, create a new systemd unit file called `polymesh.service` in the
`/etc/systemd/system/` directory. The following content should be in this unit file

```
[Unit]
Description=Polymesh Node

[Service]
ExecStart=<path to polymesh binary> <polymesh parameters>
Restart=always
MemoryLimit=<2/3 the available system RAM, e.g. ~6GB for a system with 8GB RAM>

[Install]
WantedBy=multi-user.target
```

To enable this service to automatically start on bootup run

```
systemctl daemon-reload && systemctl enable polymesh.service
```

You can also `start`, `stop`, `restart`, and check the `status` of the service with the respective `systemctl` commands, e.g.

```
systemctl start polymesh.service
```

The `journalctl` command can be used to read the systemd unit logs:

```
journalctl -u polymesh
```

See the man pages for `journalctl` for more details on how to use that command.

## Common Parameters for Running a Polymesh Node

To run a polymesh node we recommend that you make use of the following options:

* `--name <name>` (optional but recommended): Human-readable name of the nodes that is reported to the telemetry services
* `--pruning archive`: Ensure that the node maintains a full copy of the blockchain
* `--base-path <path>` (optional): Specify where Polymesh will look for its DB files and keystore
* `--db-cache <cache size in MiB>` (optional):  Improve the performance of the polymesh process by increasing its
  in-memory cache above the default `128` MiB.  On a node with 8GB RAM available a reasonable value is in the
  ballpark of `4096`.


## Running a Sentry Node

The sentry node can run with only the [common parameters](#common-parameters-for-running-a-polymesh-node).

You will need to [retrieve the sentry nodes' peer IDs](getting-the-identity-of-a-node) and public IP addresses
and provide them to the operator node(s).

It is recommended that you run at least two sentry nodes on different machines.

## Running an Operator Node

To run an operator node you will need to use the following options in additon to
the [common parameters](#common-parameters-for-running-a-polymesh-node):

* `--operator`: Enable the operator flag on the node.
* `--reserved-only`: Only connect to reserved peers.
* `--reserved-nodes` (conditionally optional - see notes): This parameter takes a space separated list of libp2p peer addresses
  in the form of `/ip4/<SENTRY_IP_ADDRESS>/tcp/30333/p2p/<SENTRY_NODE_IDENTITY>` or `/dns4/<SENTRY_RESOLVABLE_HOSTNAME>/tcp/30333/p2p/<SENTRY_NODE_IDENTITY>` to which
  the operator node will connect.  If left out then the peers must be provided via the `system_addReservedPeer`
  RPC method.  Failure to provide peers via either this parameter or the RPC method will cause the operator node
  to remain disconnected from the chain.

Next call the `author_rotateKeys` method on the operator to generate session keys for your operator node:

```
$ curl -H "Content-Type: application/json" -d '{"id":1, "jsonrpc":"2.0", "method": "author_rotateKeys", "params":[]}' http://localhost:9933 | jq -r .result
```

You will get an output similar to:

```
0x2bd908203ae740b513f5907fdcc2e29a6bd2835618da917c03d2cfe65d96745\
b54d59fe4dc5a106c130be0e677596eb023164c314d6fb5cc62ead1bcaee6a443\
fe5df859fc1de372580abaa98a22fee962bcff580bf57138adc12955aa698a5fa\
a923978d9c16014205af96da9d2e213083aefcb53982927a2756ffa83d81658
```

Take note of this string: it contains the public portion of your session keys. The private
keys are stored in a keystore on your operator server in the `/<base path>/chains/alcyone/keystore/`
directory. The filenames of those keys are the public key portion of the respective session key, and
the contents of the files are the private key portion.

NOTE: Before [activating your operator node](#setting-session-keys) please wait for all
your nodes to be fully synced and make sure that everything is production ready.

## Getting the Identity of a Node

There are two simple methods for getting the public identity of a node:

* From the operator node logs
* Via an RPC call

To get the node identity from the operator node logs start the node process and wait until the line containing the string `Local node identity` is printed:

```
2020-03-02 11:19:20 Polymesh Node
2020-03-02 11:19:20 version 2.0.0-a8676cab-x86_64-linux-gnu
2020-03-02 11:19:20 by Polymath, 2018-2020
2020-03-02 11:19:20 Chain specification: Local Testnet
2020-03-02 11:19:20 Node name: dirty-vase-9822
2020-03-02 11:19:20 Roles: AUTHORITY
2020-03-02 11:19:20 Local node identity is: 12D3KoovCz7QpYsHMug7XLZynqKcueKVWWoTxFqBCRQ487YSrrDG
2020-03-02 11:19:20 Starting BABE Authorship worker
2020-03-02 11:19:20 Grafana data source server started at 127.0.0.1:9955
...
```

After that the process can be gracefully terminated.

The above sample log tells us that that node's identity is `12D3KoovCz7QpYsHMug7XLZynqKcueKVWWoTxFqBCRQ487YSrrDG` -
your node's identity will be different. Please save this for later and terminate the operator node process.

To get the node identity via RPC call the `system_localPeerId` method and read the `result` value:

```
$ curl -s -H "Content-Type: application/json" -d '{"id":1, "jsonrpc":"2.0", "method": "system_localPeerId", "params":[]}' http://localhost:9933 | jq -r .result
12D3KoovCz7QpYsHMug7XLZynqKcueKVWWoTxFqBCRQ487YSrrDG
$
```

## Metrics

The recommended approach to getting metrics from the Polymesh node is via its built-in prometheus exporter endpoint.
This endpoint can be scraped with a prometheus-compatible server or agent.

By default the prometheus exporter will
bind to `localhost` on port `9615`.  You can expose the exporter port to additional interfaces with the
`--prometheus-external` flag to enable network based scraping or use a local agent such as `telegraf`,
`grafana-cloud-agent`, or `victoria-metrics-agent` to collect the metrics and push them to the prometheus server.

Once collected the metrics can be monitored and charted with various prometheus-compatible tools.  We provide a sample
[Grafana dashboard](https://github.com/PolymathNetwork/polymesh-tools/tree/main/grafana) which monitors the most
common metrics and includes some basic alerts for said metrics.

## Bonding POLYX

To become an operator on Polymesh, you need to bond (lock) some POLYX in the system. The
account that stores your bonded funds is called the stash account and the account that decides
what to do with the bonded funds is called the controller account.

It is highly recommended that you make your controller and stash accounts be two separate
accounts. For this, you will create two accounts and make sure each of them has at least enough
funds to pay the fees for making transactions. Keep most of your funds in the stash account since
it is meant to be the custodian of your staking funds.

**Please note that for Alcyone *Testnet* you
can use the same account for the Stash account and the Controller account.**

To bond your funds,

* Go to [Staking section](https://alcyone-app.polymesh.live/#/staking/actions)
* Click on "Account Actions"
* Click on the "+”Stash” button

![Bonding preferences](images/22079145-bec7-4e47-9154-88b0e3dfa964.png "Bonding preferences")

* **Stash account**: Select your Stash account. In this example, we will bond 100 milli POLYX - make
    sure that your Stash account contains at least this much. You can, of course, stake
    more than this.
* **Controller account**: Select the Controller account created earlier. This account will also need a small
    amount of POLYX in order to start and stop validating.
* **Value bonded**: How much POLYX from the Stash account you want to bond/stake. Note that you
    do not need to bond all of the POLYX in that account. Also, note that you can
    always bond more POLYX later. However, withdrawing any bonded amount requires
    to wait for the duration of the unbonding period.
* **Payment destination**: The account where the rewards from validating are sent.
    Once everything is filled in properly, click Bond and sign the transaction with your Stash account.
    After a few seconds, you should see an `ExtrinsicSuccess` message. You should now see a
    new card with all your accounts (note: you may need to refresh the screen). The bonded amount
    on the right corresponds to the funds bonded by the Stash account.

## Setting Session Keys

You need to tell the Polymesh blockchain what your session keys are. This is what associates your
operator with your Controller account. If you ever want to switch your operator node, you just need
to change your active session keys to the new session keys and wait for the change to become
active in the next session.

Remember the session keys we generated while setting up the operator node? It’s now time to use
those keys.

To set your Session Keys,

* Go to [Staking section](https://alcyone-app.polymesh.live/#/staking/actions)
* Click on "Account Actions"
* Click on the "Session Key" button on the bonding account you generated earlier
* Enter the result of `author_rotateKeys` that we saved earlier in the field and click "Set Session Key"
* Submit this extrinsic and you are now ready to start validating

![Set session key](images/edf14234-3474-43ad-ba3e-910ada7bca52.png "Set session key")

![Set session key](images/aad1824d-d1b9-41e7-9490-2ebf82171c24.png "Set session key")

## Activating your Operator Node

Before moving forward, please make sure that everything is set up properly via the telemetry we
set up earlier. Once this step is complete, an improper setup may lead to penalties.

If everything looks good, go ahead and click on "Validate" in the UI.

![Validate](images/87e444b5-4f94-4408-95c5-b63169fad5b9.png "Validate")

Enter the reward commission percentage and click on Validate.

![Validate](images/eddb483b-2869-4843-8407-bcc329569558.png "Validate")

**Congratulations!** Your operator has been added in the queue and will become active in the next
session.

