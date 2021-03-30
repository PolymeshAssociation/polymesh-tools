# Polymesh Operator Guide

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
IN THIS MANUAL. POLYMATH ACCEPTS NO LIABILITY AND SHALL NOT BE LIABLE FOR ANY
DAMAGES, DIRECT OR INDIRECT, RESULTING FROM YOUR USE OF THIS MANUAL. IN THE EVENT
POLYMATH, ITS AFFILIATES, LICENSORS OR SUPPLIERS IS FOUND LIABLE, SUCH LIABILITY
SHALL BE LIMITED TO 10 BARBADOS DOLLARS AND THE PAYMENT OF SUCH AMOUNT TO YOU
SHALL BE YOUR EXCLUSIVE REMEDY.

Specifications and information contained in this manual are furnished for
informational use only, and are subject to change without notice, and should
not be construed as advice by Polymath. Recipient must obtain their own
professional or specialist advice before taking, or refraining from, any action
on the basis of the information contained in this manual.

Polymath assumes no responsibility or liability for any errors or inaccuracies
that may appear in this manual and gives no undertaking, and is under no
obligation, to update this document if any errors or inaccuracies become
apparent. The information in this document has not been independently verified.


## Table of Contents

- [Polymesh Operator Guide](#polymesh-operator-guide)
  - [About](#about)
  - [Table of Contents](#table-of-contents)
  - [Introduction](#introduction)
  - [Governance](#governance)
  - [Becoming an Operator](#becoming-an-operator)
  - [Key Management](#key-management)
    - [Session Keys](#session-keys)
    - [Controller Key](#controller-key)
    - [Stash Key](#stash-key)
  - [Network Architecture](#network-architecture)
    - [Firewall traffic](#firewall-traffic)
  - [High Availability](#high-availability)
    - [Operator Node High Availability](#operator-node-high-availability)
  - [Getting the Polymesh Node Software](#getting-the-polymesh-node-software)
  - [Node Resource Requirements](#node-resource-requirements)
  - [Securing the Instances](#securing-the-instances)
  - [Upgrading or Replacing a Node](#upgrading-or-replacing-a-node)
    - [Operator Node Upgrades](#operator-node-upgrades)
  - [Backing Up a Node](#backing-up-a-node)
  - [Auto Restarting Nodes](#auto-restarting-nodes)
  - [Common Parameters for Running a Polymesh Node](#common-parameters-for-running-a-polymesh-node)
  - [Running an Operator Node](#running-an-operator-node)
    - [Generating the session keys with access to the node's RPC port](#generating-the-session-keys-with-access-to-the-nodes-rpc-port)
    - [Generating the session keys in containerised Polymesh nodes](#generating-the-session-keys-in-containerised-polymesh-nodes)
  - [Getting the Identity of a Node](#getting-the-identity-of-a-node)
  - [Metrics and Monitoring](#metrics-and-monitoring)
  - [Bonding POLYX](#bonding-polyx)
  - [Setting Session Keys](#setting-session-keys)
  - [Activating your Operator Node](#activating-your-operator-node)
  - [Stop Being an Operator](#stop-being-an-operator)
  - [Glossary](#glossary)


## Introduction

Operators perform critical functions for the network, and as such, have strict uptime requirements.
This document contains information about the recommended setup and maintenance of a
Polymesh operator node. The intended audience for this document is the operator’s IT team,
however, some business considerations were included for completeness and to provide the
operator’s IT team with the necessary context.


## Governance

Polymesh is a permissioned network meaning potential operators must go through a governance
process in order to be permissioned to work with the Polymesh network. The governance process
is on-chain and managed via the Polymesh Improvement Proposal (PIP) mechanism.

## Becoming an Operator

To become an operator on Polymesh, you may need to bond (lock) POLYX in the
system. This facilitates the economic incentives that the security of Polymesh relies on. The
account that stores your bonded POLYX is called your Stash account and the account that decides
what to do with the bonded POLYX is called your Controller account. Rewards that are generated
for running an operator node can be paid to the Stash account or another specified account.

You do not need to bond all of the POLYX in your Stash account and you can always
bond more POLYX later. However, withdrawing any bonded POLYX requires to wait for the duration
of the unbonding period, which is currently 28 days.

## Key Management

There are three main types of keys that an operator must manage:

* Session keys
* Controller key
* Stash key

The session keys are the only type of keys that the operator node needs access to. The other two
keys should be kept securely in a supported hardware wallet.

### Session Keys

The session keys are the keys that an operator node uses to sign data needed for consensus.
These keys are stored on the operator node itself. Session keys don’t hold any funds but they can
be used to perform actions that will result in a penalty, like double signing. Hence, it is important to
keep these keys secure.

These keys can either be generated offline and injected in the operator node or can be generated
within the operator node by calling the appropriate RPC method. Once generated the session keys
should be persisted.

In the future, Polymesh will support signing payloads outside the client so that keys can be stored
on another device, e.g. a hardware security module (HSM) or a secure enclave. For the time being,
however, session keys must be either stored within the client or be mounted from secure storage
via external methods.

### Controller Key

The controller key is used to manage bonded funds, vote with bonded funds and do similar actions
on chain. This key is not directly needed by the operator node and hence must never be shared
with the operator node. It should be a multisig account or a supported hardware wallet. These keys
can hold funds and directly control funds bonded by the operator and therefore these should be
kept very securely. Consider these keys to be a semi-cold wallet.

### Stash Key

This is the account which holds the POLYX that has been bonded and optionally where the operator rewards
are sent. This should be a cold wallet, never attached to the operator node.

## Network Architecture

The recommended secure operator setup for `itn` consists of the following:

* A firewalled (both ingress and egress) active [operator node](#glossary) with configured session keys
* A [warm spare operator node](#glossary) configured like an operator node but **without** session keys

A *minimum* recommended `testnet` setup would include just a single one operator node.

It is possible to limit the operator node to connecting to just a set of other trusted nodes via the `--reserved-only` and `--reserved-nodes` flags if you need to limit the operators connectivity to the public internet.

### Firewall traffic

To operate properly your Polymesh nodes should have at least the following traffic whitelisted:

* All nodes:
  * **NTP egress**: System clock drift can cause a node to fail to produce blocks due to mismatched
    timestamps between the node and the network.  Ensure your nodes are synchronised with a reliable
    NTP server.
  * **Port 443 egress (HTTPS)** (optional but recommended): Used to send basic telemetry
    to Polymath servers.
* Operator nodes:
  * **Libp2p ingress/egress**: Operator nodes should be able to send and receive p2p events from WAN or a trustedset
    of other nodes that do have WAN connectivity.

## High Availability

### Operator Node High Availability

The network is resilient to temporary outages of some of its operator nodes.  Any one operator
node experience a few minutes of downtime for upgrades, but should not have frequent or extended downtime lest
they risk getting slashed from the network.

It is imperative that only one operator node is active with the same session keys. If multiple
operator nodes with the same session keys do end up online at the same time then they will end up signing
conflicting blocks and will thus get penalised for [equivocation](#glossary).
We recommend that you do not configure automatic failover and instead maintain only a warm
spare that is failed over in a supervised manner.

There are two possible failover methods:

* Shared session key
* Unique session key

With the shared session key method the operator node session keys are added to the warm spare in
case of a primary operator node failure.  In this case the primary node **must not** come back
online. ***The penalty for equivocation is much higher than the penalty for being offline***.

The unique session key method uses different session keys for different instances of operator nodes. If
the primary operator node goes down for some reason, the controller will need to change the
active session keys on the blockchain for the secondary node to become active. Since a key
change takes effect only in the next session, you may still get penalised for being offline for one
session if your primary node went down without producing any blocks in that session. However this
approach eliminates the risk of equivocation penalties.

It is not recommended that you store your controller keys on a server for the automated signing of
the key change transaction. However, you can pre-sign an immortal transaction (a transaction
without a timeout) and store the signed transaction on a server that will broadcast it if the primary
node goes down. Please see [Upgrading or Replacing a Node](#upgrading-or-replacing-a-node) for more details.

## Getting the Polymesh Node Software

Both non-operator and operator nodes use the same binary and only differ in the parameters used to
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

## Node Resource Requirements

The following resources should be allocated to each Polymesh non-operator and operator node:

| Resource | Minimum Value | Recommended Value |
| ---------| --------------| ----------------- |
| CPU      | 2 CPU         | 4 CPU             |
| RAM      | 8 GB          | 8+ GB             |
| Storage  | 80 GB SSD     | 100+ GB low latency SSD (e.g. local NVMe)|

The storage requirements will increase over time as the blockchain grows. Sufficient spare storage
(or expandable volumes) and adequate monitoring measures should be put in place to ensure continued
operations of the node.  A long-running node will keep a large amount of write-ahead logs (WAL) in
the database directory.  These logs are compacted on node restart.  It is recommended that you reserve
an additional 40GB of disk space for the WAL.

It is not recommended that more than one node share the same resources, i.e. it is preferrable to
run two 2 CPU/8 GB RAM instances with one Polymesh node each than running one 4 CPU/16 GB RAM instance
with two Polymesh nodes.

## Securing the Instances

Best practices for securing your instances should be followed at all times. These include (but are not limited to):

* Disabling password-based SSH access
* Setting up and enabling a network firewall
* Only opening ports that are needed
* Disabling unnecessary services
* Not using the root user and disabling root login
* Keeping your system up to date
* Turning on SELinux
* Monitoring logs and metrics for signs of malicious activity
* Running periodic CIS benchmarks against your systems

Be advised: due to the constantly changing landscape of cybersecurity the above list is not and cannot be
comprehensive.  Node operators are responsible that the security of their nodes is up to date
with current best practices.

## Upgrading or Replacing a Node

### Operator Node Upgrades

The recommended upgrade process for operator nodes is to perform a failover to the warm spare
operator node. As mentioned in the [High Availability](#high-availability) section the unique
key approach is preferable to the shared-key approach.

The warm spare operator node should be upgraded first. Since this node is not actively
validating you can simply stop the Polymesh client running on it, perform the necessary upgrade,
and then resume operation.

Once your warm spare operator node is upgraded and fully synchronised, you should make it the
active node by submitting the change on the blockchain using your controller account.

To do so:

1. (If not done already) Generate a new set of session keys for the warm spare operator node
2. Go to [Staking > Account Actions](https://alcyone-app.polymesh.live/#/staking/actions)
3. Click on "Set Session Key" against your bonding account
4. Enter the session keys from the warm node in the field and click on "Set Session Key"

See [Running an Operator Nodes](#running-an-operator-node) for instructions on using the
`author_rotateKeys` RPC method to generate node session keys.

The change in operator session keys only applies in the next session. **For safety, we recommend
that you wait at least 2 sessions before continuing**. In other words, if the current session is `N`, you
should wait until session `N + 2` before proceeding with the steps below.

At this point your warm spare and active operator nodes have switched roles:  The previous warm
spare is now the active operator node and vice-versa.  Be sure to treat them accordingly henceforth.
Alternatively you may perform the failover operation again to restore your original active node
as the current active node and the original warm spare as the current warm spare.

On `testnet` you may perform an in-place upgrade if you do not have a warm spare.  **We do not
recommend in-place upgrades for** `mainnet` **due to the risk of penalisation due to downtime in the case
of a failed upgrade.**

## Backing Up a Node

Since Polymesh is a public blockchain, you do not necessarily need to backup your nodes. You can
always synchronize from scratch.

It takes quite a bit of time to synchronize a node from scratch.  To minimise the time between node
creation and node readiness may choose to back up the full blockchain DB regularly.  This process
does not need to be done on every node - a database backup performed on one node may be used
on another node as long as they have the same setting for the `--pruning` parameter. Since
operator nodes run with an implicit `--pruning archive` setting we recommend that you make that
parameter explicit on all nodes so that they can share a single database backup.

Backing up the database should be done on an offline node. A typical approach to do this would be:

* Stop the polymesh process on the backup node
* Snapshot the database directory
* Restart the polymesh process
* Sync the database snapshot to offsite storage

The database snapshot contains no confidential information as long as **only** the
`/<base path>/chains/alcyone/db` directory is backed up.

Because of the nature of how the database is stored in files, stopping/starting the polymesh
process will create partial database files.  Since an excessive amount of files in a directory
can cause performance issues we recommend to either limit snapshots to a daily frequency or
to periodically reset the backup node's database to a fresh sync from the chain.

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
* `--chain itn`: Run an ITN node. If this parameter is excluded, the default is to connect to the Alcyone network
* `--wasm-execution compiled`: Use compiled wasm to improve performance
* `--base-path <path>` (optional): Specify where Polymesh will look for its DB files and keystore
* `--db-cache <cache size in MiB>` (optional):  Improve the performance of the polymesh process by increasing its
  in-memory cache above the default `128` MiB.  On a node with 8GB RAM available a reasonable value is in the
  ballpark of `4096`.

Note - the `<name>` parameter above will be publicly visible when sending telemetry to Polymath's servers is enabled (on by default).

## Running an Operator Node

To run an operator node you will need to use the following options in additon to
the [common parameters](#common-parameters-for-running-a-polymesh-node):

* `--operator`: Enable the operator flag on the node.

If you wish to connect to just a trusted set of other nodes, you can use the below flags to control this:

* `--reserved-only`: Only connect to reserved peers.
* `--reserved-nodes` (conditionally optional - see notes): This parameter takes a space separated list of libp2p peer addresses
  in the form of `/ip4/<SENTRY_IP_ADDRESS>/tcp/30333/p2p/<SENTRY_NODE_IDENTITY>` or `/dns4/<SENTRY_RESOLVABLE_HOSTNAME>/tcp/30333/p2p/<SENTRY_NODE_IDENTITY>` to which the operator node will connect.  If left out then the peers must be provided via the `system_addReservedPeer` RPC method.  Failure to provide peers via either this parameter or the RPC method will cause the operator node to remain disconnected from the chain.

Next we will generate the node's session keys.

### Generating the session keys with access to the node's RPC port

The `author_rotateKeys` method can be called against a running operator node to generate session keys.

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

**Please wait before activating your operator node until all
your nodes are be fully synced with the chain and make sure that everything is production ready.**

### Generating the session keys in containerised Polymesh nodes

Our official container images contain a small binary to rotate the session keys without requiring the installation
of curl either in the container itself or in a sidecar.  This binary is located in `/usr/local/bin/rotate` and
when executed will produce a newline-terminated string containing the public session keys used for bonding.

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

The above sample log tells us that that node's identity is `12D3KoovCz7QpYsHMug7XLZynqKcueKVWWoTxFqBCRQ487YSrrDG` -
your node's identity will be different. Please save this for later and terminate the operator node process.

To get the node identity via RPC call the `system_localPeerId` method and read the `result` value:

```
$ curl -s -H "Content-Type: application/json" -d '{"id":1, "jsonrpc":"2.0", "method": "system_localPeerId", "params":[]}' http://localhost:9933 | jq -r .result
12D3KoovCz7QpYsHMug7XLZynqKcueKVWWoTxFqBCRQ487YSrrDG
$
```

## Metrics and Monitoring

The recommended approach to getting metrics from the Polymesh node is via its built-in prometheus exporter endpoint.
This endpoint can be scraped with a prometheus-compatible server or agent.

By default the prometheus exporter will
bind to `localhost` on port `9615`.  You can expose the exporter port to additional interfaces with the
`--prometheus-external` flag to enable network based scraping or use a local agent such as `telegraf`,
`grafana-cloud-agent`, or `victoria-metrics-agent` to collect the metrics and push them to the prometheus server.

The basic health of a node can be assessed by monitoring the following metrics:

| Metric | Used for | Operational Range | Additional Notes       |
| ------ | -------- | ----------------- | ---------------------- |
| `polymesh_block_height{status="finalized"}` | Finalised block number | +/- 3 from rest of the network | The block number for the rest of the network should be fetched from an external source. |
| `polymesh_block_height{status="best"}` | Best block time | 6s +/- 2s | The block time is the difference between the best block timestamps.  The ideal mean time is 6 seconds, but some jitter (less than 2s) is acceptable due to network latency |
| `polymesh_ready_transactions_number` | Transactions in ready queue | 0-150 | A healthy node should have zero or near-zero transactions in its ready queue.  A ready queue with a growing number of transactions can be an idicator of excessive node latency |
| `polymesh_sub_libp2p_peers_count` | Number of peers | Number of other nodes for operators | An operator node should maintain connectivity to other operators and either the public internet or a subset of trusted peers |

We have published a Grafana dashboard to monitor the metrics exposed by the Polymesh node via its Prometheus exporter.
You may download it [here](https://github.com/PolymathNetwork/polymesh-tools/tree/main/grafana). In order to use
this dashboard you will need to scrape the metrics from the Prometheus exporter and collect them in a Prometheus
server to which Grafana will connect.

In addition to the Polymesh metrics you should also monitor basic node metrics available from generic node exporters or monitoring agents:

| Metric                    | Operational Range | Additional Notes |
| ------------------------- | ----------------- | ---------------- |
| Free disk space           | 30 GB+ or > 20% volume capacity | There should always be some free disk space for the Polymesh node to consume. |
| Free RAM                  | 1 GB+ | Spikes in RAM usage are acceptable but on average, there should be at least 1 GB of free RAM available on the system for the node to consume.|
| CPU usage                 | 5-50% (overall) | This is the overall CPU usage and not per core usage. Occasional spikes above 50% are acceptable but more cores should be added if the CPU usage continuously stays above 50%. |
| Network connectivity      | <1% packet loss | Nodes should be online and reachable at all times. If they are being DDoS’d and can not respond to queries, new nodes should be deployed, or the operators connectivity limited to trusted nodes. |

## Bonding POLYX

**You should ensure that your Polymesh nodes have synced with the chain and are healthy before proceeding with
this section. Failure to do so may result in operator penalties.**

To become an operator on Polymesh, you need to bond (lock) some POLYX in the system. The
account that stores your bonded funds is called the stash account and the account that decides
what to do with the bonded funds is called the controller account.

*For ITN* `itn` **It is highly recommended that you make your controller and stash accounts be two separate
accounts.** For this, you will create two accounts and make sure each of them has at least enough
funds to pay the fees for making transactions. Keep most of your funds in the stash account since
it is meant to be the custodian of your staking funds.

*For Alcyone* `testnet` *you can use the same account for the Stash account and the Controller account.*

To bond your funds,

* Go to [Staking section](https://itn-app.polymesh.live/#/staking/actions)
* Click on "Account Actions"
* Click on the "+”Stash” button

![Bonding preferences](images/22079145-bec7-4e47-9154-88b0e3dfa964.png "Bonding preferences")

* **Stash account**: Select your Stash account. In this example, we will bond 100 milli POLYX - make
    sure that your Stash account contains at least this much. You can, of course, stake
    more than this.
* **Controller account**: Select the Controller account created earlier. This account will also need a small
    amount of POLYX in order to start and stop validating.
* **Value bonded**: How much POLYX from the Stash account you want to bond/stake. You
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

* Go to [Staking section](https://itn-app.polymesh.live/#/staking/actions)
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

## Stop Being an Operator

To stop being an operator on the Polymesh chain,

* Go to [Staking > Account Actions](https://itn-app.polymesh.live/#/staking/actions)
* Click on "Stop Validating" against your bonding account

You will be removed from the operator set in the next session. You can then safely terminate all
your operator nodes. **failure to terminate safely (e.g. by terminating
before the next session) may result in penalties.**


## Glossary

|Term|Definition|
|----|----------|
|Controller key      |Key used to manage bonded funds, vote with bonded funds and do similar actions on chain.|
|Equivocation        |When an operator node commits to two or more conflicting states.|
|Era                 |An Era is a whole number of sessions. It is the period over which operator and nominator sets are calculated, and rewards paid out.|
|Immortal transaction|A transaction that is valid at any time.|
|Operator node       |Operator nodes are permissioned network participants responsible for producing new blocks and finalising the block chain.|
|Sentry node         |Sentry nodes are full archive nodes which operator nodes use as a proxy to the wider network, limiting the operator nodes exposure to the public internet and providing data redundancy.|
|Session             |A session is a period of time that has a constant set of operators. Operators can only join or exit the operator set at a session change.|
|Session keys        |Keys that an operator node uses to sign data needed for consensus.|
|Stash key           |Account where the operator rewards are sent.|
|Warm spare node     |A node that is online and synced but not configured to be an operator.  A warm spare requires manual intervention to become an active operator.|


