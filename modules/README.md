# modules

This folder contains ready for KMM usage kernel modules.

| Kernel version | Container |
|----------------|-----------|
|5.10.198-187.748.amzn2.x86_64|kubehq/pf-ring-module:5.10.198-187.748.amzn2.x86_64|
|5.10.199-190.747.amzn2.x86_64|kubehq/pf-ring-module:5.10.199-190.747.amzn2.x86_64|
|5.14.0-362.8.1.el9_3.x86_64|kubehq/pf-ring-module:5.14.0-362.8.1.el9_3.x86_64|

# build

Where there is module for a new kernel version available, run:

1. Copy PF_RING kernel moduele into `modules` directory with name `pf-ring-<kernel version>.ko``

2. Run build

```
KERNEL_VERSION=<new kernel version>
docker build --build-arg KERNEL_VERSION=${KERNEL_VERSION} -t kubehq/pf-ring-module:${KERNEL_VERSION} .
docker push kubehq/pf-ring-module:${KERNEL_VERSION}
```

3. Update the table with supported kernel versions