#!/bin/bash

govendor fetch +missing
echo "> govendor remove +unused"
govendor remove +unused
echo "> govendor sync"
govendor sync

gocov test ./... || exit 0
echo ">> Overall"
go test -covermode=count -coverprofile=count.out fmt
