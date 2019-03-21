// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package msihttpclient

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"github.com/Azure/azure-extension-foundation/httputil"
	"github.com/Azure/azure-extension-foundation/msi"
	"io/ioutil"
	"net/http"
)

type msiHttpClient struct {
	httpClient    *http.Client
	retryBehavior httputil.RetryBehavior
	msi           *msi.Msi
	msiProvider   msi.MsiProvider
}

func NewMsiHttpClient(retryBehavior httputil.RetryBehavior) httputil.HttpClient {
	tlsConfig := &tls.Config{
		Renegotiation: tls.RenegotiateFreelyAsClient,
	}

	transport := &http.Transport{TLSClientConfig: tlsConfig}
	httpClient := &http.Client{Transport: transport}
	msiProvider := msi.NewMsiProvider(httputil.NewSecureHttpClient(retryBehavior))
	return &msiHttpClient{httpClient, retryBehavior, nil, &msiProvider}
}

func (client *msiHttpClient) Get(url string, headers map[string]string) (responseCode int, body []byte, err error) {
	return client.issueRequest(httputil.OperationGet, url, headers, nil)
}

// Post issues a post request
func (client *msiHttpClient) Post(url string, headers map[string]string, payload []byte) (responseCode int, body []byte, err error) {
	return client.issueRequest(httputil.OperationPost, url, headers, bytes.NewBuffer(payload))
}

// Put issues a put request
func (client *msiHttpClient) Put(url string, headers map[string]string, payload []byte) (responseCode int, body []byte, err error) {
	return client.issueRequest(httputil.OperationPut, url, headers, bytes.NewBuffer(payload))
}

// Delete issues a delete request
func (client *msiHttpClient) Delete(url string, headers map[string]string, payload []byte) (responseCode int, body []byte, err error) {
	return client.issueRequest(httputil.OperationDelete, url, headers, bytes.NewBuffer(payload))
}

func (client *msiHttpClient) issueRequest(operation string, url string, headers map[string]string, payload *bytes.Buffer) (int, []byte, error) {
	request, err := http.NewRequest(operation, url, nil)
	if payload != nil && payload.Len() != 0 {
		request, err = http.NewRequest(operation, url, payload)
	}

	// Initialize msi is required
	if client.msi == nil {
		msi, err := client.msiProvider.GetMsi()
		if err != nil {
			return -1, nil, err
		}
		client.msi = &msi
	}

	// Add authorization if required
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", client.msi.AccessToken))
	for key, value := range headers {
		request.Header.Add(key, value)
	}

	res, err := client.httpClient.Do(request)

	for i := 1; err != nil && (res.StatusCode == 401 || client.retryBehavior(res.StatusCode, i)); i++ {
		// refresh certificate if required
		if res.StatusCode == 401 {
			msi, err := client.msiProvider.GetMsi()
			if err != nil {
				return -1, nil, err
			}
			client.msi = &msi
			request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", client.msi.AccessToken))
		}
		res, err = client.httpClient.Do(request)
	}

	if err != nil {
		return -1, nil, err
	}

	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	code := res.StatusCode
	if err != nil {
		return -1, nil, err
	}

	return code, body, nil
}
