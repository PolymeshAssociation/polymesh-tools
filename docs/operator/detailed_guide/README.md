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
