FROM registry.gitlab.steamos.cloud/steamrt/sniper/sdk:latest

ARG GO_VERSION=1.25.5

ENV GOPATH=/app/.go
ENV GOMODCACHE=/app/.go/pkg/mod
ENV GOCACHE=/app/.go/cache
ENV PATH=$PATH:/usr/local/go/bin:$GOPATH/bin

WORKDIR /app

RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates curl git pkg-config build-essential \
    gcc-14-monolithic \
    && rm -rf /var/lib/apt/lists/*

RUN ln -sf /usr/bin/gcc-14 /usr/bin/gcc || true \
 && ln -sf /usr/bin/g++-14 /usr/bin/g++ || true

RUN curl -fsSL https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz \
    | tar -C /usr/local -xz

COPY docker-build.sh /usr/local/bin/docker-build.sh
RUN chmod +x /usr/local/bin/docker-build.sh

ENTRYPOINT ["/usr/local/bin/docker-build.sh"]
