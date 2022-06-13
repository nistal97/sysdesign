#!/bin/sh
export GO111MODULE=on
go mod tidy
go mod vendor

#NWaySetAssocCache
go build -v -o ./output/cache ./NWaySetAssocCache
go test -v NWaySetAssocCache/cache_test.go


