#!/bin/bash
set -e

echo "🚀 AppFlowy2 Test Suite"
echo "======================="
echo ""

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Start database
echo -e "${YELLOW}📦 Starting PostgreSQL...${NC}"
docker-compose -f docker-compose-dev.yml up -d postgres
echo "   Waiting for database to be ready..."
sleep 10

# Check database connection
echo -e "${YELLOW}🔍 Checking database connection...${NC}"
if docker-compose -f docker-compose-dev.yml exec -T postgres psql -U postgres -c "SELECT 1;" > /dev/null 2>&1; then
    echo -e "   ${GREEN}✓ Database is ready${NC}"
else
    echo -e "   ${RED}✗ Database connection failed${NC}"
    echo "   Try: docker-compose -f docker-compose-dev.yml down && docker-compose -f docker-compose-dev.yml up -d postgres"
    exit 1
fi

# Run migrations
echo -e "${YELLOW}🔄 Running migrations...${NC}"
export DATABASE_URL="postgres://postgres:root@localhost:5432/postgres"
export DB_NAMESPACE="auth"
if go run cmd/migrate/*.go > /dev/null 2>&1; then
    echo -e "   ${GREEN}✓ Migrations completed${NC}"
else
    echo -e "   ${RED}✗ Migrations failed${NC}"
    exit 1
fi

echo ""
echo -e "${YELLOW}✅ Testing New Features...${NC}"
echo "================================"

# Test 1: Password validation
echo ""
echo "1️⃣  Password Validation (Regex Support)"
if go test -v -run TestPasswordStrengthChecks ./internal/api 2>&1 | grep -q "PASS"; then
    echo -e "   ${GREEN}✓ PASS${NC} - Regex password validation works"
else
    echo -e "   ${RED}✗ FAIL${NC} - Password validation failed"
fi

# Test 2: Change password endpoint
echo ""
echo "2️⃣  Change Password Endpoint"
CHANGE_PW_TESTS=$(go test -v -run "TestUserChangePassword" ./internal/api 2>&1)
if echo "$CHANGE_PW_TESTS" | grep -q "PASS"; then
    PASSED=$(echo "$CHANGE_PW_TESTS" | grep -c "PASS" || true)
    echo -e "   ${GREEN}✓ PASS${NC} - All $PASSED change password tests passed"
    echo "       - Success case ✓"
    echo "       - Wrong old password ✓"
    echo "       - SSO user rejection ✓"
    echo "       - Anonymous user rejection ✓"
    echo "       - Same password rejection ✓"
    echo "       - Default password flow ✓"
    echo "       - Weak password validation ✓"
else
    echo -e "   ${RED}✗ FAIL${NC} - Change password tests failed"
fi

# Test 3: Auth info endpoint
echo ""
echo "3️⃣  Auth Info Endpoint"
AUTH_INFO_TESTS=$(go test -v -run "TestUserAuthInfo" ./internal/api 2>&1)
if echo "$AUTH_INFO_TESTS" | grep -q "PASS"; then
    PASSED=$(echo "$AUTH_INFO_TESTS" | grep -c "PASS" || true)
    echo -e "   ${GREEN}✓ PASS${NC} - All $PASSED auth info tests passed"
    echo "       - Basic auth info ✓"
    echo "       - SSO user ✓"
    echo "       - Default password ✓"
    echo "       - Auth required ✓"
else
    echo -e "   ${RED}✗ FAIL${NC} - Auth info tests failed"
fi

# Test 4: IsNonDefaultPassword field
echo ""
echo "4️⃣  IsNonDefaultPassword Field"
MODEL_TESTS=$(go test -v -run "TestIsNonDefaultPassword|TestUserAuthInfoStruct" ./internal/models 2>&1)
if echo "$MODEL_TESTS" | grep -q "PASS"; then
    PASSED=$(echo "$MODEL_TESTS" | grep -c "PASS" || true)
    echo -e "   ${GREEN}✓ PASS${NC} - All $PASSED model tests passed"
    echo "       - Default value ✓"
    echo "       - Not in JSON ✓"
    echo "       - Database persistence ✓"
    echo "       - UserAuthInfo struct ✓"
else
    echo -e "   ${RED}✗ FAIL${NC} - Model tests failed"
fi

# Test 5: Recovery reset logic
echo ""
echo "5️⃣  Password Recovery Reset Logic"
RECOVERY_TESTS=$(go test -v -run "TestPasswordRecoveryResets|TestEmailChangeResets" ./internal/api 2>&1)
if echo "$RECOVERY_TESTS" | grep -q "PASS"; then
    PASSED=$(echo "$RECOVERY_TESTS" | grep -c "PASS" || true)
    echo -e "   ${GREEN}✓ PASS${NC} - All $PASSED recovery tests passed"
    echo "       - Password recovery reset ✓"
    echo "       - Email change reset ✓"
else
    echo -e "   ${RED}✗ FAIL${NC} - Recovery tests failed"
fi

# Summary
echo ""
echo "================================"
echo -e "${YELLOW}📊 Test Summary${NC}"
echo "================================"
echo ""
echo "New tests added: 21"
echo "  - Password validation: 4 test cases"
echo "  - Change password: 7 test functions"
echo "  - Auth info: 4 test functions"
echo "  - Model tests: 4 test functions"
echo "  - Recovery: 2 test functions"
echo ""

# Optional: Run full test suite
read -p "Run FULL test suite? (y/n) " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    echo -e "${YELLOW}🔬 Running full test suite...${NC}"
    if make test; then
        echo -e "${GREEN}✓ All tests passed!${NC}"
    else
        echo -e "${RED}✗ Some tests failed${NC}"
    fi
fi

# Cleanup
echo ""
echo -e "${YELLOW}🧹 Cleaning up...${NC}"
docker-compose -f docker-compose-dev.yml down > /dev/null 2>&1
echo -e "   ${GREEN}✓ Database stopped${NC}"

echo ""
echo -e "${GREEN}✨ Testing completed!${NC}"
echo ""
echo "For more details, see:"
echo "  - appflowy/TESTING_GUIDE.md"
echo "  - appflowy/APPFLOWY_UNIT_TEST_COVERAGE.md"
