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
	operationPut = "PUT"
)

type InbuiltHttpClient struct {
	httpClient *http.Client
}

type HttpClient interface {
	Get(url string, headers map[string]string) (responseCode int, body []byte, err error)
	Post(url string, headers map[string]string, payload []byte) (responseCode int, body []byte, err error)
	Put(url string, headers map[string]string, payload []byte) (responseCode int, body []byte, err error)
	Delete(url string, headers map[string]string, payload []byte) (responseCode int, body []byte, err error)
}

func NewSecureHttpClient() HttpClient {
	tlsConfig := &tls.Config{
		Renegotiation: tls.RenegotiateFreelyAsClient,
	}

	transport := &http.Transport{TLSClientConfig: tlsConfig}
	httpClient := &http.Client{Transport: transport}
	return &InbuiltHttpClient{httpClient}
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
	return &InbuiltHttpClient{httpClient}
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

	return &InbuiltHttpClient{httpClient}
}

// Get issues a get request
func (client *InbuiltHttpClient) Get(url string, headers map[string]string) (responseCode int, body []byte, err error) {
	return client.issueRequest(operationGet, url, headers, nil)
}

// Post issues a post request
func (client *InbuiltHttpClient) Post(url string, headers map[string]string, payload []byte) (responseCode int, body []byte, err error) {
	return client.issueRequest(operationPost, url, headers, bytes.NewBuffer(payload))
}

func (client *InbuiltHttpClient) Put(url string, headers map[string]string, payload []byte) (responseCode int, body []byte, err error) {
	return client.issueRequest(operationPut, url, headers, bytes.NewBuffer(payload))
}

// Delete issues a delete request
func (client *InbuiltHttpClient) Delete(url string, headers map[string]string, payload []byte) (responseCode int, body []byte, err error) {
	return client.issueRequest(operationDelete, url, headers, bytes.NewBuffer(payload))
}

func (client *InbuiltHttpClient) issueRequest(operation string, url string, headers map[string]string, payload *bytes.Buffer) (int, []byte, error) {
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
