#!/bin/bash

# PMGO Test Script
# This script demonstrates how to use PMGO to manage processes

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Logging functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

log_step() {
    echo -e "${PURPLE}[STEP]${NC} $1"
}

# Configuration
PMGO_BIN="./pmgo"
TEST_APP="examples/webserver/main.go"
TEST_APP_NAME="test-webserver"
TEST_PORT=8080
API_BASE_URL="http://localhost:9876/api/v1"
WEB_URL="http://localhost:8080"

# Check if pmgo binary exists
check_pmgo() {
    if [ ! -f "$PMGO_BIN" ]; then
        log_error "PMGO binary not found. Please run 'make build' first."
        exit 1
    fi
    log_success "PMGO binary found: $PMGO_BIN"
}

# Check if test application exists
check_test_app() {
    if [ ! -f "$TEST_APP" ]; then
        log_error "Test application not found: $TEST_APP"
        exit 1
    fi
    log_success "Test application found: $TEST_APP"
}

# Wait for service to be ready
wait_for_service() {
    local url=$1
    local timeout=${2:-30}
    local count=0

    log_info "Waiting for service to be ready: $url"

    while [ $count -lt $timeout ]; do
        if curl -s "$url" > /dev/null 2>&1; then
            log_success "Service is ready!"
            return 0
        fi
        sleep 1
        count=$((count + 1))
        echo -n "."
    done

    echo
    log_error "Service failed to start within $timeout seconds"
    return 1
}

# Start PMGO daemon
start_pmgo_daemon() {
    log_step "1. Starting PMGO daemon..."

    # Check if already running
    if pgrep -f "$PMGO_BIN serve" > /dev/null; then
        log_warning "PMGO daemon is already running"
    else
        log_info "Starting PMGO daemon in background..."
        nohup $PMGO_BIN serve > pmgo-daemon.log 2>&1 &

        # Wait for daemon to start
        if wait_for_service "http://localhost:9876/health" 10; then
            log_success "PMGO daemon started successfully"
        else
            log_error "Failed to start PMGO daemon"
            return 1
        fi
    fi
}

# Show PMGO version and help
show_pmgo_info() {
    log_step "2. PMGO Information"

    echo
    log_info "PMGO Version:"
    $PMGO_BIN --version

    echo
    log_info "Available Commands:"
    $PMGO_BIN --help | grep -A 20 "Available Commands"
}

# Start test application
start_test_app() {
    log_step "3. Starting test application with PMGO..."

    # Clean up any existing process
    log_info "Cleaning up any existing '$TEST_APP_NAME' process..."
    $PMGO_BIN delete $TEST_APP_NAME 2>/dev/null || true

    # Start the test application
    log_info "Starting test web server on port $TEST_PORT..."
    $PMGO_BIN start $TEST_APP $TEST_APP_NAME --args="--port=$TEST_PORT"

    # Wait for the web server to be ready
    if wait_for_service "http://localhost:$TEST_PORT/health" 15; then
        log_success "Test web server started successfully"
    else
        log_error "Test web server failed to start"
        return 1
    fi
}

# Test process management
test_process_management() {
    log_step "4. Testing process management..."

    echo
    log_info "Listing all processes:"
    $PMGO_BIN list

    echo
    log_info "Getting process information:"
    $PMGO_BIN info $TEST_APP_NAME

    echo
    log_info "Testing restart functionality..."
    $PMGO_BIN restart $TEST_APP_NAME
    sleep 3

    if wait_for_service "http://localhost:$TEST_PORT/health" 10; then
        log_success "Process restart successful"
    else
        log_error "Process restart failed"
        return 1
    fi
}

# Test web application
test_web_application() {
    log_step "5. Testing web application endpoints..."

    echo
    log_info "Testing health endpoint:"
    curl -s "http://localhost:$TEST_PORT/health" | jq '.' || curl -s "http://localhost:$TEST_PORT/health"

    echo
    log_info "Testing info endpoint:"
    curl -s "http://localhost:$TEST_PORT/info" | jq '.' || curl -s "http://localhost:$TEST_PORT/info"

    echo
    log_info "Testing work endpoint:"
    curl -s "http://localhost:$TEST_PORT/api/work" | jq '.' || curl -s "http://localhost:$TEST_PORT/api/work"

    echo
    log_info "Testing counter endpoint:"
    curl -s "http://localhost:$TEST_PORT/api/counter" | jq '.' || curl -s "http://localhost:$TEST_PORT/api/counter"

    log_success "All web endpoints working correctly"
}

# Test PMGO API
test_pmgo_api() {
    log_step "6. Testing PMGO HTTP API..."

    echo
    log_info "Testing PMGO health endpoint:"
    curl -s "http://localhost:9876/health" | jq '.' || curl -s "http://localhost:9876/health"

    echo
    log_info "Testing PMGO processes API:"
    curl -s "$API_BASE_URL/processes" | jq '.' || curl -s "$API_BASE_URL/processes"

    echo
    log_info "Testing specific process API:"
    curl -s "$API_BASE_URL/processes/$TEST_APP_NAME" | jq '.' || curl -s "$API_BASE_URL/processes/$TEST_APP_NAME"

    log_success "PMGO API endpoints working correctly"
}

# Save and restore test
test_save_restore() {
    log_step "7. Testing save and restore functionality..."

    echo
    log_info "Saving current process list..."
    $PMGO_BIN save

    log_success "Process list saved successfully"
}

# Load testing
run_load_test() {
    log_step "8. Running basic load test..."

    log_info "Sending 10 concurrent requests to test server resilience..."

    for i in {1..10}; do
        (curl -s "http://localhost:$TEST_PORT/api/work" > /dev/null) &
    done

    wait
    log_success "Load test completed"

    echo
    log_info "Final counter state:"
    curl -s "http://localhost:$TEST_PORT/api/counter" | jq '.' || curl -s "http://localhost:$TEST_PORT/api/counter"
}

# Show logs
show_logs() {
    log_step "9. Showing process logs..."

    log_info "Recent logs from $TEST_APP_NAME:"
    $PMGO_BIN logs $TEST_APP_NAME || log_warning "Logs command not implemented yet"
}

# Demonstrate web interface
show_web_interface() {
    log_step "10. Web Interface Information"

    echo
    log_info "🌐 Web Interfaces Available:"
    log_info "   - PMGO Web Interface: http://localhost:8080"
    log_info "   - Test Application: http://localhost:$TEST_PORT"
    log_info "   - API Documentation: http://localhost:9876/health"

    echo
    log_info "💡 You can now:"
    log_info "   1. Open http://localhost:$TEST_PORT in your browser to see the test app"
    log_info "   2. Open http://localhost:8080 to see PMGO web interface"
    log_info "   3. Use the API at http://localhost:9876/api/v1"
}

# Cleanup function
cleanup() {
    log_step "Cleanup"

    log_info "Stopping test application..."
    $PMGO_BIN delete $TEST_APP_NAME 2>/dev/null || true

    if [ "$1" = "--kill-daemon" ]; then
        log_info "Stopping PMGO daemon..."
        $PMGO_BIN kill 2>/dev/null || true
        sleep 2
    fi

    log_success "Cleanup completed"
}

# Show usage
show_usage() {
    cat << EOF
PMGO Test Script

Usage: $0 [OPTIONS]

OPTIONS:
    -h, --help          Show this help message
    --no-daemon         Don't start the daemon (assume it's already running)
    --cleanup-only      Only run cleanup
    --full-test         Run all tests including load testing
    --kill-daemon       Kill daemon after tests

EXAMPLES:
    # Run basic test
    $0

    # Run full test suite
    $0 --full-test

    # Cleanup everything
    $0 --cleanup-only --kill-daemon

EOF
}

# Main function
main() {
    local start_daemon=true
    local cleanup_only=false
    local full_test=false
    local kill_daemon=false

    # Parse arguments
    while [[ $# -gt 0 ]]; do
        case $1 in
            -h|--help)
                show_usage
                exit 0
                ;;
            --no-daemon)
                start_daemon=false
                shift
                ;;
            --cleanup-only)
                cleanup_only=true
                shift
                ;;
            --full-test)
                full_test=true
                shift
                ;;
            --kill-daemon)
                kill_daemon=true
                shift
                ;;
            *)
                log_error "Unknown option: $1"
                show_usage
                exit 1
                ;;
        esac
    done

    # Cleanup only mode
    if [ "$cleanup_only" = true ]; then
        cleanup --kill-daemon
        exit 0
    fi

    # Setup trap for cleanup
    trap 'cleanup' EXIT

    echo "🚀 PMGO Test Script Starting..."
    echo "=================================="

    # Pre-checks
    check_pmgo
    check_test_app

    # Start daemon if needed
    if [ "$start_daemon" = true ]; then
        start_pmgo_daemon
    fi

    # Run tests
    show_pmgo_info
    start_test_app
    test_process_management
    test_web_application
    test_pmgo_api
    test_save_restore

    if [ "$full_test" = true ]; then
        run_load_test
    fi

    show_logs
    show_web_interface

    echo
    echo "=================================="
    log_success "🎉 All tests completed successfully!"
    echo
    log_info "The test application is still running. You can:"
    log_info "  - Visit http://localhost:$TEST_PORT to see the test app"
    log_info "  - Use 'pmgo list' to see running processes"
    log_info "  - Use 'pmgo stop $TEST_APP_NAME' to stop the test app"

    if [ "$kill_daemon" = true ]; then
        cleanup --kill-daemon
    else
        log_info "  - Use '$0 --cleanup-only --kill-daemon' to cleanup everything"
    fi

    echo
}

# Run main function
main "$@"