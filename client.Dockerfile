FROM golang:1.24-alpine AS builder

RUN apk add --no-cache ca-certificates

# set up nsswitch.conf for Go's "netgo" implementation
# - https://github.com/golang/go/blob/go1.9.1/src/net/conf.go#L194-L275
# - docker run --rm debian:stretch grep '^hosts:' /etc/nsswitch.conf
# see more https://github.com/golang/go/issues/22846
RUN echo "hosts: files dns" > /etc/nsswitch.conf

ENV GOPATH=/go
ENV PATH=$GOPATH/bin:/usr/local/go/bin:$PATH
ENV WORKDIR=/opt/nti

RUN mkdir -p $WORKDIR
COPY ./ $WORKDIR
WORKDIR $WORKDIR

RUN go build -o ./client cmd/client/main.go


# Final repo 
FROM alpine:3.21.3 AS final

# set up nsswitch.conf for Go's "netgo" implementation
# - https://github.com/golang/go/blob/go1.9.1/src/net/conf.go#L194-L275
# - docker run --rm debian:stretch grep '^hosts:' /etc/nsswitch.conf
# see more https://github.com/golang/go/issues/22846
RUN echo "hosts: files dns" > /etc/nsswitch.conf

ENV WORKDIR=/opt/nti

COPY --from=builder /opt/nti/client ${WORKDIR}/client

WORKDIR $WORKDIR
