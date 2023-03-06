package http

import (
	"bytes"
	"context"
	"crypto/tls"
	"log"
	"net/http"
	"net/http/httputil"
)

func Request(ctx context.Context, method string, url string, headers map[string]string, body []byte) (*http.Response, error) {

	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		log.Println("[ERROR] HTTP Create NewRequest error:", err)
		return nil, err
	}

	req = req.WithContext(ctx)

	for k, v := range headers {
		req.Header.Add(k, v)
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: tr}
	resp, err := client.Do(req)

	if err != nil {
		log.Println("[ERROR] HTTP request error:", err)
		return nil, err
	}
	dResp, _ := httputil.DumpResponse(resp, true)
	log.Printf("[INFO] HTTP Response %s", dResp)

	return resp, err
}
