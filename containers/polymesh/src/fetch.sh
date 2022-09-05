#!/bin/bash

set -e
set -x
set -o pipefail

apt-get --quiet update
apt-get --quiet --assume-yes install curl wget unzip
curl -s https://api.github.com/repos/polymeshassociation/polymesh/releases/latest | \
  grep "browser_download_url" | grep -v "wasm" | \
  awk '{print $2}' | \
  xargs wget --quiet
unzip *.zip
rm *.zip
mv polymesh-* polymesh
mkdir -p /opt/files/usr/local/bin
mkdir -p /opt/files/var/lib/polymesh
mkdir -p /opt/files/lib/x86_64-linux-gnu

mv polymesh /opt/files/usr/local/bin/

