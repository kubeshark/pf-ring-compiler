FROM amazonlinux:2.0.20230628.0

RUN yum install -y gcc gcc-c++ kernel-devel make which git wget tar gzip && \
    wget https://github.com/ntop/PF_RING/archive/refs/tags/8.4.0.tar.gz && \
    tar -xf 8.4.0.tar.gz

WORKDIR /PF_RING-8.4.0/kernel/
COPY pf-ring-al2-build.sh pf-ring-al2-build.sh

ENTRYPOINT ["/PF_RING-8.4.0/kernel/pf-ring-al2-build.sh"]