// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package httputil

import (
	"bytes"
	"crypto/tls"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const (
	OperationGet    = "GET"
	OperationPost   = "POST"
	OperationDelete = "DELETE"
	OperationPut    = "PUT"
)

type HttpClient interface {
	Get(url string, headers map[string]string) (responseCode int, body []byte, err error)
	Post(url string, headers map[string]string, payload []byte) (responseCode int, body []byte, err error)
	Put(url string, headers map[string]string, payload []byte) (responseCode int, body []byte, err error)
	Delete(url string, headers map[string]string, payload []byte) (responseCode int, body []byte, err error)
}

// for testing
type httpClientInterface interface {
	Do(req *http.Request) (*http.Response, error)
}

type Client struct {
	httpClient    httpClientInterface
	retryBehavior RetryBehavior
}

type RetryBehavior = func(statusCode int, i int) bool

// return false to end retries
// i starts from 1 keeps getting incremented while function returns true

var NoRetry = func(statusCode int, i int) bool {
	return false
}

var LinearRetryThrice = func(statusCode int, i int) bool {
	time.Sleep(time.Second * 3)
	if i < 3 {
		return true // retry if count < 3
	}
	return false
}

var ExponentialRetryThrice = func(statusCode int, i int) bool {
	delay := time.Second * time.Duration(2^(i))
	const maxDelay time.Duration = 60 * time.Second

	if delay > maxDelay {
		delay = maxDelay
	}
	time.Sleep(delay)
	if i < 5 {
		return true // retry if count < 3
	}
	return false
}

func isTransientHttpStatusCode(statusCode int) bool {
	switch statusCode {
	case 408, 429:
		return true // timeout and too many requests
	default:
		if statusCode >= 500 && statusCode != 501 && statusCode != 505 {
			return true
		}
		return false
	}
}

func IsSuccessStatusCode(statusCode int) bool {
	switch statusCode {
	case http.StatusRequestTimeout, // 408
		http.StatusTooManyRequests,     // 429
		http.StatusInternalServerError, // 500
		http.StatusBadGateway,          // 502
		http.StatusServiceUnavailable,  // 503
		http.StatusGatewayTimeout:      // 504
		return true // timeout and too many requests
	default:
		return false
	}
}

func NewSecureHttpClient(retryBehavior RetryBehavior) HttpClient {
	tlsConfig := &tls.Config{
		Renegotiation: tls.RenegotiateFreelyAsClient,
	}

	transport := &http.Transport{TLSClientConfig: tlsConfig}
	httpClient := &http.Client{Transport: transport}
	return &Client{httpClient, retryBehavior}
}

func NewSecureHttpClientWithCertificates(certificate string, key string, retryBehavior RetryBehavior) HttpClient {
	if retryBehavior == nil {
		panic("Retry policy must be specified")
	}
	cert, err := tls.LoadX509KeyPair(certificate, key)
	if err != nil {
		log.Fatal(err)
	}

	tlsConfig := &tls.Config{
		Certificates:  []tls.Certificate{cert},
		Renegotiation: tls.RenegotiateFreelyAsClient,
	}

	transport := &http.Transport{TLSClientConfig: tlsConfig}
	httpClient := &http.Client{Transport: transport}
	return &Client{httpClient, retryBehavior}
}

func NewInsecureHttpClient(certificate string, key string, retryBehavior RetryBehavior) HttpClient {
	if retryBehavior == nil {
		panic("Retry policy must be specified")
	}
	cert, err := tls.LoadX509KeyPair(certificate, key)
	if err != nil {
		log.Fatal(err)
	}

	tlsConfig := &tls.Config{
		Certificates:       []tls.Certificate{cert},
		InsecureSkipVerify: true,
		Renegotiation:      tls.RenegotiateFreelyAsClient,
	}

	transport := &http.Transport{TLSClientConfig: tlsConfig}
	httpClient := &http.Client{Transport: transport}

	return &Client{httpClient, retryBehavior}
}

// Get issues a get request
func (client *Client) Get(url string, headers map[string]string) (responseCode int, body []byte, err error) {
	return client.issueRequest(OperationGet, url, headers, nil)
}

// Post issues a post request
func (client *Client) Post(url string, headers map[string]string, payload []byte) (responseCode int, body []byte, err error) {
	return client.issueRequest(OperationPost, url, headers, bytes.NewBuffer(payload))
}

// Put issues a put request
func (client *Client) Put(url string, headers map[string]string, payload []byte) (responseCode int, body []byte, err error) {
	return client.issueRequest(OperationPut, url, headers, bytes.NewBuffer(payload))
}

// Delete issues a delete request
func (client *Client) Delete(url string, headers map[string]string, payload []byte) (responseCode int, body []byte, err error) {
	return client.issueRequest(OperationDelete, url, headers, bytes.NewBuffer(payload))
}

func (client *Client) issueRequest(operation string, url string, headers map[string]string, payload *bytes.Buffer) (int, []byte, error) {
	request, err := http.NewRequest(operation, url, nil)
	if payload != nil && payload.Len() != 0 {
		request, err = http.NewRequest(operation, url, payload)
	}

	for key, value := range headers {
		request.Header.Add(key, value)
	}

	res, err := client.httpClient.Do(request)

	if err == nil && IsSuccessStatusCode(res.StatusCode) {
		// no need to retry
	} else if err == nil && res != nil {
		// there was no error, so look at the status code to retry
		for i := 1; isTransientHttpStatusCode(res.StatusCode) && client.retryBehavior(res.StatusCode, i); i++ {
			res, err = client.httpClient.Do(request)
			if err != nil {
				break
			}
		}
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
