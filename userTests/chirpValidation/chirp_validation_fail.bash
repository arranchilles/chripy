#!/usr/bin/env bash
curl -i \
    -H "Accept: application/json" \
    -X POST -d '{"body":"Hello my opinion is always correct and no one can ever question me because it is futile and pointless and you will be destroyed for all time and never be invited to anything you pig."}' \
    http://localhost:8080/api/validate_chirp \

