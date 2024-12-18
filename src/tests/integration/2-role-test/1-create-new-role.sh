#!/bin/bash
source ../test-user-one.sh
source ../test-user-two.sh
source ../super-admin.sh
source ../admin.sh


# shellcheck disable=SC2154
TOKEN="$super_admin"
#TOKEN="$user_1"
#TOKEN="$user_2"
#TOKEN="$admin"

curl -X 'POST' \
  'http://127.0.0.1:5005/api/roles' \
  -H 'accept: application/json' \
  -H "Authorization: Bearer $TOKEN" \
  -H 'Content-Type: application/json' \
  -d '{
  "name": "efrefe"
}' |jq

