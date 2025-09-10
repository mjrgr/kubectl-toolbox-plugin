############################
# Build executable binary
############################
FROM golang:1.25-alpine AS builder

ENV APP_SRC=${GOPATH}/src/app

RUN mkdir -p ${APP_SRC}

WORKDIR ${APP_SRC}

# Get dependencies
COPY go.mod .
COPY go.sum .
RUN go mod download
RUN go mod verify
RUN go mod tidy

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix nocgo -o /tmp/kubectl-toolbox-plugin

#############################
## Build a small image
#############################
FROM scratch

LABEL org.opencontainers.image.base.name="scratch"
LABEL org.opencontainers.image.ref.name="kubectl-toolbox-plugin"
LABEL org.opencontainers.image.title="kubectl-toolbox-plugin"
LABEL org.opencontainers.image.description="Lightweight init container tool for Kubernetes checks"
LABEL org.opencontainers.image.authors="Mehdi Jr-Gr"
LABEL org.opencontainers.image.vendor="Mehdi Jr-Gr"
LABEL org.opencontainers.image.source="https://github.com/mjrgr/kubectl-toolbox-plugin"
LABEL org.opencontainers.image.licenses="Apache-2.0"

# Copy our static executable
COPY --from=builder --chmod=0755 \
    /tmp/kubectl-toolbox-plugin /kubectl-toolbox-plugin

# Copy essential files for networking and TLS
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group

# Copy DNS resolver configuration (important for hostname resolution)
COPY --from=builder /etc/resolv.conf /etc/

USER 1000:1000

ENTRYPOINT ["/kubectl-toolbox-plugin"]
CMD [ "--help" ]


