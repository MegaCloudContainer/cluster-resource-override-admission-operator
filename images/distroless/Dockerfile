# Build the manager binary
FROM golang:1.20.10 as builder

ENV CGO_ENABLED 0
ENV GOOS linux

WORKDIR /workspace
# Copy the Go Modules manifests
COPY . .

# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

# Build
RUN go build -ldflags="-s -w" -a -o cluster-resource-override-admission-operator cmd/cluster-resource-override-admission-operator/main.go && echo "Building GOARCH of $GOARCH.."

# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM gcr.io/distroless/base:nonroot
WORKDIR /
COPY --from=builder /workspace/cluster-resource-override-admission-operator /bin/cluster-resource-override-admission-operator
USER 65532:65532
