
FROM ubuntu:24.04

ENV GOLANG_VERSION=1.25.5

# System setup
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
    curl \
    git \
    jq \
    wget \
    ca-certificates \
    postgresql-client && \
    rm -rf /var/lib/apt/lists/*

# Go setup
RUN wget -q https://go.dev/dl/go${GOLANG_VERSION}.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go${GOLANG_VERSION}.linux-amd64.tar.gz && \
    rm go${GOLANG_VERSION}.linux-amd64.tar.gz

ENV PATH="/usr/local/go/bin:${PATH}"

# GO CLIs
RUN go install github.com/bootdotdev/bootdev@latest && \
    go install github.com/pressly/goose/v3/cmd/goose@latest && \
    go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest && \
    echo 'export PATH=$PATH:$HOME/go/bin' >> ~/.bashrc
