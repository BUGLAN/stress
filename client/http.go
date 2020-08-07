package client

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"io"
	"net/http"
	"net/textproto"
	"time"
)

type HTTPClient struct {
}

func NewHTTPClient() *HTTPClient {
	return &HTTPClient{}
}

type HTTP interface {
	Get(url string, header http.Header) error
	Post(url string, header http.Header, body, v interface{}) error
}

func (c *HTTPClient) Get(url string, header http.Header) error {
	rc, err := c.request("GET", url, header, bytes.NewBuffer(nil), time.Second*30)
	if err != nil || rc == nil {
		return err
	}
	defer func() { _ = rc.Close() }()
	return nil
}

// Post 发起post请求
func (c *HTTPClient) Post(url string, header http.Header, body, v interface{}) error {
	data, err := json.Marshal(body)
	if err != nil {
		return err
	}
	rc, err := c.request("POST", url, header, bytes.NewBuffer(data), 0)
	if err != nil || rc == nil {
		return err
	}
	defer func() { _ = rc.Close() }()
	return nil
}

func (c *HTTPClient) request(method, url string, header http.Header, body io.Reader, timeout time.Duration) (io.ReadCloser, error) {
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		Timeout: timeout,
	}
	// new a request
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	// set headers
	for k, v := range header {
		req.Header[textproto.CanonicalMIMEHeaderKey(k)] = v
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}
