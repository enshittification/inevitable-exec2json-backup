/**
 * Copyright 2024 Automattic, Inc.
 *
 * This program is free software; you can redistribute it and/or
 * modify it under the terms of the GNU General Public License
 * as published by the Free Software Foundation; either version 2
 * of the License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; If not, see <http://www.gnu.org/licenses/>.
 */

package main

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"os"
	"os/exec"
	"time"
)

func executeCommand(args []string) (map[string]interface{}, int) {
	if len(args) < 1 {
		log.Fatal("You must supply a command to execute")
	}
	command := exec.Command(args[0], args[1:]...)
	stdin, err := command.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}
	stderr, err := command.StderrPipe()
	if err != nil {
		log.Fatal(err)
	}
	stdout, err := command.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	startTime := time.Now()
	go func() {
		defer stdin.Close()
		io.Copy(stdin, os.Stdin)
	}()
	err = command.Start()
	if err != nil {
		log.Fatal(err)
	}
	outString, err := io.ReadAll(stdout)
	if err != nil {
		log.Fatal(err)
	}
	errString, err := io.ReadAll(stderr)
	if err != nil {
		log.Fatal(err)
	}
	if err := command.Wait(); err != nil {
		var exitError *exec.ExitError
		if !errors.As(err, &exitError) {
			log.Fatal(err)
		}
	}
	took := time.Now().Sub(startTime).Seconds()
	result := map[string]interface{}{
		"command": args,
		"stdout":  string(outString),
		"stderr":  string(errString),
		"status":  command.ProcessState.ExitCode(),
		"took":    took,
	}
	return result, command.ProcessState.ExitCode()
}

func main() {
	result, exitCode := executeCommand(os.Args[1:])
	json.NewEncoder(os.Stdout).Encode(result)
	os.Exit(exitCode)
}
