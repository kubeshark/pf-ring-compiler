#!/bin/bash

kernel_version=$(uname -r | awk -F"." '{print $1 "." $2}')
amazon-linux-extras install -y kernel-${kernel_version}
yum install -y kernel-devel-$(uname -r)
make KERNEL_SRC=/usr/src/kernels/$(uname -r)

# BELOW LINES SHOULD BE THE SAME FOR ANY BUILD CONTAINER
cp /PF_RING-8.4.0/kernel/pf_ring.ko /tmp/pf-ring-$(uname -r).ko
echo "Kernel module is ready at: /tmp/pf-ring-$(uname -r).ko"
sleep infinity