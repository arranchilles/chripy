#!/usr/bin/env bash

curl -i \
    -X POST \
    -H "Accept:application/json" \
    -d '{"email":"user@example.com"}' \
    http://localhost:8080/api/users