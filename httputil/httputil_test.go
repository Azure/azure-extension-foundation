// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package httputil

import (
	"io"
	"net/http"
	"testing"
)

type mockHttpClient struct {
	AttemptCount *int
	DoFunc       func(i *int, req *http.Request) (*http.Response, error)
}

func (client *mockHttpClient) Do(req *http.Request) (*http.Response, error) {
	*client.AttemptCount++
	return client.DoFunc(client.AttemptCount, req)
}

type noBody struct {
}

var return401 = func(i *int, req *http.Request) (*http.Response, error) {
	return &http.Response{StatusCode: 401, Body: noBody{}}, nil
}

func (noBody) Read(bytes []byte) (int, error)   { return 0, io.EOF }
func (noBody) Close() error                     { return nil }
func (noBody) WriteTo(io.Writer) (int64, error) { return 0, nil }

func TestRetryNever(t *testing.T) {
	attemptCount := 0
	mockClient := mockHttpClient{&attemptCount, return401}
	client := Client{&mockClient, NoRetry}
	client.Get("fake address", make(map[string]string))
	if *mockClient.AttemptCount != 1 {
		t.Fatal("Retry was attemped when none was specified")
	}
}

func TestRetryThrice(t *testing.T) {
	attemptCount := 0
	mockClient := mockHttpClient{&attemptCount, return401}
	client := Client{&mockClient, LinearRetryThrice}
	client.Get("fake address", make(map[string]string))
	if *mockClient.AttemptCount != 3 {
		t.Fatal("httpclient didn't retry thrice")
	}
}
