# PF_RING compiler

The PF_RING compiler is a tool designed to compile the PF_RING kernel module for specific target distributions.
It is implemented in Go and uses Kubernetes to run the compilation jobs.

# Overview

The tool is structured into two main packages: `cmd` and `pkg/compiler`.

The `cmd`` package contains the command-line interface (CLI) for the tool. It defines the flags that can be passed to the tool, the runner that executes the main logic, and the main entry point of the application.

The `pkg/compiler` package contains the logic for creating, monitoring, and cleaning up the Kubernetes jobs that compile the PF_RING kernel module.

# Usage

To use the PF_RING Compiler, you need to run the main application with the `--target`` flag specifying the target for which the PF_RING kernel module should be compiled.
The target in this context is Linux distribution, used in Kubernetes cluster, where Kubeshark is expected to be used.
The supported targets are:
- al2 (Amazon Linux 2)
- ...

Here is an example of how to run the tool:

```bash
./pfring-compiler --target al2
```

# Build containers

## Overview

The `getCompileContainerImage` function is a helper function that maps the target specified by the user to the corresponding Docker image that should be used for the compilation job.

The function is defined in the `pkg/compiler/compiler.go` file and takes a single argument, `target`, which is a string representing the target for which the PF_RING kernel module should be compiled.

The function uses a map to associate each supported target with its corresponding Docker image.

| Target | Docker Image |
|--------|--------------|
| al2 | coreset/build:kubeshark-pf-ring-al2-builder |

Build containers Dockerfiles are defined in `builders` folder.

## Adding new build containers

To add the build container for the new target, follow the example in [al2 Dockerfile](builders/al2/Dockerfile).
The build container script in general consists does 3 steps:
1. Install kernel headers for the current kernel version.
2. Build pf_ring module using kernel headers for the current kernel version.
3. Put pf_ring.ko file into /tmp/pf-ring-<kernel version>.ko file.

Steps 1 and 2 are specific to the target distribution.
Step 3 should be implemented same way for all the build containers by adding lines below into `entrypoint.sh` script:

```bash
cp /PF_RING-8.4.0/kernel/pf_ring.ko /tmp/pf-ring-$(uname -r).ko
echo "Kernel module is ready at: /tmp/pf-ring-$(uname -r).ko"
sleep infinity
```

These 3 lines:
- set predictable path for the PF_RING kernel module path
- set predictable log line to determine the end of the build process
- give enough time for CLI to copy kernel module from pod to the local file system

After build container is ready:
1. Build container
2. Push into container registry
3. Add new target-to-container mapping into `getCompileContainerImage` function in the `pkg/compiler/compiler.go` file.
4. Create new PR, merge into `main` branch.
5. Create new Github release

# Building CLI tool

To build the tool, you need to have Go installed. You can then use the `go build` command.
