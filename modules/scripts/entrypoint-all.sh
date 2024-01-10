#!/bin/sh

# Get the current kernel version
current_kernel_version=$(uname -r)

# Check if the module directory exists for the current kernel
module_path="/opt/lib/modules/${current_kernel_version}/pf_ring.ko"

if [ -f "$module_path" ]; then
    # Check if the module is already loaded
    if ! lsmod | grep -q pf_ring; then
        echo "Loading pf_ring module for kernel ${current_kernel_version}"
        insmod $module_path
    else
        echo "pf_ring module is already loaded for kernel ${current_kernel_version}"
        exit 0
    fi
else
    echo "No pf_ring module found for the current kernel version ${current_kernel_version}"
    echo "Falling back to libpcap"
fi