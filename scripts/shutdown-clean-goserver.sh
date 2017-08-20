#!/usr/bin/env bash -x

/bin/kill $(ps aux | grep java | head -n 1 | cut -d ' ' -f 4)

rm -rf /godata/*