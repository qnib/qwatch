#!/bin/bash
set -e

govendor fetch +missing
echo "> govendor remove +unused"
govendor remove +unused
echo "> govendor fetch +missing"
govendor fetch +missing
echo "> govendor update +local"
govendor update +local
#echo "> go get github.com/stretchr/testify/assert"
#go get github.com/stretchr/testify/assert
#echo "> go get -d"
#go get -d
if [ ! -d resources/coverity ];then
    mkdir -p resources/coverity
fi
echo "> go test -cover -coverprofile=resources/coverity/qwatch.cover"
go test -cover -coverprofile=resources/coverity/qwatch.cover >>/dev/null
COVER_FILES="resources/coverity/qwatch.cover"
for x in $(find . -maxdepth 1 -type d |egrep -v "(\.$|\.git|vendor|bin|resources)");do
    echo "> go test -cover -coverprofile=resources/coverity/${x}.cover ${x}"
    go test -cover -coverprofile=resources/coverity/${x}.cover ${x} >>/dev/null
    COVER_FILES="${COVER_FILES} resources/coverity/${x}.cover"
done
echo "> coveraggregator -o resources/coverity/coverage-all.out ${COVER_FILES}"
coveraggregator -o resources/coverity/coverage-all.out ${COVER_FILES} >>/dev/null
#go tool cover -func=coverage-all.outcover |tee ./resources/coverity/coverage-all.out
#go tool cover -html=coverage-all.out -o resources/coverity/all.html
