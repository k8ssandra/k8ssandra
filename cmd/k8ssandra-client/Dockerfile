# Build the binary
FROM golang:1.15 as builder

WORKDIR /workspace

# Go Modules
COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

# Copy the go source
COPY cmd/k8ssandra-client/main.go main.go
COPY pkg/ pkg/

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -a -o k8ssandra-client main.go

# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=builder /workspace/k8ssandra-client .
USER nonroot:nonroot

ENTRYPOINT ["/k8ssandra-client"]
