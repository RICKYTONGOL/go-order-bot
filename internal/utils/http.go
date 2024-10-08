// utils package contains utility functions
package utils

import (
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

func DoPost(url string, contentType string, body io.Reader) ([]byte, error) {
	tr := &http.Transport{
		IdleConnTimeout: 120 * time.Second,
		DialContext: (&net.Dialer{
			Timeout:   120 * time.Second,
			KeepAlive: 120 * time.Second,
		}).DialContext,
		DisableKeepAlives:      false,
		ResponseHeaderTimeout:  10 * time.Minute,
		ExpectContinueTimeout:  10 * time.Minute,
		ReadBufferSize:         1024 * 1024 * 100,
		WriteBufferSize:        1024 * 1024 * 100,
		TLSHandshakeTimeout:    10 * time.Minute,
		MaxResponseHeaderBytes: 1024 * 1024 * 1000,
	}

	client := &http.Client{
		Timeout:   time.Second * 120,
		Transport: tr,
	}

	resp, err := client.Post(url, contentType, body)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	// if resp.ContentLength > 100000 {
	// 	buf := make([]byte, resp.ContentLength)
	// 	n, _ := io.ReadFull(resp.Body, buf)
	// 	n = n + 1
	// }

	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return resBody, nil
}
