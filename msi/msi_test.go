// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package msi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/Azure/azure-extension-foundation/httputil"
)

func TestSuccessfulGetMsi(t *testing.T) {
	const tokenValue = "token"
	httpClient := httputil.MockHttpClient{Getfunc: func(url string, headers map[string]string) (responseCode int, body []byte, err error) {
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
	// metadata service will return 400 if MSI is disabled
	httpClient := httputil.MockHttpClient{Getfunc: func(url string, headers map[string]string) (responseCode int, body []byte, err error) {
		return 400, nil, nil
	}}
	provider := NewMsiProvider(&httpClient)

	_, err := provider.GetMsi()
	if err == nil {
		t.FailNow()
	}
}

func TestCanGetMsi(t *testing.T) {
	t.Skip() // for testing on Azure VM only
	outdir := "./testoutput"
	os.Mkdir(outdir, 0777)
	secureHttpClient := httputil.NewSecureHttpClient(httputil.NoRetry)
	msiProvider := NewMsiProvider(secureHttpClient)
	msi, err := msiProvider.GetMsi()
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Logf("Successfully got msi token.\nClientId was : %s", msi.ClientID)
	msiJsonBytes, err := json.Marshal(msi)
	if err != nil {
		t.Fatal(err.Error())
	}
	err = ioutil.WriteFile(fmt.Sprintf("%s/msi1.json", outdir), msiJsonBytes[:], 0700)
	if err != nil {
		t.Fatal(err.Error())
	}
}

func TestCanGetMsiForStorage(t *testing.T) {
	t.Skip() // for testing on Azure VM only
	outdir := "./testoutput"
	os.Mkdir(outdir, 0777)
	secureHttpClient := httputil.NewSecureHttpClient(httputil.NoRetry)
	msiProvider := NewMsiProvider(secureHttpClient)
	msi, err := msiProvider.GetMsiForResource("https://storage.azure.com/")
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Logf("Successfully got msi token.\nClientId was : %s", msi.ClientID)
	msiJsonBytes, err := json.Marshal(msi)
	if err != nil {
		t.Fatal(err.Error())
	}
	err = ioutil.WriteFile(fmt.Sprintf("%s/msi2.json", outdir), msiJsonBytes[:], 0700)
	if err != nil {
		t.Fatal(err.Error())
	}
}

func TestCanGetMsiForKeyVault(t *testing.T) {
	t.Skip() // for testing on Azure VM only
	outdir := "./testoutput"
	os.Mkdir(outdir, 0777)
	secureHttpClient := httputil.NewSecureHttpClient(httputil.NoRetry)
	msiProvider := NewMsiProvider(secureHttpClient)
	msi, err := msiProvider.GetMsiForResource("https://vault.azure.net")
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Logf("Successfully got msi token.\nClientId was : %s", msi.ClientID)
	msiJsonBytes, err := json.Marshal(msi)
	if err != nil {
		t.Fatal(err.Error())
	}
	err = ioutil.WriteFile(fmt.Sprintf("%s/msi3.json", outdir), msiJsonBytes[:], 0700)
	if err != nil {
		t.Fatal(err.Error())
	}
}

func TestCanGetMsiForKeyVaultWithClientId(t *testing.T) {
	t.Skip() // for testing on Azure VM only
	outdir := "./testoutput"
	os.Mkdir(outdir, 0777)
	secureHttpClient := httputil.NewSecureHttpClient(httputil.NoRetry)
	msiProvider := NewMsiProvider(secureHttpClient)
	// change client id to match your user managed identity
	msi, err := msiProvider.GetMsiUsingClientId("31b403aa-c364-4240-a7ff-d85fb6cd7232", "https://vault.azure.net")
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Logf("Successfully got msi token.\nClientId was : %s", msi.ClientID)
	msiJsonBytes, err := json.Marshal(msi)
	if err != nil {
		t.Fatal(err.Error())
	}
	err = ioutil.WriteFile(fmt.Sprintf("%s/msi4.json", outdir), msiJsonBytes[:], 0700)
	if err != nil {
		t.Fatal(err.Error())
	}
}

func TestCanGetMsiForStoragetWithObjectId(t *testing.T) {
	t.Skip() // for testing on Azure VM only
	outdir := "./testoutput"
	os.Mkdir(outdir, 0777)
	secureHttpClient := httputil.NewSecureHttpClient(httputil.NoRetry)
	msiProvider := NewMsiProvider(secureHttpClient)
	// change object id to match your user managed identity
	msi, err := msiProvider.GetMsiUsingObjectId("20931e04-e65f-4526-8c01-9d627f577263", "https://storage.azure.com/")
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Logf("Successfully got msi token.\nClientId was : %s", msi.ClientID)
	msiJsonBytes, err := json.Marshal(msi)
	if err != nil {
		t.Fatal(err.Error())
	}
	err = ioutil.WriteFile(fmt.Sprintf("%s/msi5.json", outdir), msiJsonBytes[:], 0700)
	if err != nil {
		t.Fatal(err.Error())
	}
}

func TestCanGetImdsVariable(t *testing.T) {

	//without setting env var - should return metadataURL
	url := GetMetadataIdentityURL()
	if url != metadataIdentityURL {
		t.Fatal("mismatch URLs")
	}

	//setting env var - should return test (or whatever var is set to)
	os.Setenv(identityEnvVar, "test")
	url = GetMetadataIdentityURL()
	if url != "test" {
		t.Fatal("mismatch URL")
	}

	os.Unsetenv(identityEnvVar)
}
