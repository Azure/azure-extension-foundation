// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package msi

import (
	"encoding/json"
	"fmt"
)

const metadataMsiURL = "http://169.254.169.254/metadata/identity/oauth2/token?api-version=2018-02-01&resource=https://management.core.windows.net/"

type Msi struct {
	AccessToken  string `json:"access_token"`
	ClientID     string `json:"client_id"`
	ExpiresIn    string `json:"expires_in"`
	ExpiresOn    string `json:"expires_on"`
	ExtExpiresIn string `json:"ext_expires_in"`
	NotBefore    string `json:"not_before"`
	Resource     string `json:"resource"`
	TokenType    string `json:"token_type"`
}

type provider struct {
	httpClient httpClient
}

type httpClient interface {
	Get(url string, headers map[string]string) (responseCode int, body []byte, err error)
}

func NewMsiProvider(client httpClient) provider {
	return provider{httpClient: client}
}

func (p *provider) GetMsi() (Msi, error) {
	var msi = Msi{}
	code, body, err := p.httpClient.Get(metadataMsiURL, map[string]string{"Metadata": "true"})
	if err != nil {
		return msi, err
	}

	if code != 200 {
		return msi, fmt.Errorf("unable to get msi, metadata service response code %v", code)
	}

	err = json.Unmarshal(body, &msi)
	if err != nil {
		return msi, fmt.Errorf("unable to deserialize metadata service response")
	}

	return msi, nil
}

func (msi *Msi) GetJson() string {
	jsonBytes, err := json.Marshal(msi)
	if err != nil {
		panic(err.Error())
	}
	return string(jsonBytes[:])
}