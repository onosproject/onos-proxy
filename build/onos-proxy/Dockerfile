# SPDX-FileCopyrightText: 2021-present Open Networking Foundation <info@opennetworking.org>
#
# SPDX-License-Identifier: Apache-2.0

FROM onosproject/golang-build:v1.3.0 as build

ENV GO111MODULE=on
COPY . /go/src/github.com/onosproject/onos-proxy
RUN cd /go/src/github.com/onosproject/onos-proxy && GOFLAGS=-mod=vendor make build

FROM alpine:3.11
RUN apk add libc6-compat

USER nobody

COPY --from=build /go/src/github.com/onosproject/onos-proxy/build/_output/onos-proxy /usr/local/bin/onos-proxy

ENTRYPOINT ["onos-proxy"]
