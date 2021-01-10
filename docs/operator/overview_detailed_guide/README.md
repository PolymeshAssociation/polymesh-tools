# Polymesh Operator - Overview

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

## Bonding POLYX

To become an operator on Polymesh, you need to bond (lock) a minimum of 5 million POLYX in the
system. This facilitates the economic incentives that the security of Polymesh relies on. The
account that stores your bonded POLYX is called your Stash account and the account that decides
what to do with the bonded POLYX is called your Controller account. Rewards that are generated
for running an operator node can be paid to the Stash account or another specified account.

Note that you do not need to bond all of the POLYX in your Stash account and you can always
bond more POLYX later. However, withdrawing any bonded POLYX requires to wait for the duration
of the unbonding period, which is currently 28 days.

## Network Architecture

The recommended secure operator setup for `mainnet` consists of the following:

* A firewalled (both ingress and egress) active operator node with configured session keys
* A warm spare operator node **without** session keys
* Two or more Internet-facing sentry nodes

The operator node needs to only connect with its sentry nodes.

Sentry nodes are essentially full archive nodes that act as the gatekeeper between your operator
node and the outside world. This setup is intended to isolate the operator node from public
networks, mitigating the risk of a DDoS and other attacks on your operator node.

The operator and sentry nodes do not need to be co-located, but the network between the nodes
should be secured and should allow two-way traffic between the sentries and operators. This may
be achieved via solutions like firewalls, VPN, or a cloud provider’s private networking and peering
solutions. Traffic encryption is preferred but not required.

A *minimum* recommended `testnet` setup would include one operator node and one sentry node.

## High Availability

### Sentry Nodes

The internet-facing sentry nodes should be highly available. An operator node should have at least
two sentry nodes. Two or more operator nodes may share their sentry nodes, but the amount of sentry
nodes should be scaled to provide sufficient redundancy / load balancing / DDoS protection for their assigned
operator nodes.

An operator node needs at least one sentry online at all times so you must make
sure that your sentry nodes are highly available. You can set up as many active sentry nodes for
your operator nodes as you like.

### Operator Node

The network is resilient to temporary outages of some of its operator nodes.  Any one operator
node experience a few minutes of downtime for upgrades, but should not have frequent or extended downtime lest
they risk getting slashed from the network.

It is imperative that only one operator node is active with the same session keys. If multiple
operator nodes with the same session keys do end up online at the same time then they will end up signing
conflicting blocks and will thus get penalised for [equivocation](#terminology).
We recommend that you do not configure automatic failover and instead maintain only a warm
spare that is failed over in a supervised manner.

There are two recommended failover methods:

* Shared session key
* Unique session key

With the shared session key method the operator node session keys are added to the warm spare in
case of a primary operator node failure.  In this case the primary node **must not** come back
online. ***The penalty for equivocation is much higher than the penalty for being offline***.

The uniqye session key method uses different session keys for different instances of operator nodes. If
the primary operator node goes down for some reason, the controller will need to change the
active session keys on the blockchain for the secondary node to become active. Since a key
change takes effect only in the next session, you may still get penalised for being offline for one
session if your primary node went down without producing any blocks in that session. However this
approach eliminates the risk of equivocation penalties.

It is not recommended that you store your controller keys on a server for the automated signing of
the key change transaction. However, you can pre-sign an immortal transaction (a transaction
without a timeout) and store the signed transaction on a server that will broadcast it if the primary
node goes down. Please see [Upgrading or Replacing a Node](#upgrading-or-replacing-a-node) for more details.


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

This is the account where the operator rewards are sent. This should be a cold wallet, never
attached to the operator node.

## Resource requirements

The following resources should be allocated to each Polymesh node:

| Resource | Minimum Value | Recommended Value |
| ---------| --------------| ----------------- |
| CPU      | 2 CPU         | 4 CPU             |
| RAM      | 8 GB          | 8+ GB              |
| Storage  | 80 GB SSD     | 100+ GB local NVMe SSD |

The storage requirements will increase over time as the blockchain grows. Sufficient spare storage
(or expandable volumes) and adequate monitoring measures should be put in place to ensure continued
operations of the node.

## Monitoring

The basic health of a node can be assessed by monitoring the following metrics:

| Parameter | Operational Range | Additional Information |
| --------- | ----------------- | ---------------------- |
| Finalised block number    |+/- 3 from rest of the network|The block number for the rest of the network should be fetched from an external source.|
|Best block time            |6s +/- 1s|The block time is the difference between the best block timestamps.  The ideal mean time is 6 seconds, but some jitter (less than 2s) is acceptable due to network latency|
|Transactions in ready queue|0-150|A healthy node should have zero or near-zero transactions in its ready queue.  A ready queue with a growing number of transactions can be an idicator of excessive node latency|
| Number of peers           |Number of sentries for operators, >2 for sentry nodes|An operator node should have a deterministic number of peers equal to the number of its sentries. A sentry node should have at least two peers (its operator and another network node)|
| Free disk space           |30 GB+ or > 20% volume capacity|There should always be some free disk space for the Polymesh node to consume.|
| Free RAM|1 GB+            |Spikes in RAM usage are acceptable but on average, there should be at least 1 GB of free RAM available on the system for the node to consume.|
| CPU usage|5-50% (overall)|This is the overall CPU usage and not per core usage. Occasional spikes above 50% are acceptable but more cores should be added if the CPU usage continuously stays above 50%.|
| Network connectivity      |<1% packet loss|This mainly applies to sentry nodes. They should be online and reachable at all times. If they are being DDoS’d and can not respond to queries, new sentry nodes should be deployed.|

We have published a Grafana dashboard to monitor the metrics exposed by the Polymesh node via its Prometheus exporter.
You may download it [here](https://github.com/PolymathNetwork/polymesh-tools/tree/main/grafana). In order to use
this dashboard you will need to scrape the metrics from the Prometheus exporter and collect them in a Prometheus
server to which Grafana will connect.

## Upgrading or Replacing a Node

### Sentry Nodes

The upgrade process for sentry nodes varies depending on your network topology.

* If you have only a single sentry node (i.e. when running the minimal `testnet` setup) or low
  redundancy for your sentry nodes it is recommended to create a replacement sentry node first,
  connect the operator node to it, and then terminate the original sentry node (or do a rolling
  upgrade if more than one sentry node requires upgrading). All precautions outlined in
  the [High Availability](#high-availability) section should be observed.
* If you have sufficient redundancy you may just do a rolling upgrade of your sentries. Do ensure
  that your operator nodes reconnect to the upgraded sentries before proceeding to upgrading
  the next sentry.

### Operator Node

The recommended upgrade process for operator nodes is to perform a failover to the warm spare
operator node.

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

On `testnet` you may perform an in-place upgrade if you do not have a warm spare.  We do not
recommend this approach for mainnet due to the risk of penalisation due to downtime in the case
of a failed upgrade.

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

## Stop Being an Operator

To stop being an operator on the Polymesh chain,

* Go to [Staking > Account Actions](https://alcyone-app.polymesh.live/#/staking/actions)
* Click on "Stop Validating" against your bonding account

You will be removed from the operator set in the next session. You can then safely terminate all
your operator and sentry nodes. **failure to terminate safely (e.g. by terminating
before the next session) may result in penalties.**

## Securing the instance

Best practices for securing your instances should be followed at all times. These include:

* disabling password-based SSH
* Setting up and enabling a firewall
* Only opening ports that are needed
* Disabling unnecessary services
* Not using the root user and disabling root login
* Keeping your system up to date
* Turning on SELinux
* Reviewing log files regularly

## Terminology

|Term|Definition|
|----|----------|
|Equivocation|When an operator node commits to two or more conflicting states.|
|Immortal Transaction|A transaction that is valid at any time.|
|Session|A session is a period of time that has a constant set of operators. Operators can only join or exit the operator set at a session change.|
|Era|An Era is a whole number of sessions. It is the period over which operator and nominator sets are calculated, and rewards paid out.|
|Operator Node|Operator nodes are permissioned network participants responsible for producing new blocks and finalising the block chain.|
|Sentry Node|Sentry nodes are full archive nodes which operator nodes use as a proxy to the wider network, limiting the operator nodes exposure to the public internet and providing data redundancy.|

# Polymesh Operator - Detailed Guide

## Overview

This guide will show how a Polymesh operator node can be set up. Please see “Polymesh Operator
- Overview” before continuing to get an overview of the architecture. This guide assumes that you
have already set up a VPN connection between the sentry nodes and the operator nodes.

## Getting the Polymesh node

You can directly build the latest release from source code. You can find the instructions on
https://github.com/PolymathNetwork/Polymesh

## Getting the identity of a Operator node

You need to execute the following command to start your operator and then copy the node's
identity. Once you have copied the identity, close the operator node.

```
./polymesh --operator
```

The above command will start the operator node and it will look something like

```
./polymesh --operator
2020-03-02 11:19:20 Polymesh Node
2020-03-02 11:19:20 version 2.0.0-a8676cab-x86_64-linux-gnu
2020-03-02 11:19:20 by Polymath, 2018-2020
2020-03-02 11:19:20 Chain specification: Local Testnet
2020-03-02 11:19:20 Node name: dirty-vase-9822
2020-03-02 11:19:20 Roles: AUTHORITY
2020-03-02 11:19:20 Local node identity is: QmZ1vCz7QpYsHMug7XLZynqKcueKVWWoTxFqBCRQ487YSr
2020-03-02 11:19:20 Starting BABE Authorship worker
2020-03-02 11:19:20 Grafana data source server started at 127.0.0.1:9955
```

From that, you can see your node's identity (QmZ1vCz7QpYsHMug7XLZynqKcueKVWWoTxFqBCRQ487YSr in this
case). Please save this for later and terminate the operator node.

You can also retrieve the node’s identity after starting it by querying the `system_networkState`
method on the node’s JSONRPC port and reading the result.peerId value, for example:

```
curl \
-s \
-H "Content-Type: application/json" \
-d '{"id":1, "jsonrpc":"2.0", "method": "system_networkState", "params":[]}' \
http://localhost:9933 \
| jq -r .result.peerId
```

## Telemetry Server

Details around telemetry are yet to be decided. A full guide with proper instructions will follow
later.

## Running a Sentry Node

To run a sentry node, you will need to make use of the following options:

* `--sentry`: This parameter enumerates the operators to which the sentry will connect. Each
operator address uses libp2p format, for example
`/ip4/OPERATOR_IP_ADDRESS/tcp/30333/p2p/OPERATOR_PEER_ID`. Multiple operators can be
enumerated with a single `--sentry` parameter, provided that their addresses are separated with a
space.
* `--name`: Human-readable name of the nodes that is reported to the telemetry services.
* `--telemetry-url`: This provides the optional telemetry server to report the node stats and resource
usage. It requires a second parameter to specify the verbosity in the form or a number (0=lowest
verbosity, default=1)

Execute the following command to start a sentry node

```
./polymesh \
--name "Sentry-A" \
--sentry /ip4/OPERATOR_IP_ADDRESS/tcp/30333/p2p/OPERATOR_NODE_IDENTITY \
--telemetry-url ws://TELEMETRY_SERVER_IP:TELEMETRY_SERVER_PORT 0
```

Make sure to write down the identity of the sentry nodes as well.
It is recommended that you run at least two sentry nodes. To start your second sentry node, spin
up a new instance, download/build the Polymesh node on it and run

```
./polymesh \
--name "Sentry-B" \
--sentry /ip4/OPERATOR_IP_ADDRESS/tcp/30333/p2p/OPERATOR_NODE_IDENTITY \
--telemetry-url ws://TELEMETRY_SERVER_IP:TELEMETRY_SERVER_PORT 0
```

## Running an Operator Node

To run an operator node, you need to use the --operator flag in the Polymesh node. You should
also use the --reserved-only flag so that the node only connects to the reserved trusted peers.
The following command will start a Polymesh operator node.

```
./polymesh \
--name "Operator" \
--operator \
--reserved-nodes /ip4/SENTRY_A_IP_ADDRESS/tcp/30333/p2p/SENTRY_A_NODE_IDENTITY \
/ip4/SENTRY_B_IP_ADDRESS/tcp/30333/p2p/SENTRY_B_NODE_IDENTITY \
--reserved-only \
--telemetry-url ws://TELEMETRY_SERVER_IP:TELEMETRY_SERVER_PORT
```

It is recommended that you also set up a failover operator node using the same config on a
different server. You can enumerate both operator addresses in the sentries or update the sentry
configuration at failover time.

Now you need to generate session keys for your operator. Run the following command on the
same machine to generate session keys for your operator.

```
curl -H "Content-Type: application/json" -d '{"id":1, "jsonrpc":"2.0", "method":
"author_rotateKeys", "params":[]}' http://localhost:9933
```

You will get an output similar to:

```
{"jsonrpc":"2.0","result":"0x2bd908203ae740b513f5907fdcc2e29a6bd2835618da917c03d2cfe65
d96745b54d59fe4dc5a106c130be0e677596eb023164c314d6fb5cc62ead1bcaee6a443fe5df859fc1de37
2580abaa98a22fee962bcff580bf57138adc12955aa698a5faa923978d9c16014205af96da9d2e213083ae
fcb53982927a2756ffa83d81658","id":1}
```

Take note of the “result” field. In example above, it is:

```
0x2bd908203ae740b513f5907fdcc2e29a6bd2835618da917c03d2cfe65d96745b54d59fe4dc5a106c130b
e0e677596eb023164c314d6fb5cc62ead1bcaee6a443fe5df859fc1de372580abaa98a22fee962bcff580b
f57138adc12955aa698a5faa923978d9c16014205af96da9d2e213083aefcb53982927a2756ffa83d81658
```

These are the public keys of your session keys. The private keys are stored in a keystore on your
operator server.

NOTE: Before proceeding to the final step that activates your operator node, please wait for all
your nodes to fully sync and make sure that everything has been set up properly.

## Auto restarting nodes

All your nodes should be run using services like systemd so that they are automatically restarted
when a failure happens or the server restarts. You are free to use alternatives but we’ll be using
systemd for this guide.

To get started, create a new systemd unit file called `polymesh-node.service` in
`/etc/systemd/system/`. You may use your favourite text editor, e.g.
`nano /etc/systemd/system/polymesh-node.service`

The following content should be in the unit file

```
[Unit]
Description=Polymesh Node
[Service]
ExecStart=PATH_TO_POLYMESH_BIN POLYMESH_FLAGS_MENTIONED_ABOVE
Restart=always
[Install]
WantedBy=multi-user.target
```

To enable this to autostart on bootup run:

```
systemctl enable polymesh-node.service
```

Start it manually with:

```
systemctl start polymesh-node.service
```

You can check the status of the service with:

```
systemctl status polymesh-node.service
```

## Other Configurations

It is recommended to cap the node’s memory use to 2/3rd of the system RAM or 6144MB,
whichever is greater. This can be achieved in systemd unit files with the MemoryLimit setting.
The `--db-cache` parameter to the polymesh binary can be used to improve the performance on
busy nodes. It is recommended that it be set no lower than the default of 128 (MiB), and that it be
capped at the lesser of 1/2 system RAM or `(system RAM - 4GiB)`.

## Bonding POLYX

To become an operator on Polymesh, you need to bond (lock) some POLYX in the system. The
account that stores your bonded funds is called the stash account and the account that decides
what to do with the bonded funds is called the controller account.

It is highly recommended that you make your controller and stash accounts be two separate
accounts. For this, you will create two accounts and make sure each of them has at least enough
funds to pay the fees for making transactions. Keep most of your funds in the stash account since
it is meant to be the custodian of your staking funds. Please note that for Alcyone Testnet you
can use the same account for the Stash account and the Controller account.

To bond your funds,

* go to the Staking section,
* click on "Account Actions",
* click on the "+”Stash” button.

![Bonding preferences](images/22079145-bec7-4e47-9154-88b0e3dfa964.png "Bonding preferences")

* Stash account
  * Select your Stash account. In this example, we will bond 100 milli POLYX - make
    sure that your Stash account contains at least this much. You can, of course, stake
    more than this.
* Controller account
  * Select the Controller account created earlier. This account will also need a small
    amount of POLYX in order to start and stop validating.
* Value bonded
  * How much POLYX from the Stash account you want to bond/stake. Note that you
    do not need to bond all of the POLYX in that account. Also, note that you can
    always bond more POLYX later. However, withdrawing any bonded amount requires
    to wait for the duration of the unbonding period.
* Payment destination
  * The account where the rewards from validating are sent.
    Once everything is filled in properly, click Bond and sign the transaction with your Stash account.
    After a few seconds, you should see an "ExtrinsicSuccess" message. You should now see a
    new card with all your accounts (note: you may need to refresh the screen). The bonded amount
    on the right corresponds to the funds bonded by the Stash account.

##  Setting Session Keys

You need to tell the Polymesh blockchain what your session keys are. This is what associates your
operator with your Controller account. If you ever want to switch your operator node, you just need
to change your active session keys to the new session keys and wait for the change to become
active in the next session.

Remember the session keys we generated while setting up the operator node? It’s now time to use
those keys.

To set your Session Keys,

* go to the Staking section,
* click on "Account Actions",
* click on the "Session Key" button on the bonding account you generated earlier.
* enter the result of `author_rotateKeys` that we saved earlier in the field and click "Set Session Key",
* submit this extrinsic and you are now ready to start validating.

![Set session key](images/edf14234-3474-43ad-ba3e-910ada7bca52.png "Set session key")

![Set session key](images/aad1824d-d1b9-41e7-9490-2ebf82171c24.png "Set session key")

## Activating your Operator Node

Before moving forward, please make sure that everything is set up properly via the telemetry we
set up earlier. Once this step is complete, an improper setup may lead to penalties.

If everything looks good, go ahead and click on "Validate" in the UI.

![Validate](images/87e444b5-4f94-4408-95c5-b63169fad5b9.png "Validate")

Enter the reward commission percentage and click on Validate.

![Validate](images/eddb483b-2869-4843-8407-bcc329569558.png "Validate")

Congratulations! Your operator has been added in the queue and will become active in the next
session.
