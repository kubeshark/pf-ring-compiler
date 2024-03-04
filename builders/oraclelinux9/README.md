# oraclelinux9

This builder is used to build the kernel module for Kubernetes workers running Oracle Linux 9(OpenShift platform).
It is based on `oraclielinux:9` container image.

# build

```
docker build -t kubehq/pf-ring-builder:oraclelinux9 .
docker push kubehq/pf-ring-builder:oraclelinux9
```
