#!/bin/bash
BASE="http://localhost:8080"
TOKEN="testplayer"

echo "=== 1. Health ==="
curl -s $BASE/health | jq .

echo -e "\n=== 2. Login (new user) ==="
curl -s -X POST $BASE/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d "{\"token\":\"$TOKEN\"}" | jq .

echo -e "\n=== 3. Login (repeat) ==="
curl -s -X POST $BASE/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d "{\"token\":\"$TOKEN\"}" | jq .

echo -e "\n=== 4. Benefits ==="
curl -s $BASE/api/v1/benefits | jq .

echo -e "\n=== 5. Benefits (skins only) ==="
curl -s "$BASE/api/v1/benefits?type=skin" | jq .

echo -e "\n=== 6. Profile ==="
curl -s $BASE/api/v1/profile \
  -H "Authorization: Bearer $TOKEN" | jq .

echo -e "\n=== 7. Purchase Arctic Fox ==="
curl -s -X POST "$BASE/api/v1/shop/purchase/[UUID70]" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Idempotency-Key: test-$(date +%s)" | jq .

echo -e "\n=== 8. Purchase free Early Bird ==="
curl -s -X POST "$BASE/api/v1/shop/purchase/[UUID71]" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Idempotency-Key: test2-$(date +%s)" | jq .

echo -e "\n=== 9. Inventory ==="
INVENTORY=$(curl -s $BASE/api/v1/inventory \
  -H "Authorization: Bearer $TOKEN")
echo $INVENTORY | jq .

ITEM_ID=$(echo $INVENTORY | jq -r '.data[0].id')
echo "First item ID: $ITEM_ID"

echo -e "\n=== 10. Equip ==="
curl -s -X PATCH "$BASE/api/v1/inventory/$ITEM_ID/equip" \
  -H "Authorization: Bearer $TOKEN" | jq .

echo -e "\n=== 11. Unequip ==="
curl -s -X PATCH "$BASE/api/v1/inventory/$ITEM_ID/unequip" \
  -H "Authorization: Bearer $TOKEN" | jq .

echo -e "\n=== 12. Transactions ==="
curl -s "$BASE/api/v1/transactions?limit=10" \
  -H "Authorization: Bearer $TOKEN" | jq .

echo -e "\n=== 13. Final profile ==="
curl -s $BASE/api/v1/profile \
  -H "Authorization: Bearer $TOKEN" | jq .

echo -e "\n=== 14. Unauthorized ==="
curl -s $BASE/api/v1/profile | jq .

echo -e "\n=== Done ==="