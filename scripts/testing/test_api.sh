#!/bin/bash

# API Testing Script for Reolink Server
# This script tests all major API endpoints

BASE_URL="http://localhost:8080"
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo "========================================="
echo "Reolink Server API Testing"
echo "========================================="
echo ""

# Test 1: Health Check
echo -e "${YELLOW}Test 1: Health Check${NC}"
RESPONSE=$(curl -s ${BASE_URL}/health)
if echo "$RESPONSE" | jq -e '.success == true' > /dev/null 2>&1; then
    echo -e "${GREEN}✓ Health check passed${NC}"
    echo "$RESPONSE" | jq .
else
    echo -e "${RED}✗ Health check failed${NC}"
    echo "$RESPONSE"
fi
echo ""

# Test 2: Readiness Check
echo -e "${YELLOW}Test 2: Readiness Check${NC}"
RESPONSE=$(curl -s ${BASE_URL}/ready)
if echo "$RESPONSE" | jq -e '.success == true' > /dev/null 2>&1; then
    echo -e "${GREEN}✓ Readiness check passed${NC}"
    echo "$RESPONSE" | jq .
else
    echo -e "${RED}✗ Readiness check failed${NC}"
    echo "$RESPONSE"
fi
echo ""

# Test 3: Login
echo -e "${YELLOW}Test 3: Login (admin/admin)${NC}"
LOGIN_RESPONSE=$(curl -s -X POST ${BASE_URL}/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin"}')

if echo "$LOGIN_RESPONSE" | jq -e '.success == true' > /dev/null 2>&1; then
    echo -e "${GREEN}✓ Login successful${NC}"
    TOKEN=$(echo "$LOGIN_RESPONSE" | jq -r '.data.token')
    echo "Token: ${TOKEN:0:50}..."
    echo "$LOGIN_RESPONSE" | jq '.data.user'
else
    echo -e "${RED}✗ Login failed${NC}"
    echo "$LOGIN_RESPONSE" | jq .
    exit 1
fi
echo ""

# Test 4: List Cameras (Protected)
echo -e "${YELLOW}Test 4: List Cameras (Protected Endpoint)${NC}"
RESPONSE=$(curl -s ${BASE_URL}/api/v1/cameras \
  -H "Authorization: Bearer $TOKEN")
if echo "$RESPONSE" | jq -e '.success == true' > /dev/null 2>&1; then
    echo -e "${GREEN}✓ List cameras successful${NC}"
    echo "$RESPONSE" | jq .
else
    echo -e "${RED}✗ List cameras failed${NC}"
    echo "$RESPONSE" | jq .
fi
echo ""

# Test 5: List Events (Protected)
echo -e "${YELLOW}Test 5: List Events (Protected Endpoint)${NC}"
RESPONSE=$(curl -s ${BASE_URL}/api/v1/events \
  -H "Authorization: Bearer $TOKEN")
if echo "$RESPONSE" | jq -e '.success == true' > /dev/null 2>&1; then
    echo -e "${GREEN}✓ List events successful${NC}"
    echo "$RESPONSE" | jq .
else
    echo -e "${RED}✗ List events failed${NC}"
    echo "$RESPONSE" | jq .
fi
echo ""

# Test 6: Unauthorized Access (No Token)
echo -e "${YELLOW}Test 6: Unauthorized Access (No Token)${NC}"
RESPONSE=$(curl -s ${BASE_URL}/api/v1/cameras)
if echo "$RESPONSE" | jq -e '.success == false' > /dev/null 2>&1; then
    echo -e "${GREEN}✓ Correctly rejected unauthorized request${NC}"
    echo "$RESPONSE" | jq .
else
    echo -e "${RED}✗ Should have rejected unauthorized request${NC}"
    echo "$RESPONSE" | jq .
fi
echo ""

# Test 7: Invalid Login
echo -e "${YELLOW}Test 7: Invalid Login (Wrong Password)${NC}"
RESPONSE=$(curl -s -X POST ${BASE_URL}/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"wrongpassword"}')
if echo "$RESPONSE" | jq -e '.success == false' > /dev/null 2>&1; then
    echo -e "${GREEN}✓ Correctly rejected invalid credentials${NC}"
    echo "$RESPONSE" | jq .
else
    echo -e "${RED}✗ Should have rejected invalid credentials${NC}"
    echo "$RESPONSE" | jq .
fi
echo ""

echo "========================================="
echo "API Testing Complete!"
echo "========================================="

