#!/usr/bin/env bash

curl -i \
    -X POST \
    -H "Accept:application/json" \
    -H "Authorization: Bearer bc95535487b3ac43b64b0bbf3c66e254a3934ea0017ef254c45f41826f10ef61" \
    -d '{"email":"user@example.com"}' \
    http://localhost:8080/api/refresh