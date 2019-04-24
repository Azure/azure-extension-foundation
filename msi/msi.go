// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package msi

import (
	"encoding/json"
	"fmt"
	"github.com/Azure/azure-extension-foundation/httputil"
	"strconv"
	"time"
)

const metadataMsiURL = "http://169.254.169.254/metadata/identity/oauth2/token?api-version=2018-02-01&resource=https://management.core.windows.net/"

type Msi struct {
	AccessToken  string `json:"access_token"`
	ClientID     string `json:"client_id"`
	ExpiresIn    string `json:"expires_in"`
	ExpiresOn    string `json:"expires_on"` // expressed in seconds from epoch
	ExtExpiresIn string `json:"ext_expires_in"`
	NotBefore    string `json:"not_before"`
	Resource     string `json:"resource"`
	TokenType    string `json:"token_type"`
}

type MsiProvider interface {
	GetMsi() (Msi, error)
}

type provider struct {
	httpClient httputil.HttpClient
}

func NewMsiProvider(client httputil.HttpClient) provider {
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

// check expiry of MSI token based on time
func (msi *Msi) IsMsiTokenExpired() (bool, error) {
	expiryTime, err := msi.GetExpiryTime()
	if err != nil {
		return false, err
	}

	// Consider token expired 2 minutes before expiry time
	expiryTime = expiryTime.Add(-2 * time.Minute)

	if time.Now().After(expiryTime) {
		return true, nil
	} else {
		return false, nil
	}
}

func (msi *Msi) GetExpiryTime() (time.Time, error) {
	expiryTimeInSeconds, err := strconv.ParseInt(msi.ExpiresOn, 10, 64)
	if err != nil {
		return time.Unix(0, 0), err
	}
	expiryTime := time.Unix(expiryTimeInSeconds, 0)
	return expiryTime, nil
}

func (msi *Msi) GetJson() (string, error) {
	jsonBytes, err := json.Marshal(msi)
	return string(jsonBytes[:]), err
}
