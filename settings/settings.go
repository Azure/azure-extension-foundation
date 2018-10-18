// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package settings

import "github.com/Azure/azure-extension-foundation/internal/settings"

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

// GetExtensionSettings reads the settings for the provided sequenceNumber and assigns the settings to the
// respective structure reference
func GetExtensionSettings(sequenceNumber int, publicSettings, protectedSettings interface{}) error {
	return settings.GetExtensionSettings(sequenceNumber, publicSettings, protectedSettings)
}

// GetHandlerEnvironment returns the handler environment properties
func GetHandlerEnvironment() (HandlerEnvironment, error) {
	// temporary work around since type alias is avail in 1.9 and build box only support 1.8
	he, err := settings.GetEnvironment()
	return HandlerEnvironment(he), err
}
