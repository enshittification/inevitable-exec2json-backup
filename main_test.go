package main

import (
	"encoding/json"
	"os"
	"testing"
)

func TestCommandExecutor(t *testing.T) {
	// Save the original os.Args and restore it after the test
	originalArgs := os.Args
	defer func() { os.Args = originalArgs }()

	// Set up test cases
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
			// Set up os.Args for the test
			os.Args = append([]string{"command-executor"}, tc.args...)

			// Capture stdout
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			// Run the main function
			main()

			// Restore stdout
			w.Close()
			os.Stdout = oldStdout

			// Read the captured output
			var output map[string]interface{}
			json.NewDecoder(r).Decode(&output)

			// Check the results
			if status, ok := output["status"].(float64); !ok || int(status) != tc.expectedStatus {
				t.Errorf("Expected status %d, got %v", tc.expectedStatus, output["status"])
			}

			if stdout, ok := output["stdout"].(string); !ok || stdout != tc.expectedStdout {
				t.Errorf("Expected stdout %q, got %q", tc.expectedStdout, stdout)
			}

			if stderr, ok := output["stderr"].(string); !ok || stderr != tc.expectedStderr {
				t.Errorf("Expected stderr %q, got %q", tc.expectedStderr, stderr)
			}

			if _, ok := output["took"].(float64); !ok {
				t.Errorf("Expected 'took' to be a float64, got %T", output["took"])
			}

			if command, ok := output["command"].([]interface{}); !ok || len(command) != len(tc.args) {
				t.Errorf("Expected command %v, got %v", tc.args, command)
			}
		})
	}
}
