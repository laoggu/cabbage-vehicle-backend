#!/usr/bin/env bash
set -euo pipefail

command -v jq >/dev/null 2>&1 || { echo "jq is required but not installed. Aborting." >&2; exit 1; }

HOST=${HOST:-http://localhost:8080}

echo "1. healthz"
curl -s "$HOST/healthz" | grep -q ok

echo "2. wx login"
CODE="mock_code_for_test"
TOK=$(curl -s -X POST "$HOST/v1/auth/wx_login" \
       -d "{\"code\":\"$CODE\"}" | jq -r .access_token)

echo "3. create vehicle"
ID=$(curl -s -X POST "$HOST/v1/vehicles" \
      -H "Authorization: Bearer $TOK" \
      -H "X-Idempotency-Key: ${UUID:-$(uuidgen 2>/dev/null || cat /proc/sys/kernel/random/uuid)}" \
      -d '{"plate":"沪A12345","model":"东风","cold_chain":true,"ice_box_no":"IC001"}' \
     | jq -r .id)

echo "4. get vehicle"
curl -s -H "Authorization: Bearer $TOK" "$HOST/v1/vehicles/$ID" \
  | jq -e ".id==\"$ID\""

echo "5. list vehicles"
curl -s -H "Authorization: Bearer $TOK" "$HOST/v1/vehicles?offset=0&limit=10" \
  | jq -e '.items|length>0'

echo "6. upload url"
curl -s -H "Authorization: Bearer $TOK" "$HOST/v1/vehicles/upload_url?suffix=jpg" \
  | jq -e '.url!=""

echo "7. idempotency replay"
CODE=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$HOST/v1/vehicles" \
          -H "Authorization: Bearer $TOK" \
          -H "X-Idempotency-Key: same-key" \
          -d '{"plate":"沪B11111","model":"解放"}')
[[ $CODE -eq 200 ]] || { echo "idempotency test failed: $CODE"; exit 1; }

echo "All M1 tests passed!"