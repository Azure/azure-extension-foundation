package httputil

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

const (
	operationGet    = "GET"
	operationPost   = "POST"
	operationDelete = "DELETE"
)

// Get issues a get request
var Get = func(url string, headers map[string]string) (responseCode int, body []byte, err error) {
	return issueRequest(operationGet, url, headers, bytes.NewBufferString(""))
}

// Post issues a post request
var Post = func(url string, headers map[string]string, payload []byte) (responseCode int, body []byte, err error) {
	return issueRequest(operationPost, url, headers, bytes.NewBuffer(payload))
}

// Delete issues a delete request
var Delete = func(url string, headers map[string]string, payload []byte) (responseCode int, body []byte, err error) {
	return issueRequest(operationDelete, url, headers, bytes.NewBuffer(payload))
}

func issueRequest(operation string, url string, headers map[string]string, payload *bytes.Buffer) (int, []byte, error) {
	client := &http.Client{}

	request, err := http.NewRequest(operation, url, nil)
	if payload.Len() != 0 {
		request, err = http.NewRequest(operation, url, payload)
	}

	for key, value := range headers {
		request.Header.Add(key, value)
	}

	res, err := client.Do(request)
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
