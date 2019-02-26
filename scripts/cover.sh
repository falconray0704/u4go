#!/usr/bin/env bash

#set -x
set -e
echo "" > cover.out


for d in $(go list $@); do
#    go test -race -coverprofile=profile.out $d
    go test -coverprofile=profile.out $d
    if [ -f profile.out ]; then
        go tool cover -html=profile.out -o=cover.html
        cat profile.out >> cover.out
        rm profile.out
    fi
done
