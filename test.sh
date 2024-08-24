#!/bin/bash

set -e

# Check if main executable exists
if [ ! -f "./main" ]; then
    echo "Error: main executable not found. Make sure to build it before running tests."
    exit 1
fi

# Function to run a test case
run_test() {
    local test_name="$1"
    local command="$2"
    local expected_stdout="$3"
    local expected_stderr="$4"
    local expected_status="$5"

    echo "Running test: $test_name"
    output=$(./main $command | jq -r '.')
    
    # ... rest of the function remains the same
}

# ... rest of the test cases remain the same
