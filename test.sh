#!/bin/bash

set -e

# Function to run a test case
run_test() {
    local test_name="$1"
    local command="$2"
    local expected_stdout="$3"
    local expected_stderr="$4"
    local expected_status="$5"

    echo "Running test: $test_name"
    output=$(./main $command | jq -r '.')
    
    actual_stdout=$(echo "$output" | jq -r '.stdout' | tr -d '\n')
    actual_stderr=$(echo "$output" | jq -r '.stderr' | tr -d '\n')
    actual_status=$(echo "$output" | jq -r '.status')
    actual_command=$(echo "$output" | jq -r '.command | join(" ")')

    if [ "$actual_stdout" != "$expected_stdout" ] || 
       [ "$actual_stderr" != "$expected_stderr" ] || 
       [ "$actual_status" != "$expected_status" ] ||
       [ "$actual_command" != "$command" ]; then
        echo "Test case '$test_name' failed"
        echo "Command: $command"
        echo "Expected stdout: $expected_stdout"
        echo "Actual stdout: $actual_stdout"
        echo "Expected stderr: $expected_stderr"
        echo "Actual stderr: $actual_stderr"
        echo "Expected status: $expected_status"
        echo "Actual status: $actual_status"
        echo "Actual command: $actual_command"
        exit 1
    fi

    echo "Test case '$test_name' passed"
    echo
}

# Test cases
run_test "Echo command" "echo 'Hello, World!'" "Hello, World!" "" "0"
run_test "ls command" "ls -l" ".*" "" "0"
run_test "Non-existent command" "nonexistentcommand" "" ".*: command not found" "127"
run_test "Exit with non-zero status" "bash -c 'exit 42'" "" "" "42"
run_test "Command with arguments" "echo arg1 arg2 arg3" "arg1 arg2 arg3" "" "0"
run_test "Command with quoted arguments" "echo 'arg with spaces'" "arg with spaces" "" "0"
run_test "Command with environment variables" "bash -c 'echo \$PATH'" ".*" "" "0"
run_test "Command with stdin input" "cat" "test input" "" "0" < <(echo "test input")
run_test "Command with stderr output" "bash -c 'echo error >&2; echo output'" "output" "error" "0"
run_test "Long-running command" "sleep 2" "" "" "0"

echo "All tests passed!"