#!/usr/bin/env bash -e

set -x

ROOT_DIR=$(pwd)/
COVERAGE_PATH=${ROOT_DIR}/coverage.txt

echo "" > ${COVERAGE_PATH}

for d in $(go list ./... | grep -v vendor | grep -v gocd-response-links); do
    go test ${TESTARGS} -v $d
    r=$?
    if [ $r -ne 0 ]; then
        exit $r
    elif [ -f profile.out ]; then
        cat profile.out >> ${COVERAGE_PATH}
        rm profile.out
    fi
done
