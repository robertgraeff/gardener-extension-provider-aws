############# builder
FROM golang:1.16.3 AS builder

WORKDIR /go/src/github.com/gardener/gardener-extension-provider-aws/landscaper
COPY . .
RUN make install

############# base
FROM alpine:3.13.2 AS base

############# deployer-aws
FROM base AS provider-aws-deployer

COPY --from=builder /go/bin/landscaper-extension-provider-aws /provider-aws-deployer
ENTRYPOINT ["/provider-aws-deployer"]
