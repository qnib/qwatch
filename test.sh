#!/bin/bash

govendor fetch +missing
echo "> govendor remove +unused"
govendor remove +unused
echo "> govendor sync"
govendor sync
if [ ! -d resources/coverity ];then
    mkdir -p resources/coverity
fi
go test -coverprofile=coverage.out
for x in $(find . -maxdepth 1 -type d |egrep -v "(\.$|\.git|vendor|bin|resources)");do
    go test -coverprofile=resources/coverity/${x}.out ${x} >>/dev/null
done
coveraggregator -o resources/coverity/coverage-all.out $(find . -name coverage.out |egrep -v "(\.$|\.git|vendor|bin)") >>/dev/null
go tool cover -func=resources/coverity/coverage-all.out |tee ./coverage-all.out
go tool cover -html=resources/coverity/coverage-all.out -o resources/coverity/all.html
