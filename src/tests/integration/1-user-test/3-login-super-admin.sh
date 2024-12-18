#!/bin/bash

LOGIN_URL="http://127.0.0.1:5005/api/login"
TOKEN_FILE="../super-admin.sh"
VAR_NAME="super_admin"

EMAIL="super_admin@exampel.com"
PASSWORD="admin"



LOGIN_RESPONSE=$(curl -s -X POST \
  "$LOGIN_URL" \
  -H 'accept: */*' \
  -H 'Content-Type: application/json' \
  -d "{
    \"email\": \"$EMAIL\",
    \"password\": \"$PASSWORD\"
  }")

echo "$LOGIN_RESPONSE" | jq

TOKEN=$(echo "$LOGIN_RESPONSE" | jq -r '.data')

if [ "$TOKEN" != "null" ] && [ -n "$TOKEN" ]; then

  echo "$VAR_NAME=\"$TOKEN\"" > "$TOKEN_FILE"

else
  echo "error save file"
  exit 1
fi
