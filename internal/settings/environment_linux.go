// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package settings

import (
	"encoding/json"
	"fmt"
	"github.com/Azure/azure-extension-foundation/errorhelper"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// HandlerEnvFileName is the file name of the Handler Environment as placed by the
// Azure Linux Guest Agent.
const HandlerEnvFileName = "HandlerEnvironment.json"

// HandlerEnvironment describes the handler environment configuration presented
// to the extension handler by the Azure Linux Guest Agent.
type HandlerEnvironment struct {
	Version            float64 `json:"version"`
	Name               string  `json:"name"`
	HandlerEnvironment struct {
		HeartbeatFile string `json:"heartbeatFile"`
		StatusFolder  string `json:"statusFolder"`
		ConfigFolder  string `json:"configFolder"`
		LogFolder     string `json:"logFolder"`
	}
}

// GetEnvironment locates the HandlerEnvironment.json file by assuming it lives next to or one level above
// the extension handler (read: this) executable, reads, parses and returns it
func GetEnvironment() (environment HandlerEnvironment, _ error) {
	dir, err := scriptDirectory()
	if err != nil {
		return environment, errorhelper.AddStackToError(fmt.Errorf("vmextension: cannot find base directory of the running process: %v", err))
	}
	paths := []string{
		filepath.Join(dir, HandlerEnvFileName),       // this level (i.e. executable is in [EXT_NAME]/.)
		filepath.Join(dir, "..", HandlerEnvFileName), // one up (i.e. executable is in [EXT_NAME]/bin/.)
	}
	var b []byte
	for _, p := range paths {
		o, err := ioutil.ReadFile(p)
		if err != nil && !os.IsNotExist(err) {
			return environment, errorhelper.AddStackToError(fmt.Errorf("vmextension: error examining HandlerEnvironment at '%s': %v", p, err))
		} else if err == nil {
			b = o
			break
		}
	}
	if b == nil {
		return environment, errorhelper.AddStackToError(fmt.Errorf("vmextension: Cannot find HandlerEnvironment at paths: %s", strings.Join(paths, ", ")))
	}
	return parseEnvironmentManifest(b)
}

// scriptDirectory returns the absolute path of the running process
func scriptDirectory() (string, error) {
	p, err := filepath.Abs(os.Args[0])
	if err != nil {
		return "", err
	}
	return filepath.Dir(p), nil
}

// parseEnvironmentManifest parses the /var/lib/waagent/[extension]/HandlerEnvironment.json format
func parseEnvironmentManifest(b []byte) (he HandlerEnvironment, _ error) {
	var hf []HandlerEnvironment

	if err := json.Unmarshal(b, &hf); err != nil {
		return he, errorhelper.AddStackToError(fmt.Errorf("vmextension: failed to parse handler env: %v", err))
	}
	if len(hf) != 1 {
		return he, errorhelper.AddStackToError(fmt.Errorf("vmextension: expected 1 config in parsed HandlerEnvironment, found: %v", len(hf)))
	}
	return hf[0], nil
}
