#!/bin/bash -e
GOPATH=$HOME
go generate ./...
go test -cover -coverprofile /tmp/c.out
