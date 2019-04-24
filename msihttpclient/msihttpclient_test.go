// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package msihttpclient

import (
	"fmt"
	"github.com/Azure/azure-extension-foundation/httputil"
	"github.com/Azure/azure-extension-foundation/metadata"
	"github.com/Azure/azure-extension-foundation/msi"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"testing"
	"time"
)

var getDummyMsiFunc = func() *msi.Msi {
	return &msi.Msi{
		AccessToken:  "dummy access token",
		ExpiresIn:    "1234",
		NotBefore:    time.Now().String(),
		ClientID:     "dummy client Id",
		ExpiresOn:    strconv.FormatInt(time.Now().Unix(), 10),
		ExtExpiresIn: time.Duration(time.Hour).String(),
		Resource:     "dummy resource",
		TokenType:    "Bearer",
	}
}

type mockMsiProvider struct {
	dummyMsi     *msi.Msi
	timesInvoked int
}

func (prov *mockMsiProvider) GetMsi() (msi.Msi, error) {
	prov.timesInvoked++
	if prov.dummyMsi == nil {
		prov.dummyMsi = getDummyMsiFunc()
	}
	expiryTime, _ := prov.dummyMsi.GetExpiryTime()
	prov.dummyMsi.ExpiresOn = strconv.FormatInt(expiryTime.Add(1*time.Minute).Unix(), 10)
	return *prov.dummyMsi, nil
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

var mdata = metadata.Metadata{
	Compute: metadata.MetadataCompute{
		VmId:              "vmid",
		SubscriptionId:    "subId",
		ResourceGroupName: "resourceGroupName",
		Name:              "vmName",
	},
	Network: metadata.MetadataNetwork{},
}

func TestAddVmIdQueryParameterToUrl(t *testing.T) {
	dummyMsi := getDummyMsiFunc()
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
	msiHttp := msiHttpClient{httpClient: getHttpClientFunc(), retryBehavior: httputil.DefaultRetryBehavior, msiProvider: &mockMsiProvider{timesInvoked: 0}, metadata: &mdata}
	modifiedUrl, err := msiHttp.addVmIdQueryParameterToUrl("http://foo.bar.com?query1=val1&query2=val2&speed=100")
	if err != nil {
		t.Fatal(err)
	}
	if len(modifiedUrl) == 0 {
		t.Fatal(fmt.Errorf("modified url was of length 0"))
	}

	u, _ := url.Parse(modifiedUrl)

	if u.Query().Get("vmResourceId") != mdata.GetAzureResourceId() {
		t.Fatal(fmt.Errorf("the modified query does not contain the query parameter for ARM id"))
	}
}

func TestNewMsiHttpClientHeaders(t *testing.T) {
	dummyMsi := getDummyMsiFunc()
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
	msiHttp := NewMsiHttpClient(&mockMsiProvider{timesInvoked: 0}, &mdata, httputil.DefaultRetryBehavior)
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
					return &http.Response{StatusCode: 429, Body: noBody{}}, nil
				default:
					return &http.Response{StatusCode: 200, Body: noBody{}}, nil
				}
			},
		}
	}
	msiHttp := NewMsiHttpClient(&mockMsi, &mdata, httputil.DefaultRetryBehavior)
	msiHttp.Get("", make(map[string]string))
	if mockMsi.timesInvoked < 3 {
		t.Fatal("retry logic didn't invoke msiProvider for retries")
	}
}
