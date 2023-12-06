#!/bin/bash

kernel_version=$(uname -r | awk -F"." '{print $1 "." $2}')
amazon-linux-extras install -y kernel-${kernel_version}
yum install -y kernel-devel-$(uname -r)
make KERNEL_SRC=/usr/src/kernels/$(uname -r)

echo "Kernel module is ready."
echo "Run from terminal: kubectl cp $(hostname):/PF_RING-8.4.0/kernel/pf_ring.ko pf_ring.ko"

sleep infinity