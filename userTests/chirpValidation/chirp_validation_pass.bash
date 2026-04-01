#!/usr/bin/env bash

curl -i \
    -H "Accept: application/json" \
    -X POST -d '{"body": "Good Morning Vietnam!"}' \
    http://localhost:8080/api/validate_chirp
