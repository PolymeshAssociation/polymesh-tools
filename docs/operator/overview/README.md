# Polymesh Operator - Overview

## About
Copyright © 2020 Polymath Inc. All Rights Reserved.

No part of this manual, including the products and software described in it, may be reproduced,
transmitted or transcribed to a third-party, or translated into any language in any form or by any
means without the express written permission of Polymath Inc. (“Polymath”).

THIS MANUAL IS PROVIDED “AS-IS” WITHOUT WARRANTY OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING 
BUT NOT LIMITED TO THE IMPLIED WARRANTIES OR CONDITIONS OF COMPLETENESS, ACCURACY, MERCHANTABILITY 
OR FITNESS FOR A PARTICULAR PURPOSE. IN NO EVENT SHALL POLYMATH, ITS AFFILIATES OR ANY OF THEIR 
DIRECTORS, OFFICERS, EMPLOYEES OR AGENTS BE LIABLE FOR ANY INDIRECT, SPECIAL, INCIDENTAL OR CONSEQUENTIAL 
DAMAGES (INCLUDING DAMAGES FOR LOSS OF PROFITS, LOSS OF BUSINESS, LOSS OF USE OR DATA, INTERRUPTION 
OF BUSINESS AND THE LIKE), EVEN IF POLYMATH HAS BEEN ADVISED OF THE POSSIBILITY OF SUCH DAMAGES 
ARISING FROM ANY DEFECT OR ERROR IN THIS MANUAL. POLYMATH ACCEPTS NO LIABILITY AND SHALL NOT BE LIABLE 
FOR ANY DAMAGES, DIRECT OR INDIRECT, RESULTING FROM YOUR USE OF THIS MANUAL. IN THE EVENT POLYMATH, ITS 
AFFILIATES, LICENSORS OR SUPPLIERS IS FOUND LIABLE, SUCH LIABILITY SHALL BE LIMITED TO 10 BARBADOS 
DOLLARS AND THE PAYMENT OF SUCH AMOUNT TO YOU SHALL BE YOUR EXCLUSIVE REMEDY.

Specifications and information contained in this manual are furnished for informational use only,
and are subject to change without notice, and should not be construed as advice by Polymath.
Recipient must obtain their own professional or specialist advice before taking, or refraining from,
any action on the basis of the information contained in this manual. 

Polymath assumes no responsibility or liability for any errors or inaccuracies that may appear in
this manual and gives no undertaking, and is under no obligation, to update this document if any
errors or inaccuracies become apparent. The information in this 
document has not been independently verified.

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

See [Running an Operator Nodes](https://github.com/PolymathNetwork/polymesh-tools/tree/main/docs/operator/detailed_guide/README.md#running-an-operator-node) for instructions on using the
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

## Terminology

|Term|Definition|
|----|----------|
|Equivocation        |When an operator node commits to two or more conflicting states.|
|Immortal transaction|A transaction that is valid at any time.|
|Session             |A session is a period of time that has a constant set of operators. Operators can only join or exit the operator set at a session change.|
|Era                 |An Era is a whole number of sessions. It is the period over which operator and nominator sets are calculated, and rewards paid out.|
|Operator node       |Operator nodes are permissioned network participants responsible for producing new blocks and finalising the block chain.|
|Sentry node         |Sentry nodes are full archive nodes which operator nodes use as a proxy to the wider network, limiting the operator nodes exposure to the public internet and providing data redundancy.|
|Warm spare node     |A node that is online and synced but not configured to be an operator.  A warm spare requires manual intervention to become an active operator.|
