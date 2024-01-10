# modules

Folder structure:

- `ko/` - contains all available `pf_ring.ko` built by Kubeshark for different kernel versions
- `scripts/` - contains scripts used to build container with all modules in `ko` folder
- `Dockerfile.all` - Dockerfile for building container with all modules in `ko` folder
- `Dockerfile.single` - Dockerfile for building container with `pf_ring.ko` for selected kernel version.

## Build

### Container with all available modules

```bash
docker build -t kubeshark/pf-ring-module:all -f Dockerfile.all
```

### Container for specific PF_RING kernel version

This requires target kernel module to exist at `ko/pf-ring-<kernel-version>.ko`.

```bash
kernel_version=5.15.0-1050-aws # example
docker build --build-arg KERNEL_VERSION=${kernel_version} -t kubeshark/pf-ring-module:${kernel_version}  -f Dockerfile.single
```

## Usage

Kubeshark maintains `kubeshark/pf-ring-module:all` and `kubeshark/pf-ring-module:${kernel_version}` containers with modules available under `ko`.
These containers are used in Kubeshark Helm chart.
Please refer to [Documentation](https://github.com/kubeshark/kubeshark/tree/master/helm-chart) for more details.
