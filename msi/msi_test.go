// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package msi

import (
	"encoding/json"
	"testing"
)

type httpClientMock struct {
	get_f func(url string, headers map[string]string) (responseCode int, body []byte, err error)
}

func (h httpClientMock) Get(url string, headers map[string]string) (responseCode int, body []byte, err error) {
	return h.get_f(url, headers)
}

func TestSuccessfulGetMsi(t *testing.T) {
	const tokenValue = "token"
	httpClient := httpClientMock{get_f: func(url string, headers map[string]string) (responseCode int, body []byte, err error) {
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
	}}
	provider := NewMsiProvider(&httpClient)

	msi, err := provider.GetMsi()
	if err != nil {
		t.FailNow()
	}

	if msi.AccessToken != tokenValue {
		t.FailNow()
	}
}

func TestGetMsiReturns400(t *testing.T) {
	// metadata service will return 400 if MSI is disable
	httpClient := httpClientMock{get_f: func(url string, headers map[string]string) (responseCode int, body []byte, err error) {
		return 400, nil, nil
	}}
	provider := NewMsiProvider(&httpClient)

	_, err := provider.GetMsi()
	if err == nil {
		t.FailNow()
	}
}
