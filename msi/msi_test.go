// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package msi

import (
	"encoding/json"
	"github.com/Azure/azure-extension-foundation/httputil"
	"testing"
)

func TestGetMsi(t *testing.T) {
	m, err := GetMsi()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	t.Log(m.AccessToken)
	t.Log("Hello")
}

func TestSuccessfulGetMsi(t *testing.T) {
	const tokenValue = "token"
	httputil.Get = func(url string, headers map[string]string) (int, []byte, error) {
		m := Msi{AccessToken: tokenValue,
			ClientID:     "",
			ExpiresIn:    "",
			ExpiresOn:    "",
			ExtExpiresIn: "",
			NotBefore:    "",
			Resource:     "",
			TokenType:    ""}

		o, err := json.Marshal(m)
		if err != nil {
			return 0, nil, err
		}
		return 200, o, nil
	}

	msi, err := GetMsi()
	if err != nil {
		t.FailNow()
	}

	if msi.AccessToken != tokenValue {
		t.FailNow()
	}
}

func TestGetMsiReturns400(t *testing.T) {
	// metadata service will return 400 if MSI is disable
	httputil.Get = func(url string, headers map[string]string) (int, []byte, error) {
		return 400, nil, nil
	}

	_, err := GetMsi()
	if err == nil {
		t.FailNow()
	}
}
