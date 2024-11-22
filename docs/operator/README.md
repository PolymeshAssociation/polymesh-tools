# Polymesh Operator Guide

## Table of Contents

- [Polymesh Operator Guide](#polymesh-operator-guide)
  - [Table of Contents](#table-of-contents)
  - [Introduction](#introduction)
  - [Governance](#governance)
  - [Becoming an Operator](#becoming-an-operator)
  - [Key Management](#key-management)
    - [Session Keys](#session-keys)
    - [Permissioned Identity](#permissioned-identity)
    - [Primary key](#primary-key)
    - [Secondary keys](#secondary-keys)
    - [Stash Key](#stash-key)
    - [Controller Key](#controller-key)
  - [Network Architecture](#network-architecture)
    - [Firewall Traffic](#firewall-traffic)
  - [High Availability](#high-availability)
    - [Operator Node High Availability](#operator-node-high-availability)
  - [Getting the Polymesh Node Software](#getting-the-polymesh-node-software)
  - [Node Resource Requirements](#node-resource-requirements)
  - [Securing the Instances](#securing-the-instances)
  - [Upgrading or Replacing a Node](#upgrading-or-replacing-a-node)
    - [Operator Node Upgrades](#operator-node-upgrades)
  - [Backing Up a Node](#backing-up-a-node)
  - [Auto-Restarting Nodes](#auto-restarting-nodes)
    - [Container-Based Nodes](#container-based-nodes)
    - [Binary-Based Nodes](#binary-based-nodes)
      - [Setting Up `systemd`](#setting-up-systemd)
  - [Common Parameters for Running a Polymesh Node](#common-parameters-for-running-a-polymesh-node)
  - [Running an Operator Node](#running-an-operator-node)
    - [Generating Session Keys with Access to the Node's RPC Port](#generating-session-keys-with-access-to-the-nodes-rpc-port)
    - [Generating Session Keys in Containerized Polymesh Nodes](#generating-session-keys-in-containerized-polymesh-nodes)
    - [Getting the Identity of a Node](#getting-the-identity-of-a-node)
      - [From the Operator Node Logs](#from-the-operator-node-logs)
      - [Via RPC Call](#via-rpc-call)
    - [Metrics and Monitoring](#metrics-and-monitoring)
      - [Polymesh-Specific Metrics](#polymesh-specific-metrics)
      - [General Node Metrics](#general-node-metrics)
    - [Bonding POLYX](#bonding-polyx)
      - [Steps to Bond Funds](#steps-to-bond-funds)
    - [Setting Session Keys](#setting-session-keys)
      - [Steps to Set Session Keys](#steps-to-set-session-keys)
    - [Activating Your Operator Node](#activating-your-operator-node)
      - [Steps to Activate](#steps-to-activate)
    - [Stopping Operator Activity](#stopping-operator-activity)
      - [Reward Destination Considerations](#reward-destination-considerations)
  - [Glossary](#glossary)

## Introduction

Operators perform critical functions for the network and have strict uptime requirements. This document provides information about the recommended setup and maintenance of a Polymesh operator node. The intended audience is the operator’s IT team; however, some business considerations have been included for completeness and to provide the necessary context.

## Governance

Polymesh is a permissioned network, meaning potential operators must undergo a governance process to be permitted to work with the Polymesh network. This governance process is on-chain and managed via the Polymesh Improvement Proposal (PIP) mechanism.

## Becoming an Operator

The node operator role on Polymesh requires a permission to be assigned. All Operators satisfy selection criteria determined by the Polymesh Association and be approved by the the Polymesh Governing Council.

To become an operator on Polymesh, you also need to bond (lock) POLYX in the system. This facilitates the economic incentives on which Polymesh's security relies. The account that stores your bonded POLYX is called the Stash account, while the account that manages the bonded POLYX is called the Controller account. Rewards generated for running an operator node can be sent to the Stash account or another specified account.

You do not need to bond all the POLYX in your Stash account and can bond more later. However, withdrawing any bonded POLYX requires waiting for the unbonding period, currently set at 28 days.

## Key Management

**NB**: It is recommended that operators use Ledger Nano S Plus, Nano X, Flex or Stax devices to store their keys. The discontinued Ledger Nano S should **NOT** be used.

The Nano S does not support setting session keys and is therefore unsuitable for operators.

There are three main types of keys that an operator must manage:

- Session keys
- Controller key
- Stash key

The session keys are the only type of key that the operator node needs access to. The other two keys should be securely stored in a supported hardware wallet.

### Session Keys

Session keys are used by the operator node to sign data necessary for consensus. These keys are stored on the operator node itself. Although session keys do not hold any funds, they can be used to perform actions that could result in penalties, such as double signing. It is therefore critical to keep these keys secure.

Session keys can either be generated offline and injected into the operator node or generated within the operator node by calling the appropriate RPC method. Once generated, session keys should be persisted.

Session keys must either be stored within the client or mounted from secure storage via external methods.

### Permissioned Identity

On Polymesh all keys (excluding Session keys) must be linked to a on chain DID with a valid CDD claim to be able to receive POLYX and submit transactions. A node operators DID must be granted an additional role making it a permissioned identity. Only keys linked to this identity will be allowed to be Stash keys.

### Primary key

This key is a special key under the onchain identity. It has additional capabilities of adding and removing secondary keys from an identity. It cannot access POLYX on secondary keys. It is recommended to keep the Primary key as in a cold wallet and **NOT** use it as a stash key.

### Secondary keys

Secondary keys are authorized by a primary key to join an identity. For the purposes of staking/POLYX transfers they are no different to primary keys. Secondary keys can be given permissions to perform specific onchain actions. e.g. they may be allowed to interact only with specific assets, only with specific portfolios or only use specific transactions.

### Stash Key

The stash key is the account that holds the POLYX bonded by the operator and, optionally, receives operator rewards. This should be a cold wallet and must never be attached to the operator node.

Stash keys can be either a primary or a secondary key but **MUST** be linked to the permissioned DID of the operator.

Stash keys can:

- Bond POLYX
- Bond extra POLYX
- Set the Controller key of the stash. If a controller key is not provided the Stash key is automatically made the Controller.

### Controller Key

The controller key is used to manage bonded funds, vote with bonded funds, and perform similar on-chain actions. This key is not directly required by the operator node and should never be shared with it. It is recommended to use a multisig account or a supported hardware wallet for the controller key. These keys can hold funds and directly control bonded funds, so they must be stored securely. Consider the controller key a semi-cold wallet.

Controller keys can be a Primary or Secondary key and can be associated with any identity. They do not need to be associated with the permissioned identity of the Stash.

Controller keys can:

- Set/update session keys
- Set commission and validate
- Set a destination for reward payments - i.e. stash with automatically stake, stash unstaked, controller or other account.
- Commence unbonding of stash tokens
- Withdraw unbonded tokens after the 28 day waiting period to the stash account. i.e. the controller cannot transfer unbonded tokens from the stash.
- Rebond tokens which are in the process of unbonding
- "Chill" the node i.e. gracefully stop validating after the last era the node was elected to ends.
- Purge Session keys

## Network Architecture

The recommended secure operator setup for Testnet/Mainnet includes the following:

- An active [operator node](#glossary) configured with session keys
- A [warm spare operator node](#glossary) configured like an operator node but **without** session keys

A _minimum_ setup consists of a single operator node.

The `--reserved-only` flag, combined with the `--reserved-nodes` parameter, allows you to restrict connections to a whitelist of nodes that the operator node may peer with.

### Firewall Traffic

To function properly, Polymesh nodes require the following traffic to be whitelisted:

- **All Nodes**:

  - **NTP Egress**: System clock drift can cause a node to fail to produce blocks due to mismatched timestamps. Ensure your nodes are synchronized with a reliable NTP server.
  - **Port 443 Egress (HTTPS)** (optional but recommended): Used to send basic telemetry to Polymesh servers. Enabling telemetry allows your node to appear on the [Polymesh Telemetry page](https://stats.polymesh.network/).

- **Operator Nodes**:
  - **Libp2p Ingress/Egress**: Operator nodes **must** send and receive P2P events from the WAN or a trusted set of nodes with WAN connectivity. _(Default Port: 30333)_

## High Availability

### Operator Node High Availability

The network tolerates simultaneous outages of multiple operator nodes, provided a critical threshold is maintained. An operator node can experience brief downtime for maintenance or upgrades. However, frequent or prolonged downtime—or outages of multiple nodes simultaneously—risks incurring slashing penalties.

**Only one operator node may be active with the same session keys at a time.** If multiple operator nodes with identical session keys are online simultaneously, they may sign conflicting blocks, resulting in [equivocation](#glossary) penalties. Therefore, automatic failover is discouraged. Instead, maintain a warm spare node for supervised failover.

Two failover methods are available:

- **Shared Session Key**: The session keys are shared with the warm spare, which is activated if the primary node fails. The primary node **must not** come back online. **_The penalty for equivocation is much higher than for being offline._**
- **Unique Session Key**: Each operator node instance has a unique session key. If the primary node fails, the controller must update the session keys on-chain for the secondary node. As key changes take effect in the next session, there may still be a period where the node is offline for one session. This method eliminates the risk of equivocation penalties.

Storing controller keys on a server for automated key change transactions is **not** recommended. However, you may pre-sign an immortal transaction (without a timeout) and store it on a server to broadcast if the primary node fails. See [Upgrading or Replacing a Node](#upgrading-or-replacing-a-node) for details.

## Getting the Polymesh Node Software

All Polymesh nodes use the same binary, differing only in runtime parameters.

There are several ways to obtain the node binary:

- **Prebuilt Container Images**: Fetch from the [Polymesh Docker Hub repository](https://hub.docker.com/r/polymeshassociation/polymesh). Two flavors are available: `debian` (with a shell for easier debugging) and `distroless` (reduced attack surface, no shell). Images are tagged as `<flavor>` or `<flavor>-<version>`. Use versioned tags for deterministic updates. See our [sample Docker Compose files](https://github.com/PolymeshAssociation/polymesh-tools/tree/main/docker-compose). Refer to the Polymesh Developer Documentation for a guide to [running a Polymesh node with Docker](https://developers.polymesh.network/polymesh-docs/network/running-a-node-docker/).

- **Precompiled Binary**: Download from the [GitHub releases page](https://github.com/PolymeshAssociation/Polymesh/releases). Releases include the Polymesh binary, its checksum, and optional runtime archives. The runtimes are already included in the binary.

- **Build from Source**: Clone the [release branch](https://github.com/PolymeshAssociation/Polymesh/tree/mainnet) and follow the repository instructions to compile the binary.

## Node Resource Requirements

At the time of writing, each Polymesh node should have the following minimum resources:

| Resource | Minimum Value | Recommended Value                          |
| -------- | ------------- | ------------------------------------------ |
| CPU      | 2 CPUs        | 4 CPUs                                     |
| RAM      | 8 GB          | 8+ GB                                      |
| Storage  | 250 GB SSD    | 400+ GB low-latency SSD (e.g., local NVMe) |

As the blockchain grows, storage requirements will increase. Ensure sufficient spare storage or expandable volumes and monitor disk usage. A long-running node accumulates write-ahead logs (WAL) in the database directory. These logs are compacted upon node restart. Reserve an additional 40 GB of disk space for WAL.

Running multiple nodes on shared resources is not recommended. For example, it is preferable to run two 2-CPU/8-GB RAM instances (one node each) than a single 4-CPU/16-GB RAM instance hosting two nodes.

## Securing the Instances

Best practices for securing your instances should always be followed. These include (but are not limited to):

- Disabling password-based SSH access
- Setting up and enabling a network firewall
- Only opening required ports
- Disabling unnecessary services
- Avoiding the use of the root user and disabling root login
- Keeping your system up to date
- Enabling SELinux
- Monitoring logs and metrics for signs of malicious activity
- Running periodic CIS benchmarks against your systems

**Be advised:** Due to the constantly evolving cybersecurity landscape, the above list is not, and cannot be, comprehensive. Node operators are responsible for ensuring that their nodes remain secure and adhere to current best practices.

## Upgrading or Replacing a Node

### Operator Node Upgrades

The recommended process for upgrading operator nodes is to perform a failover to the warm spare operator node. As mentioned in the [High Availability](#high-availability) section, the unique key approach is preferable to the shared-key approach.

Begin by upgrading the warm spare operator node. Since this node is not actively validating, you can stop the Polymesh client, perform the necessary upgrade, and then resume operation.

Once the warm spare node is upgraded and fully synchronized, make it the active node by submitting a change of session keys associated with your stash to those stored in the warm spare's keystore, using your controller account.

To do so:

1. (If not already done) Generate a new set of session keys for the warm spare operator node.
2. Navigate to [Staking > Account Actions](https://mainnet-app.polymesh.network/#/staking/actions).
3. Click "Set Session Key" for your bonding account.
4. Enter the session keys from the warm spare node in the field and click "Set Session Key."

See [Running an Operator Node](#running-an-operator-node) for instructions on using the `author_rotateKeys` RPC method to generate session keys.

The change in operator session keys will only take effect in the next session. **For safety, we recommend waiting at least two sessions before proceeding.** If the current session is `N`, wait until session `N + 2` before continuing.

At this point, the warm spare and active operator nodes will have switched roles: the previous warm spare becomes the active node, and vice versa. Ensure they are treated accordingly going forward. Alternatively, you may perform the failover operation again to restore the original active node as the current active node and the original warm spare as the current warm spare.

If a warm spare is unavailable, you may perform an in-place upgrade by stopping the Polymesh client, performing the necessary client upgrade, and restarting the node client. **We do not recommend in-place upgrades due to the risk of encountering a failure during the upgrade.**

## Backing Up a Node

Since Polymesh is a public blockchain, node backups are not strictly required, as you can always synchronize from scratch. However, syncing from scratch can be time-consuming. To minimize the delay between node creation and readiness, you may choose to back up the full blockchain database regularly.

This does not need to be done for every node—a single database backup can be used across nodes, provided they use the same `--pruning` setting. Operator nodes run with an implicit `--pruning archive` setting, so we recommend explicitly setting this parameter on all nodes to allow sharing a single database backup.

Database backups should be performed on an offline node. A typical approach is:

1. Stop the Polymesh process on the backup node.
2. Snapshot the database directory.
3. Restart the Polymesh process.
4. Sync the database snapshot to offsite storage.

The database snapshot contains no confidential information as long as **only** the `db` directory is backed up (e.g., `/<base path>/chains/mainnet/db`).

Because the database uses file-based storage, stopping/starting Polymesh may create partial files. Excessive file accumulation can cause performance issues. We recommend limiting snapshots to daily intervals and periodically resetting the backup node's database with a fresh sync from the chain.

## Auto-Restarting Nodes

Nodes should automatically restart in the event of an intermittent failure.

### Container-Based Nodes

For container-based nodes, use your container runtime's features:

- `restart_policy.condition: any` for `docker-compose`
- `restartPolicy: Always` for `kubernetes`

### Binary-Based Nodes

For binary-based nodes, we recommend using a supervisor process. Most modern Linux distributions use `systemd`, which we will focus on, though other options are also viable.

#### Setting Up `systemd`

1. Create a new unit file called `polymesh.service` in `/etc/systemd/system/` with the following content:

   ```ini
   [Unit]
   Description=Polymesh Node

   [Service]
   ExecStart=<path to polymesh binary> <polymesh parameters>
   Restart=always
   MemoryLimit=<2/3 of available system RAM, e.g., ~6GB for an 8GB system>

   [Install]
   WantedBy=multi-user.target
   ```

2. Enable automatic startup with:

   ```bash
   systemctl daemon-reload && systemctl enable polymesh.service
   ```

3. Manage the service with commands such as:

   ```bash
   systemctl start polymesh.service
   ```

4. View logs using `journalctl`:

   ```bash
   journalctl -u polymesh
   ```

Refer to the `journalctl` man pages for additional details.

## Common Parameters for Running a Polymesh Node

Recommended options for running a Polymesh node include:

- `--name <name>` (optional): Human-readable name reported to telemetry services.
- `--pruning archive`: Maintain a full blockchain copy.
- `--chain mainnet`: Run a Mainnet node (default is Testnet if omitted).
- `--wasm-execution compiled`: Use compiled WASM for better performance.
- `--base-path <path>` (optional): Specify the location for DB files and the keystore.
- `--db-cache <cache size in MiB>` (optional): Increase in-memory cache for better performance. On a node with 8GB of available RAM, a reasonable value is `4096`. (default `128` MiB)
- `--db-max-total-wal-size <WAL database size in MiB>` (optional): Limit the total storage capacity that the database can use for WAL files. Recommended minimum value `1024`.

**Note:** The `<name>` parameter will be publicly visible when telemetry is enabled (default setting).

To see a full list of available options and their descriptions you can viewed by including the `--help` command.

## Running an Operator Node

To run an operator node, use the following in addition to the [common parameters](#common-parameters-for-running-a-polymesh-node):

- `--operator`: Enable operator mode.

To connect only to trusted peers, use these options:

- `--reserved-only`: Restrict connections to reserved peers.
- `--reserved-nodes`: A space-separated list of libp2p peer addresses in the format `/ip4/<IP>/tcp/30333/p2p/<Node ID>` or `/dns4/<Hostname>/tcp/30333/p2p/<Node ID>`. If omitted, peers must be added via the `system_addReservedPeer` RPC method.

Next, generate the node's session keys.

### Generating Session Keys with Access to the Node's RPC Port

The `author_rotateKeys` method can be called against a running operator node to generate session keys.

```bash
curl -H "Content-Type: application/json" -d '{"id":1, "jsonrpc":"2.0", "method": "author_rotateKeys", "params":[]}' http://localhost:9933 | jq -r .result
```

You will get an output similar to:

```plaintext
0x2bd908203ae740b513f5907fdcc2e29a6bd2835618da917c03d2cfe65d96745\
b54d59fe4dc5a106c130be0e677596eb023164c314d6fb5cc62ead1bcaee6a443\
fe5df859fc1de372580abaa98a22fee962bcff580bf57138adc12955aa698a5fa\
a923978d9c16014205af96da9d2e213083aefcb53982927a2756ffa83d81658
```

Take note of this string: it contains the **public** portion of your session keys. The private keys are stored in a keystore on your operator server in the `/<base path>/chains/<chain name>/keystore/` directory. The filenames of these keys are the public key portion of the respective session key, and the file contents represent the private key.

**Wait until your operator node is fully synced with the chain and production-ready before activation.**

### Generating Session Keys in Containerized Polymesh Nodes

Our official container images include a utility to rotate session keys without requiring additional tools like `curl` or exposing unsafe RPC methods outside the container. This utility is located at `/usr/local/bin/rotate`. Running it will produce a newline-terminated string containing the public session keys used for bonding.

To run the utility inside your container:

```bash
docker exec <container_name> /usr/local/bin/rotate
```

See our [guide to running a node with Docker](https://developers.polymesh.network/polymesh-docs/network/running-a-node-docker/#generating-node-session-keys) for more details.

### Getting the Identity of a Node

When your node is first initialized a Local node identity is generated if none is present in your nodes `network` folder. The key used to generate it is stored so your nodes peer ID persists across restarts.

There are two ways to obtain the public identity of a node:

- From the operator node logs
- Via an RPC call

#### From the Operator Node Logs

Start the node process and look for a line containing `Local node identity`:

```plaintext
2024-11-15 18:50:11 Reserved nodes: []
2024-11-15 18:50:11 Polymesh Node
2024-11-15 18:50:11 ✌️  version 7.0.0
2024-11-15 18:50:11 ❤️  by PolymeshAssociation, 2017-2024
2024-11-15 18:50:11 📋 Chain specification: Polymesh Testnet
2024-11-15 18:50:11 🏷  Node name: woebegone-galley-5149
2024-11-15 18:50:11 👤 Role: FULL
2024-11-15 18:50:11 💾 Database: RocksDb at /var/lib/polymesh/chains/testnet/db/full
2024-11-15 18:50:11 ⛓  Native runtime: polymesh_testnet-7000005 (polymesh_testnet-0.tx7.au1)
2024-11-15 18:50:12 🔨 Initializing Genesis block/state (state: 0xcba3…bce0, header-hash: 0x2ace…d0d6)
2024-11-15 18:50:12 👴 Loading GRANDPA authority set from genesis on what appears to be first startup.
2024-11-15 18:50:13 👶 Creating empty BABE epoch changes on what appears to be first startup.
2024-11-15 18:50:13 🏷  Local node identity is: 12D3KooWSDAHjBmA6j2GyBZPktEz2gLZmtJAc2bWnDV7eCCcgcbC
2024-11-15 18:50:13 💻 Operating system: linux
...
```

In this example, the node's identity is `12D3KooWSDAHjBmA6j2GyBZPktEz2gLZmtJAc2bWnDV7eCCcgcbC`. Save this value, then terminate the process.

#### Via RPC Call

Call the `system_localPeerId` method and read the `result` value:

```bash
curl -s -H "Content-Type: application/json" -d '{"id":1, "jsonrpc":"2.0", "method": "system_localPeerId", "params":[]}' http://localhost:9933 | jq -r .result
```

### Metrics and Monitoring

The recommended method for obtaining metrics is through the node's built-in Prometheus exporter. By default, it binds to `localhost` on port `9615`. Use the `--prometheus-external` flag to expose the exporter port for network-based scraping, or deploy a local agent (e.g., `telegraf`, `grafana-cloud-agent`, or `victoria-metrics-agent`) to collect metrics.

#### Polymesh-Specific Metrics

The basic health of a node can be assessed by monitoring the following metrics:

| Metric                                      | Purpose                     | Range                   | Notes                                                                                                                     |
| ------------------------------------------- | --------------------------- | ----------------------- | ------------------------------------------------------------------------------------------------------------------------- |
| `polymesh_block_height{status="finalized"}` | Finalized block number      | ±3 from network average | Verify against the network block number using an external source.                                                         |
| `polymesh_block_height{status="best"}`      | Best block time             | 6s ±2s                  | Ideal mean block time is 6 seconds; minor variations are normal.                                                          |
| `polymesh_ready_transactions_number`        | Transactions in ready queue | 0–150                   | A growing queue may indicate node latency issues.                                                                         |
| `polymesh_sub_libp2p_peers_count`           | Number of peers             | > Minimum peer count    | Nodes should maintain connectivity with other operator nodes, ideally with a maximum of three hops to any other operator. |

Guides for monitoring other Substrate-based chains, such as [Polkadot](https://wiki.polkadot.network/docs/maintain-guides-how-to-monitor-your-node), can be referenced for additional approaches to node monitoring.

#### General Node Metrics

In addition to Polymesh-specific metrics, you should monitor basic node health metrics available from generic node exporters or monitoring agents:

| Metric               | Range                   | Notes                                                                                                                   |
| -------------------- | ----------------------- | ----------------------------------------------------------------------------------------------------------------------- |
| Free disk space      | 30 GB+ or >20% capacity | Ensure sufficient free space for the node to function properly.                                                         |
| Free RAM             | 1 GB+                   | Spikes are acceptable, but an average of 1 GB of free RAM should remain available.                                      |
| CPU usage            | 5–50% (overall)         | If usage consistently exceeds 50%, consider increasing CPU resources.                                                   |
| Network connectivity | >0.1 Mbps bandwidth     | Maintain stable connectivity. If under attack (e.g., DDoS), deploy new nodes or restrict connectivity to trusted peers. |

### Bonding POLYX

**Ensure your Polymesh nodes are fully synced and healthy before proceeding. Failure to do so may result in operator penalties.**

To become an operator on Polymesh, you must bond (lock) some POLYX. The account holding the bonded funds is called the **stash account**, while the account managing these funds is the **controller account**.

> **Recommendation:** Use separate accounts for your controller and stash accounts. Ensure both accounts have enough POLYX to cover transaction fees, while keeping the majority of your funds in the stash account, as it serves as the custodian of your staking funds.

**Note:** The steps in the following sections can also be completed in a single batched transaction by selecting the **[+ Validator]** button. The instructions below focus on each individual step.

#### Steps to Bond Funds

1. Navigate to the [Staking section](https://mainnet-app.polymesh.network/#/staking/actions).
2. Click **Account Actions**.
3. Select the **[+ Stash]** button.

   ![Bonding preferences](images/bonding-preferences.png 'Bonding preferences')

4. **Complete the following fields:**

   - **Stash account**: Choose your [stash account](#stash-key). Ensure it has a balance exceeding the amount of POLYX you plan to bond. You can bond additional POLYX later if needed.
   - **Controller account**: Select your [controller account](#controller-key). This account requires a small amount of POLYX for managing validation.
   - **Value bonded**: Specify the amount of POLYX to bond/stake from the stash account. Note that withdrawing bonded POLYX requires waiting for the unbonding period to elapse. At the time of writing, a minimum of 50,000 POLYX is required for operator staking.
   - **Payment destination**: Specify where validation rewards should be sent.

5. Click **Bond** and sign the transaction using your stash account. After a few seconds, an `ExtrinsicSuccess` message should appear. Refresh the page if necessary to view your new bonded account details.

   ![Stash bonded](images/stash-bonded.png 'Stash bonded')

### Setting Session Keys

Session keys link your operator node to your controller account. If you switch operator nodes, you can update these keys to reflect the new setup, which will take effect in the next session.

Use the session keys generated during your operator node setup.

#### Steps to Set Session Keys

1. Go to the [Staking section](https://mainnet-app.polymesh.network/#/staking/actions).
2. Click **Account Actions**.
3. Locate your bonded account and click **Session Key**.
4. Paste the `author_rotateKeys` result ([generated earlier](#generating-session-keys-with-access-to-the-nodes-rpc-port)) into the provided field.
5. Click **Set Session Key** and submit the extrinsic.

   ![Set session key](images/set-session-keys.png 'Set session key')

Once the transaction succeeds, your session keys will be updated, and you’re ready to proceed with validation.

### Activating Your Operator Node

Before activation, ensure your setup is verified via telemetry. Improper configurations may result in penalties.

#### Steps to Activate

1. In the UI, click **Validate**.

   ![Validate](images/validate.png 'Validate')

2. Enter your reward commission percentage, then click **Validate**. At the time of writing, the maximum allowable commission is capped at 10%.

   ![Validator Preferences](images/validator-preferences.png 'Validator Preferences')

**Congratulations!** Your operator has been added to the queue of available operators. It will be eligible for election in the next operator cycle and will become active at the start of the subsequent era. On Mainnet, this process may take up to 29 hours from setting your validator preferences.

### Stopping Operator Activity

To stop operating as a validator:

1. Navigate to [Staking > Account Actions](https://mainnet-app.polymesh.network/#/staking/actions).
2. Locate your bonded account and click **Stop**.
3. Sign the transaction to confirm the action.

This action removes your node from the list of operators available for election. Your node will remain active for the current era but will become inactive at the start of the next era following its last election.

**Important:** Only stop your node client **after** it is no longer active in the current era to avoid penalties.

#### Reward Destination Considerations

If your reward destination is set to your **Stash account (increase the amount at stake)**:

- Wait until the last reward payment has been received **before** starting the unbonding process.
- Alternatively, change the reward destination to **Stash account (do not increase the amount at stake)** before unbonding to ensure the final reward payment is not bonded.

By doing so, you can ensure the smooth receipt of all pending rewards without interruptions.

## Glossary

| **Term**                 | **Definition**                                                                                                                |
| ------------------------ | ----------------------------------------------------------------------------------------------------------------------------- |
| **Equivocation**         | Occurs when an operator node commits to two or more conflicting states.                                                       |
| **Era**                  | A fixed number of sessions. It determines operator and nominator sets and distributes rewards. (24 hours on Polymesh mainnet) |
| **Immortal transaction** | A transaction valid at any time. (note consideration must be given to account nonce when using an immortal transaction)       |
| **Operator node**        | Permissioned network participants responsible for producing new blocks and finalizing the blockchain.                         |
| **Session**              | A fixed period with a constant set of operators. Operators can only join or leave the set at the start of a session.          |
| **Session keys**         | Keys used by operator nodes to sign consensus-related data.                                                                   |
| **Warm spare node**      | A synced node ready to replace an active operator manually.                                                                   |
