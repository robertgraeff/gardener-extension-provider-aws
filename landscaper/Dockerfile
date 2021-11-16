############# builder
FROM golang:1.16.3 AS builder

WORKDIR /go/src/github.com/gardener-extension-provider-aws/landscaper
COPY . .
RUN make install

############ base
FROM alpine:3.13.2 AS base

############ gardener-extension-provider-aws
FROM base AS aws-deployer

COPY --from=builder /go/bin/landscaper-extension-provider-aws /aws-deployer
ENTRYPOINT ["/aws-deployer"]
