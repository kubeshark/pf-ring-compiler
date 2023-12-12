#!/bin/bash

apt install -y linux-headers-$(uname -r)
make KERNEL_SRC=/usr/src/linux-headers-$(uname -r)

# BELOW LINES SHOULD BE THE SAME FOR ANY BUILD CONTAINER
cp /PF_RING-8.4.0/kernel/pf_ring.ko /tmp/pf-ring-$(uname -r).ko
echo "Kernel module is ready at: /tmp/pf-ring-$(uname -r).ko"
sleep infinity
