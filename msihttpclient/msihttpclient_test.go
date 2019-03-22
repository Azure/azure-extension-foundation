// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package msihttpclient

import (
	"fmt"
	"github.com/Azure/azure-extension-foundation/httputil"
	"github.com/Azure/azure-extension-foundation/msi"
	"io"
	"net/http"
	"testing"
	"time"
)

var dummyMsi = msi.Msi{
	AccessToken:  "dummy access token",
	ExpiresIn:    "1234",
	NotBefore:    time.Now().String(),
	ClientID:     "dummy client Id",
	ExpiresOn:    time.Now().Add(time.Duration(time.Hour)).String(),
	ExtExpiresIn: time.Duration(time.Hour).String(),
	Resource:     "dummy resource",
	TokenType:    "Bearer",
}

type mockMsiProvider struct {
	timesInvoked int
}

func (prov *mockMsiProvider) GetMsi() (msi.Msi, error) {
	prov.timesInvoked++
	return dummyMsi, nil
}

type mockHttpClient struct {
	AttemptCount *int
	DoFunc       func(i *int, req *http.Request) (*http.Response, error)
}

func (client *mockHttpClient) Do(req *http.Request) (*http.Response, error) {
	return client.DoFunc(client.AttemptCount, req)
}

type noBody struct {
}

func (noBody) Read(bytes []byte) (int, error)   { return 0, io.EOF }
func (noBody) Close() error                     { return nil }
func (noBody) WriteTo(io.Writer) (int64, error) { return 0, nil }

func TestNewMsiHttpClientHeaders(t *testing.T) {
	getHttpClientFunc = func() httpClientInterface {
		return &mockHttpClient{
			DoFunc: func(i *int, req *http.Request) (*http.Response, error) {
				authorization := req.Header.Get("Authorization")
				if authorization != fmt.Sprintf("Bearer %s", dummyMsi.AccessToken) {
					t.Fatal("authorization header didn't match")
				}
				return &http.Response{StatusCode: 200, Body: noBody{}}, nil
			},
		}
	}
	msiHttp := NewMsiHttpClient(&mockMsiProvider{timesInvoked: 0}, httputil.ExponentialRetryThrice)
	msiHttp.Get("", make(map[string]string))
}

func TestRetryLogic(t *testing.T) {
	mockMsi := mockMsiProvider{timesInvoked: 0}
	i := 0
	getHttpClientFunc = func() httpClientInterface {
		return &mockHttpClient{
			AttemptCount: &i,
			DoFunc: func(i *int, req *http.Request) (*http.Response, error) {
				(*i)++
				switch *i {
				case 1, 2:
					return &http.Response{StatusCode: 401, Body: noBody{}}, nil
				default:
					return &http.Response{StatusCode: 200, Body: noBody{}}, nil
				}
			},
		}
	}
	msiHttp := NewMsiHttpClient(&mockMsi, httputil.ExponentialRetryThrice)
	msiHttp.Get("", make(map[string]string))
	if mockMsi.timesInvoked < 2 {
		t.Fatal("retry logic didn't invoke msiProvider for retries")
	}
}
