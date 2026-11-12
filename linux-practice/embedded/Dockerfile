FROM ubuntu:24.04

ENV DEBIAN_FRONTEND=noninteractive
ENV container=docker

RUN apt-get update && apt-get install -y --no-install-recommends \
    systemd \
    systemd-sysv \
    dbus \
    sudo \
    curl \
    iproute2 \
    procps \
    net-tools \
    iputils-ping \
    dnsutils \
    ufw \
    vim \
    nano \
    less \
    ca-certificates \
    python3 \
    locales \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

RUN locale-gen en_US.UTF-8
ENV LANG=en_US.UTF-8
ENV LC_ALL=en_US.UTF-8

# Configure student user with passwordless sudo
RUN useradd -m -s /bin/bash student && \
    echo "student:student" | chpasswd && \
    usermod -aG sudo student && \
    echo "student ALL=(ALL) NOPASSWD:ALL" > /etc/sudoers.d/student && \
    chmod 0440 /etc/sudoers.d/student

STOPSIGNAL SIGRTMIN+3

CMD ["/lib/systemd/systemd"]
