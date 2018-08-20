# Build the manager binary
FROM golang:1.10.3 as builder

# Copy in the go src
WORKDIR /go/src/github.com/mengqiy/webhook-example/
COPY example/    example/
COPY vendor/ vendor/

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o webhook github.com/mengqiy/webhook-example/example/

# Copy the controller-manager into a thin image
FROM ubuntu:latest
WORKDIR /root/
COPY --from=builder /go/src/github.com/mengqiy/webhook-example/webhook .
ENTRYPOINT ["./webhook"]

