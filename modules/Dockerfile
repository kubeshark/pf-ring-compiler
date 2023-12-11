FROM alpine:3.18
ARG KERNEL_VERSION

COPY pf-ring-${KERNEL_VERSION}.ko /opt/lib/modules/${KERNEL_VERSION}/pf_ring.ko
RUN apk add kmod

RUN depmod -b /opt ${KERNEL_VERSION}