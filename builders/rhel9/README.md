# rhel9

This builder is used to build the kernel module for Kubernetes workers running RHEL9/Rocky Linux 9.
It is based on Rocky Linux 9 container image.
`ubi9/ubi` container image is not used as build runtime due to `kernel-devel` package repository is not available without Red Hat subscription.

# build

```
docker build -t kubehq/pf-ring-builder:rhel9 .
docker push kubehq/pf-ring-builder:rhel9
```