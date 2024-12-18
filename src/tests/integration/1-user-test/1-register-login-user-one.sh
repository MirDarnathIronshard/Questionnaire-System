#!/bin/bash

REGISTER_URL="http://127.0.0.1:5005/api/register"
LOGIN_URL="http://127.0.0.1:5005/api/login"
TOKEN_FILE="../test-user-one.sh"
VAR_NAME="user_1"

EMAIL="test@test.com"
GENDER="male"
NATIONAL_ID="1360718583"
PASSWORD="123456"

REGISTER_RESPONSE=$(curl -s -X POST \
  "$REGISTER_URL" \
  -H 'accept: */*' \
  -H 'Content-Type: application/json' \
  -d "{
    \"email\": \"$EMAIL\",
    \"gender\": \"$GENDER\",
    \"national_id\": \"$NATIONAL_ID\",
    \"password\": \"$PASSWORD\"
  }")

echo "$REGISTER_RESPONSE" | jq

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
