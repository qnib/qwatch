#!/bin/bash

for x in cmd collectors handlers server types utils;do
    pushd ${x} >>/dev/null
    go test -coverprofile=coverage.out
    popd >>/dev/null
done
