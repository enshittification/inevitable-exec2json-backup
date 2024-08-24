#!/bin/bash

set -e

# Check if exec2json executable exists
if [ ! -f "./exec2json" ]; then
    echo "Error: exec2json executable not found. Make sure to build it before running tests."
    exit 1
fi

# Function to run a test case
run_test() {
    local test_name="$1"
    local command="$2"
    local expected_stdout_pattern="$3"
    local expected_stderr_pattern="$4"
    local expected_status="$5"

    echo "Running test: $test_name"
    output=$(./exec2json $command | jq -r '.')
    
    actual_stdout=$(echo "$output" | jq -r '.stdout')
    actual_stderr=$(echo "$output" | jq -r '.stderr')
    actual_status=$(echo "$output" | jq -r '.status')
    actual_command=$(echo "$output" | jq -r '.command | join(" ")')

    # Remove surrounding quotes if present
    actual_stdout="${actual_stdout#\'}"
    actual_stdout="${actual_stdout%\'}"

    if [[ ! "$actual_stdout" =~ $expected_stdout_pattern ]] || 
       [[ ! "$actual_stderr" =~ $expected_stderr_pattern ]] || 
       [ "$actual_status" != "$expected_status" ] ||
       [ "$actual_command" != "$command" ]; then
        echo "Test case '$test_name' failed"
        echo "Command: $command"
        echo "Expected stdout pattern: $expected_stdout_pattern"
        echo "Actual stdout:"
        echo "$actual_stdout"
        echo "Expected stderr pattern: $expected_stderr_pattern"
        echo "Actual stderr:"
        echo "$actual_stderr"
        echo "Expected status: $expected_status"
        echo "Actual status: $actual_status"
        echo "Actual command: $actual_command"
        exit 1
    fi

    echo "Test case '$test_name' passed"
    echo
}

# Test cases
run_test "Echo command" "echo 'Hello, World!'" "^Hello, World!$" "^$" "0"
run_test "ls command" "ls -l" "^total [0-9]+\n-.*" "^$" "0"
run_test "Non-existent command" "nonexistentcommand" "^$" ".*: command not found" "127"
run_test "Exit with non-zero status" "bash -c 'exit 42'" "^$" "^$" "42"
run_test "Command with arguments" "echo arg1 arg2 arg3" "^arg1 arg2 arg3$" "^$" "0"
run_test "Command with quoted arguments" "echo 'arg with spaces'" "^arg with spaces$" "^$" "0"
run_test "Command with environment variables" "bash -c 'echo \$PATH'" ".+" "^$" "0"
run_test "Command with stdin input" "cat" "^test input$" "^$" "0" < <(echo "test input")
run_test "Command with stderr output" "bash -c 'echo error >&2; echo output'" "^output$" "^error$" "0"
run_test "Long-running command" "sleep 2" "^$" "^$" "0"

echo "All tests passed!"
