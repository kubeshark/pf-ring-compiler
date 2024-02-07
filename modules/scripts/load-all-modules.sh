#!/bin/sh

apt update
apt install kmod -y

for f in $(ls /tmp/modules/*.ko); do
    kernel_version=$(echo $f | sed -E 's/.*pf-ring-(.*)\.ko/\1/');
    mkdir /opt/lib/modules/${kernel_version} -p
    mv ${f} /opt/lib/modules/${kernel_version}/pf_ring.ko
    echo "Loading module for ${kernel_version} kernel"
    depmod -b /opt ${kernel_version}
done