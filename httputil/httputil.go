package httputil

import (
	"bytes"
	"crypto/tls"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	operationGet    = "GET"
	operationPost   = "POST"
	operationDelete = "DELETE"
)

type HttpClient struct {
	httpClient *http.Client
}

func NewSecureHttpClient() HttpClient {
	tlsConfig := &tls.Config{
		Renegotiation: tls.RenegotiateFreelyAsClient,
	}

	transport := &http.Transport{TLSClientConfig: tlsConfig}
	httpClient := &http.Client{Transport: transport}
	return HttpClient{httpClient}
}

func NewSecureHttpClientWithCertificates(certificate string, key string) HttpClient {
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
	return HttpClient{httpClient}
}

func NewInsecureHttpClientWithCertificates(certificate string, key string) HttpClient {
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

	return HttpClient{httpClient}
}

// Get issues a get request
func (client *HttpClient) Get(url string, headers map[string]string) (responseCode int, body []byte, err error) {
	return client.issueRequest(operationGet, url, headers, nil)
}

// Post issues a post request
func (client *HttpClient) Post(url string, headers map[string]string, payload []byte) (responseCode int, body []byte, err error) {
	return client.issueRequest(operationPost, url, headers, bytes.NewBuffer(payload))
}

// Delete issues a delete request
func (client *HttpClient) Delete(url string, headers map[string]string, payload []byte) (responseCode int, body []byte, err error) {
	return client.issueRequest(operationDelete, url, headers, bytes.NewBuffer(payload))
}

func (client *HttpClient) issueRequest(operation string, url string, headers map[string]string, payload *bytes.Buffer) (int, []byte, error) {
	request, err := http.NewRequest(operation, url, nil)
	if payload != nil && payload.Len() != 0 {
		request, err = http.NewRequest(operation, url, payload)
	}

	for key, value := range headers {
		request.Header.Add(key, value)
	}

	res, err := client.httpClient.Do(request)
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
