// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package msi

import (
	"encoding/json"
	"fmt"
	"github.com/Azure/azure-extension-foundation/httputil"
	"io/ioutil"
	"os"
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

func TestCanGetMsi(t *testing.T) {
	//t.Skip() // for testing on Azure VM only
	outdir := "./testoutput"
	os.Mkdir(outdir, 0777)
	secureHttpClient := httputil.NewSecureHttpClient()
	msiProvider := NewMsiProvider(&secureHttpClient)
	msi, err := msiProvider.GetMsi()
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Logf("Successfully got msi token.\nClientId was : %s", msi.ClientID)
	msiJsonBytes, err := json.Marshal(msi)
	if err != nil {
		t.Fatal(err.Error())
	}
	err = ioutil.WriteFile(fmt.Sprintf("%s/msi.json", outdir), msiJsonBytes[:], 0700)
	if err != nil {
		t.Fatal(err.Error())
	}
}
