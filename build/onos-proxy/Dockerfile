FROM onosproject/golang-build:v0.6.3 as build

ENV GO111MODULE=on
COPY . /go/src/github.com/onosproject/onos-proxy
RUN cd /go/src/github.com/onosproject/onos-proxy && GOFLAGS=-mod=vendor make build

FROM alpine:3.11
RUN apk add libc6-compat

USER nobody

COPY --from=build /go/src/github.com/onosproject/onos-proxy/build/_output/onos-proxy /usr/local/bin/onos-proxy

ENTRYPOINT ["onos-proxy"]
