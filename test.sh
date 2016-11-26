#!/bin/bash

for x in cmd collectors handlers server types utils;do
    pushd ${x} >>/dev/null
    go test -coverprofile=coverage.out
    if [ -f coverage.out ];then
        go tool cover -func=coverage.out -o ../resources/coverage/${x}.html
    fi
    popd >>/dev/null
done
