#!/usr/bin/env bash

ENDPOINT="http://127.0.0.1:8153/go/api/admin/config.xml"
function get_status {
    curl -H 'Accept: application/vnd.go.cd.v3+json' \
        --write-out %{http_code} \
        --silent \
        --output /dev/null \
        ${ENDPOINT}
}

counter=0
wait_length=5
elapsed=0

while [ $counter -lt 30 ]; do

    code=$(get_status)
    if [ "200" == "$code" ]; then
        echo "Got status ${code}. Exiting."
        exit 0
    fi

    if [ "$elapsed" == "120" ]; then
        cat godata/server/logs/*.log
        curl -H 'Accept: application/vnd.go.cd.v3+json' \
            ${ENDPOINT}
    fi

    echo "Got status ${code}. Elapsed: '${elapsed}' seconds."

    sleep "${wait_length}"
    elapsed=$((elapsed+wait_length))

done

exit 1