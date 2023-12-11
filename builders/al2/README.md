# al2

This builder is used to build the kernel module for default EKS workers.
It is based on Amazon Linux 2.

# build

```
docker build -t kubehq/pf-ring-builder:al2 .
docker push kubehq/pf-ring-builder:al2
```