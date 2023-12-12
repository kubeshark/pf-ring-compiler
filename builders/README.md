# Builders

This folder contains specific Dockerfiles, used to build pf_ring.ko kernel module for a target kernel.

Supported builders:
- [al2](al2/README.md)
- [rhel9](rhel9/README.md)
- [ubuntu](ubuntu/README.md)

# Manual pf_ring.ko build

If there is no available automated builder, PF_RING can be compiled manually following steps:

1. Run debug container on the target node

Select a proper debug container for the target node 
(e.g. if the node is running Ubuntu - use Ubuntu container)

```bash
kubectl debug node/<node name> -it --attach=true --image=<container>
```

2. Install kernel headers for the respective kernel version

The installation process is different for different distributions.
Follow the examples in Dockerfiles of the automated builders.
After installations is completed, the most common kernel headers path is `/usr/src/kernels/<kernel version>`

3. Download PF_RING stable release

```
wget https://github.com/ntop/PF_RING/archive/refs/tags/8.4.0.tar.gz && \
tar -xf 8.4.0.tar.gz
cd PF_RING-8.4.0/kernel/ 
```

4. Run kernel module compilation

```
make KERNEL_SRC=/usr/src/kernels/$(uname -r)
pwd # print kernel module path
```

5. Copy pf_ring.ko to local filesystem

```
kubectl debug pod/<debug pod name>/<kernel module path>/pf_ring.ko pf_ring.ko
```
