############# builder
FROM golang:1.16.3 AS builder

WORKDIR /go/src/github.com/gardener-extension-provider-aws/landscaper
COPY . .
RUN make install

############ base
FROM alpine:3.13.2 AS base

############ gardener-extension-provider-aws
FROM base AS landscaper-extension-provider-aws

COPY --from=builder /go/bin/landscaper-extension-provider-aws /landscaper-extension-provider-aws
ENTRYPOINT ["/landscaper-extension-provider-aws"]
