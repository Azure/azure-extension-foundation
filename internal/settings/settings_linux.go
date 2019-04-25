// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package settings

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os/exec"
	"path/filepath"
)

const settingsFileSuffix = ".settings"

type handlerSettingsFile struct {
	RuntimeSettings []struct {
		HandlerSettings handlerSettings `json:"handlerSettings"`
	} `json:"runtimeSettings"`
}

type handlerSettings struct {
	PublicSettings          map[string]interface{} `json:"publicSettings"`
	ProtectedSettingsBase64 string                 `json:"protectedSettings"`
	SettingsCertThumbprint  string                 `json:"protectedSettingsCertThumbprint"`
}

func GetExtensionSettings(sequenceNumber int, publicSettings, protectedSettings interface{}) error {
	publicSettingsJSON, protectedSettingsJSON, err := readSettings(sequenceNumber)
	if err != nil {
		return errorhelper.AddStackToError(fmt.Errorf("error reading handler settings: %v", err))
	}

	if err := unmarshalHandlerSettings(publicSettingsJSON, protectedSettingsJSON, &publicSettings, &protectedSettings); err != nil {
		return errorhelper.AddStackToError(fmt.Errorf("error parsing handler settings: %v", err))
	}

	return nil
}

// ReadSettings locates the .settings file and returns public settings
// JSON, and protected settings JSON (by decrypting it with the keys in
// configFolder).
func readSettings(sequenceNumber int) (public, protected map[string]interface{}, _ error) {
	hEnv, err := GetEnvironment()
	if err != nil {
		return nil, nil, errorhelper.AddStackToError(fmt.Errorf("unable to get handler environment: %v", err))
	}
	configFolderPath := hEnv.HandlerEnvironment.ConfigFolder

	cf, err := settingsFilePath(configFolderPath, sequenceNumber)
	if err != nil {
		return nil, nil, errorhelper.AddStackToError(fmt.Errorf("cannot locate settings file: %v", err))
	}
	hs, err := parseHandlerSettingsFile(cf)
	if err != nil {
		return nil, nil, errorhelper.AddStackToError(fmt.Errorf("error parsing settings file: %v", err))
	}

	public = hs.PublicSettings
	if err := unmarshalProtectedSettings(configFolderPath, hs, &protected); err != nil {
		return nil, nil, errorhelper.AddStackToError(fmt.Errorf("failed to parse protected settings: %v", err))
	}
	return public, protected, nil
}

// unmarshalHandlerSettings unmarshals given publicSettings/protectedSettings types
// assumed underlying values are JSON into references publicV/protectedV respectively
// (of struct types that contain structured fields for settings).
func unmarshalHandlerSettings(publicSettings, protectedSettings map[string]interface{}, publicV, protectedV interface{}) error {
	if err := unmarshalSettings(publicSettings, &publicV); err != nil {
		return errorhelper.AddStackToError(fmt.Errorf("failed to unmarshal public settings: %v", err))
	}
	if err := unmarshalSettings(protectedSettings, &protectedV); err != nil {
		return errorhelper.AddStackToError(fmt.Errorf("failed to unmarshal protected settings: %v", err))
	}
	return nil
}

// settingsFilePath returns the full path to the .settings file with the
// highest sequence number found in configFolder.
func settingsFilePath(configFolder string, sequenceNumber int) (string, error) {
	return filepath.Join(configFolder, fmt.Sprintf("%d%s", sequenceNumber, settingsFileSuffix)), nil
}

// parseHandlerSettings parses a handler settings file (e.g. 0.settings) and
// returns it as a structured object.
func parseHandlerSettingsFile(path string) (h handlerSettings, _ error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return h, errorhelper.AddStackToError(fmt.Errorf("error reading setting's file %s: %v", path, err))
	}
	if len(b) == 0 { // if no config is specified, we get an empty file
		return h, nil
	}

	var f handlerSettingsFile
	if err := json.Unmarshal(b, &f); err != nil {
		return h, errorhelper.AddStackToError(fmt.Errorf("error parsing json: %v", err))
	}
	if len(f.RuntimeSettings) != 1 {
		return h, errorhelper.AddStackToError(fmt.Errorf("wrong runtimeSettings count. expected:1, got:%d", len(f.RuntimeSettings)))
	}
	return f.RuntimeSettings[0].HandlerSettings, nil
}

// unmarshalSettings makes a round-trip JSON marshaling and unmarshaling
// from in (assumed map[interface]{}) to v (actual settings type).
func unmarshalSettings(in interface{}, v interface{}) error {
	s, err := json.Marshal(in)
	if err != nil {
		return errorhelper.AddStackToError(fmt.Errorf("failed to marshal into json: %v", err))
	}
	if err := json.Unmarshal(s, &v); err != nil {
		return errorhelper.AddStackToError(fmt.Errorf("failed to unmarshal json: %v", err))
	}

	return nil
}

// unmarshalProtectedSettings decodes the protected settings from handler
// runtime settings JSON file, decrypts it using the certificates and unmarshals
// into the given struct v.
func unmarshalProtectedSettings(configFolder string, hs handlerSettings, v interface{}) error {
	if hs.ProtectedSettingsBase64 == "" {
		return nil
	}
	if hs.SettingsCertThumbprint == "" {
		return errorhelper.AddStackToError(fmt.Errorf("handlerSettings has protected settings but no cert thumbprint"))
	}

	decoded, err := base64.StdEncoding.DecodeString(hs.ProtectedSettingsBase64)
	if err != nil {
		return errorhelper.AddStackToError(fmt.Errorf("failed to decode base64: %v", err))
	}

	// go two levels up where certs are placed (/var/lib/waagent)
	crt := filepath.Join(configFolder, "..", "..", fmt.Sprintf("%s.crt", hs.SettingsCertThumbprint))
	prv := filepath.Join(configFolder, "..", "..", fmt.Sprintf("%s.prv", hs.SettingsCertThumbprint))

	// we use os/exec instead of azure-docker-extension/pkg/executil here as
	// other extension handlers depend on this package for parsing handler
	// settings.
	cmd := exec.Command("openssl", "smime", "-inform", "DER", "-decrypt", "-recip", crt, "-inkey", prv)
	var bOut, bErr bytes.Buffer
	cmd.Stdin = bytes.NewReader(decoded)
	cmd.Stdout = &bOut
	cmd.Stderr = &bErr

	if err := cmd.Run(); err != nil {
		return errorhelper.AddStackToError(fmt.Errorf("decrypting protected settings failed: error=%v stderr=%s", err, string(bErr.Bytes())))
	}

	// decrypted: json object for protected settings
	if err := json.Unmarshal(bOut.Bytes(), &v); err != nil {
		return errorhelper.AddStackToError(fmt.Errorf("failed to unmarshal decrypted settings json: %v", err))
	}
	return nil
}
