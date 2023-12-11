FROM rockylinux:9

RUN yum install -y gcc gcc-c++ make which git wget tar gzip && \
    wget https://github.com/ntop/PF_RING/archive/refs/tags/8.4.0.tar.gz && \
    tar -xf 8.4.0.tar.gz

WORKDIR /PF_RING-8.4.0/kernel/
COPY entrypoint.sh entrypoint.sh

ENTRYPOINT ["/PF_RING-8.4.0/kernel/entrypoint.sh"]