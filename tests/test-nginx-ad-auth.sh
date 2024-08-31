#!/bin/bash

# Configuration
AUTH_SERVER="http://localhost:8080"
AUTH_ENDPOINT="/auth"
VALIDUSER=${VALIDUSER:-"validuser"}
CORRECTPASSWORD=${CORRECTPASSWORD:-"correctpassword"}
# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Function to make a curl request and check the response
test_auth() {
    local test_name="$1"
    local user="$2"
    local pass="$3"
    local protocol="$4"
    local expected_status="$5"

    echo "Running test: $test_name"
    
    response=$(curl -s -D - -X GET \
        -H "Auth-User: $user" \
        -H "Auth-Pass: $pass" \
        -H "Auth-Protocol: $protocol" \
        -H "Auth-Login-Attempt: 1" \
        -H "Client-IP: 192.0.2.42" \
        -H "Client-Host: client.example.org" \
        "${AUTH_SERVER}${AUTH_ENDPOINT}")

    echo "Response: $response"
    if echo "$response" | grep -q "Auth-Status: $expected_status"; then
        echo -e "${GREEN}Test passed${NC}"
    else
        echo -e "${RED}Test failed${NC}"
        echo "Expected: Auth-Status: $expected_status"
        echo "Received: $response"
    fi
    echo
}

# Test cases
test_auth "Valid IMAP login" "$VALIDUSER" "$CORRECTPASSWORD" "imap" "OK"
test_auth "Invalid IMAP login" "in$VALIDUSER" "wrongpassword" "imap" "Invalid login or password"
test_auth "Valid POP3 login" "$VALIDUSER" "$CORRECTPASSWORD" "pop3" "OK"
test_auth "Valid SMTP login" "$VALIDUSER" "$CORRECTPASSWORD" "smtp" "OK"
test_auth "Missing username" "" "somepassword" "imap" "No login or password"
test_auth "Missing password" "someuser" "" "imap" "No login or password"

# Test for checking returned Auth-Server and Auth-Port
echo "Checking Auth-Server and Auth-Port for valid login"
response=$(curl -s -D - -X GET \
    -H "Auth-User: $VALIDUSER" \
    -H "Auth-Pass: $CORRECTPASSWORD" \
    -H "Auth-Protocol: imap" \
    -H "Auth-Login-Attempt: 1" \
    -H "Client-IP: 192.0.2.42" \
    -H "Client-Host: client.example.org" \
    "${AUTH_SERVER}${AUTH_ENDPOINT}")

if echo "$response" | grep -q "Auth-Status: OK" && \
   echo "$response" | grep -q "Auth-Server:" && \
   echo "$response" | grep -q "Auth-Port:"; then
    echo -e "${GREEN}Test passed${NC}"
    echo "$response"
else
    echo -e "${RED}Test failed${NC}"
    echo "Expected Auth-Status: OK, Auth-Server, and Auth-Port in response"
    echo "Received: $response"
fi