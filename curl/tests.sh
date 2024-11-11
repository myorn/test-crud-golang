#!/bin/bash

# Base URL
BASE_URL="http://localhost:8080/equipment"

# Create Equipment
echo "Creating Equipment..."
CREATE_RESPONSE=$(curl -s -w "\nHTTP Status: %{http_code}\n" -X POST -H "Content-Type: application/json" -d '{
  "type": "TypeA",
  "status": "Active",
  "parameters": {"key1": "value1"}
}' $BASE_URL)
echo "Create Response: $CREATE_RESPONSE"

# Extract ID from Create Response
EQUIPMENT_ID=$(echo $CREATE_RESPONSE | jq -r '.id')
echo "Created Equipment ID: $EQUIPMENT_ID"

echo

# Get Equipment
echo "Getting Equipment..."
GET_RESPONSE=$(curl -s -w "\nHTTP Status: %{http_code}\n" -X GET $BASE_URL/$EQUIPMENT_ID)
echo "Get Response: $GET_RESPONSE"

echo

# Update Equipment
echo "Updating Equipment..."
UPDATE_RESPONSE=$(curl -s -w "\nHTTP Status: %{http_code}\n" -X PUT -H "Content-Type: application/json" -d '{
  "id": "'$EQUIPMENT_ID'",
  "type": "TypeB",
  "status": "Inactive",
  "parameters": {"key1": "value1", "key2": "value2"}
}' $BASE_URL/$EQUIPMENT_ID)
echo "Update Response: $UPDATE_RESPONSE"

echo

# Delete Equipment
echo "Deleting Equipment..."
DELETE_RESPONSE=$(curl -s -w "\nHTTP Status: %{http_code}\n" -X DELETE $BASE_URL/$EQUIPMENT_ID)
echo "Delete Response: $DELETE_RESPONSE"

echo

# Get Equipment after Deletion
echo "Getting Equipment after Deletion..."
GET_AFTER_DELETE_RESPONSE=$(curl -s -w "\nHTTP Status: %{http_code}\n" -X GET $BASE_URL/$EQUIPMENT_ID)
echo "Get after Delete Response: $GET_AFTER_DELETE_RESPONSE"

echo

# Restore Equipment
echo "Restoring Equipment..."
RESTORE_RESPONSE=$(curl -s -w "\nHTTP Status: %{http_code}\n" -X PATCH $BASE_URL/restore/$EQUIPMENT_ID)
echo "Restore Response: $RESTORE_RESPONSE"

echo

# Get Equipment after Restoration
echo "Getting Equipment after Restoration..."
GET_AFTER_RESTORE_RESPONSE=$(curl -s -w "\nHTTP Status: %{http_code}\n" -X GET $BASE_URL/$EQUIPMENT_ID)
echo "Get after Restore Response: $GET_AFTER_RESTORE_RESPONSE"