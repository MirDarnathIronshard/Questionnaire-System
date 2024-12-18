#!/bin/bash


curl -X 'POST' \
  'http://127.0.0.1:5005/api/register' \
  -H 'accept: */*' \
  -H 'Content-Type: application/json' \
  -d '{
  "email": "test@test.com",
  "gender": "male",
  "national_id": "",
  "password": "123456"
}' | jq


curl -X 'POST' \
  'http://127.0.0.1:5005/api/register' \
  -H 'accept: */*' \
  -H 'Content-Type: application/json' \
  -d '{
  "email": "test@test.com",
  "gender": "male",
  "national_id": "1362056170",
  "password": ""
}' | jq