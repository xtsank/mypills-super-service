#!/usr/bin/env bash
set -euo pipefail

base_url=${E2E_BASE_URL:-http://127.0.0.1:18081}
illness_id=00000000-0000-4000-8000-000000000041
safe_medicine_id=00000000-0000-4000-8000-000000000042
unsafe_medicine_id=00000000-0000-4000-8000-000000000043
response=$(mktemp)
trap 'rm -f "$response"' EXIT

request() {
  local expected=$1 method=$2 endpoint=$3 body=${4:-} token=${5:-}
  local status
  local -a options=(--silent --show-error --max-time 10 --noproxy '*' --output "$response" --write-out '%{http_code}' --request "$method")
  if [[ -n $body ]]; then
    options+=(--header 'Content-Type: application/json' --data "$body")
  fi
  if [[ -n $token ]]; then
    options+=(--header "Authorization: Bearer $token")
  fi
  status=$(curl "${options[@]}" "$base_url$endpoint")
  if [[ $status != "$expected" ]]; then
    printf '%s %s: expected HTTP %s, got %s\n' "$method" "$endpoint" "$expected" "$status" >&2
    cat "$response" >&2
    exit 1
  fi
}

echo '1/5 Register the driver'
request 201 POST /api/v2/users '{"login":"e2e-driver","email":"e2e-driver@example.test","password":"strongpass123","age":30,"weight":70,"is_pregnant":false,"is_driver":true,"allergies":[],"illnesses":[]}'
token=$(jq -er '.token | select(length > 0)' "$response")

echo '2/5 Create a cabinet'
request 201 POST /api/v2/cabinets '{"name":"E2E home cabinet"}' "$token"
cabinet_id=$(jq -er '.id | select(length > 0)' "$response")

manufactured=$(date -u -d '7 days ago' +%F)
for medicine_id in "$safe_medicine_id" "$unsafe_medicine_id"; do
  if [[ $medicine_id == "$safe_medicine_id" ]]; then
    echo '3/5 Add the driver-safe medicine'
  else
    echo '4/5 Add the driver-unsafe medicine'
  fi
  body=$(jq -nc --arg medicine_id "$medicine_id" --arg manufactured "$manufactured" \
    '{medicine_id:$medicine_id,date_of_manufacture:$manufactured,quantity:10,minimum_quantity:1}')
  request 201 POST "/api/v2/cabinets/$cabinet_id/items" "$body" "$token"
  item_id=$(jq -er '.id | select(length > 0)' "$response")
  if [[ $medicine_id == "$safe_medicine_id" ]]; then
    safe_item_id=$item_id
  fi
done

echo '5/5 Request medicine recommendations'
request 200 GET "/api/v2/recommendations?illness_id=$illness_id&cabinet_id=$cabinet_id" '' "$token"
jq -e --arg medicine_id "$safe_medicine_id" --arg cabinet_id "$cabinet_id" --arg item_id "$safe_item_id" '
  .items | length == 1 and
  .[0].medicine_id == $medicine_id and
  .[0].cabinet_id == $cabinet_id and
  .[0].cabinet_item_id == $item_id and
  .[0].name == "Driver-safe medicine" and
  .[0].dosage == 1 and
  .[0].unit_name == "piece" and
  .[0].frequency_per_day == 2 and
  .[0].quantity_in_cabinet == 10
' "$response" >/dev/null
echo 'The E2E scenario passed: only the driver-safe medicine was recommended.'
