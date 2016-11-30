#!/bin/bash

govendor fetch +missing
echo "> govendor remove +unused"
govendor remove +unused
echo "> govendor sync"
govendor sync

gocov test ./... || exit 0
