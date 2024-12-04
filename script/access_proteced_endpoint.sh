#!/bin/bash

# 颜色定义
GREEN='\033[0;32m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

usage() {
    echo -e "${BLUE}Usage:${NC}"
    echo -e "  $0 [command] [token]"
    echo
    echo -e "${BLUE}Commands:${NC}"
    echo "  login                   - Get new access token with test credentials"
    echo "  protected <token>       - Access protected endpoint with given token"
    echo "  refresh <refresh_token> - Refresh access token with refresh token"
    echo
    echo -e "${BLUE}Examples:${NC}"
    echo "  $0 login"
    echo "  $0 protected eyJhbGciOiJIUzI1NiIs..."
    echo "  $0 refresh eyJhbGciOiJIUzI1NiIs..."
}

login() {
    echo -e "${BLUE}Logging in...${NC}"
    response=$(curl -s -X POST http://localhost:8080/auth/login \
        -H "Content-Type: application/json" \
        -d '{"username":"test","password":"test123"}')
    
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}Login successful!${NC}"
        echo -e "${BLUE}Response:${NC}"
        echo $response | jq '.'
        
        echo $response > .last_token.json
        echo -e "${GREEN}Token saved to .last_token.json${NC}"
    else
        echo -e "${RED}Login failed!${NC}"
    fi
}

access_protected() {
    local token=$1
    if [ -z "$token" ]; then
        echo -e "${RED}Error: Token is required${NC}"
        echo "Usage: $0 protected <token>"
        exit 1
    fi

    echo -e "${BLUE}Accessing protected endpoint...${NC}"
    response=$(curl -s -X GET http://localhost:8080/api/protected \
        -H "Authorization: Bearer $token")
    
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}Request successful!${NC}"
        echo -e "${BLUE}Response:${NC}"
        echo $response | jq '.'
    else
        echo -e "${RED}Request failed!${NC}"
    fi
}

refresh_token() {
    local refresh_token=$1
    if [ -z "$refresh_token" ]; then
        echo -e "${RED}Error: Refresh token is required${NC}"
        echo "Usage: $0 refresh <refresh_token>"
        exit 1
    fi

    echo -e "${BLUE}Refreshing token...${NC}"
    response=$(curl -s -X POST http://localhost:8080/auth/refresh \
        -H "Content-Type: application/json" \
        -d "{\"refresh_token\":\"$refresh_token\"}")
    
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}Token refresh successful!${NC}"
        echo -e "${BLUE}Response:${NC}"
        echo $response | jq '.'
        
        echo $response > .last_token.json
        echo -e "${GREEN}New token saved to .last_token.json${NC}"
    else
        echo -e "${RED}Token refresh failed!${NC}"
    fi
}

if ! command -v jq &> /dev/null; then
    echo -e "${RED}Error: jq is required but not installed.${NC}"
    echo "Please install jq first:"
    echo "  Ubuntu/Debian: sudo apt-get install jq"
    echo "  MacOS: brew install jq"
    exit 1
fi

case "$1" in
    "login")
        login
        ;;
    "protected")
        access_protected "$2"
        ;;
    "refresh")
        refresh_token "$2"
        ;;
    *)
        usage
        exit 1
        ;;
esac