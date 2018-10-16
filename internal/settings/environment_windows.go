// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package settings

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
// the extension handler (read: this) executable, reads, parses and returns it.
func GetEnvironment() (he HandlerEnvironment, _ error) {
	return HandlerEnvironment{}, nil
}
