package main

import (
	"testing"
)

func TestExecuteCommand(t *testing.T) {
	testCases := []struct {
		name           string
		args           []string
		expectedStatus int
		expectedStdout string
		expectedStderr string
		expectError    bool
	}{
		{
			name:           "Echo command",
			args:           []string{"echo", "Hello, World!"},
			expectedStatus: 0,
			expectedStdout: "Hello, World!\n",
			expectedStderr: "",
			expectError:    false,
		},
		{
			name:           "Non-existent command",
			args:           []string{"non_existent_command"},
			expectedStatus: -1, // The status will not be set for a non-existent command
			expectedStdout: "",
			expectedStderr: "",
			expectError:    true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := executeCommand(tc.args)

			if tc.expectError && err == nil {
				t.Errorf("Expected an error, but got none")
			} else if !tc.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if result != nil {
				if status, ok := result["status"].(int); !ok || status != tc.expectedStatus {
					t.Errorf("Expected status %d, got %v", tc.expectedStatus, result["status"])
				}

				if stdout, ok := result["stdout"].(string); !ok || stdout != tc.expectedStdout {
					t.Errorf("Expected stdout %q, got %q", tc.expectedStdout, stdout)
				}

				if stderr, ok := result["stderr"].(string); !ok || stderr != tc.expectedStderr {
					t.Errorf("Expected stderr %q, got %q", tc.expectedStderr, stderr)
				}

				if _, ok := result["took"].(float64); !ok {
					t.Errorf("Expected 'took' to be a float64, got %T", result["took"])
				}

				if command, ok := result["command"].([]string); !ok || len(command) != len(tc.args) {
					t.Errorf("Expected command %v, got %v", tc.args, command)
				}
			} else if !tc.expectError {
				t.Errorf("Expected a result, but got nil")
			}
		})
	}
}
