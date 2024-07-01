FROM debian:bookworm

# Updating packages and installing dependencies
RUN apt-get update && apt-get install -y \
    apt-transport-https \
    curl \
    make \
    gnupg \
    unzip \
    git \
    sudo \
    ca-certificates \
    lsb-release \
    procps \
    htop \
    && rm -rf /var/lib/apt/lists/*

# Installing Go
ENV GO_VERSION=1.22.0
RUN curl -OL https://golang.org/dl/go${GO_VERSION}.linux-amd64.tar.gz \
    && tar -C /usr/local -xzf go${GO_VERSION}.linux-amd64.tar.gz \
    && rm go${GO_VERSION}.linux-amd64.tar.gz
ENV PATH=$PATH:/usr/local/go/bin

# Add Docker's official GPG key: \
RUN install -m 0755 -d /etc/apt/keyrings \
    && curl -fsSL https://download.docker.com/linux/debian/gpg -o /etc/apt/keyrings/docker.asc \
    && chmod a+r /etc/apt/keyrings/docker.asc \
    # Add the repository to Apt sources:
    && echo \
    "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.asc] https://download.docker.com/linux/debian bookworm stable" | \
    tee /etc/apt/sources.list.d/docker.list > /dev/null \
    && apt-get update && apt-get install -y \
    # docker-ce \
    docker-ce-cli \
    # containerd.io \
    # docker-buildx-plugin \
    docker-compose-plugin

# Add a new non-root user
RUN useradd -m devuser && echo "devuser:devuser" | chpasswd && adduser devuser sudo \
    && echo 'devuser ALL=(ALL) NOPASSWD:ALL' > /etc/sudoers.d/devuser \
    && chmod 0440 /etc/sudoers.d/devuser

# Switch to the non-root user
USER devuser
WORKDIR /home/devuser/
