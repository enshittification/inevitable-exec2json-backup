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
	}{
		{
			name:           "Echo command",
			args:           []string{"echo", "Hello, World!"},
			expectedStatus: 0,
			expectedStdout: "Hello, World!\n",
			expectedStderr: "",
		},
		{
			name:           "Non-existent command",
			args:           []string{"non_existent_command"},
			expectedStatus: 1,
			expectedStdout: "",
			expectedStderr: "exec: \"non_existent_command\": executable file not found in $PATH\n",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, exitCode := executeCommand(tc.args)

			if exitCode != tc.expectedStatus {
				t.Errorf("Expected status %d, got %d", tc.expectedStatus, exitCode)
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
		})
	}
}
