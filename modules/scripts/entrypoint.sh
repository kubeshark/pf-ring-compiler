#!/bin/sh

# Get the current kernel version
kernel_version=$(uname -r)

# Check if the module directory exists for the current kernel
module_path="/opt/modules/pf-ring-${kernel_version}.ko"

if [ -f "$module_path" ]; then
    # Check if the module is already loaded
    if ! lsmod | grep -q pf_ring; then
        echo "Loading pf_ring module for kernel ${kernel_version}"
        module_load_path="/opt/modules/pf-ring-${kernel_version}.ko"
        mkdir /opt/lib/modules/${kernel_version} -p
        cp $module_path ${module_load_path}
        depmod -b /opt ${kernel_version}
        insmod ${module_load_path}
        if [ $? -ne 0 ]; then
            echo "Failed to load module from ${module_load_path}"
            echo "Falling back to AF_PACKET"
            exit 0
        fi
    else
        echo "pf_ring module is already loaded for kernel ${kernel_version}"
        exit 0
    fi
else
    echo "No pf_ring module found for the current kernel version ${kernel_version}"
    echo "Falling back to AF_PACKET"
fi