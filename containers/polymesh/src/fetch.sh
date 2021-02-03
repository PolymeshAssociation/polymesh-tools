#!/bin/bash

set -e
set -x
set -o pipefail

apt-get --quiet update
apt-get --quiet --assume-yes install curl wget
curl -s https://api.github.com/repos/polymathnetwork/polymesh/releases/latest | \
  grep "browser_download_url.*linux-amd64" | \
  awk '{print $2}' | \
  xargs wget --quiet
sha256sum -c *.sha256
tar xvzf *.tgz
mkdir -p /opt/files/usr/local/bin
mkdir -p /opt/files/var/lib/polymesh
mkdir -p /opt/files/lib/x86_64-linux-gnu

mv $(find . -type f -executable -print) /opt/files/usr/local/bin/polymesh 
touch /opt/files/var/lib/polymesh/.keep
cp -a /lib/x86_64-linux-gnu/* /opt/files/lib/x86_64-linux-gnu/
LDLIBS=$(ldd /opt/files/usr/local/bin/polymesh | grep -o '/\S*')
for LIB in $LDLIBS; do
    mkdir -p /opt/files/$(dirname $LIB | cut -c 2-)
    cp $LIB  /opt/files/$(dirname $LIB | cut -c 2-)/
done

