package kik

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

func (k *KikClient) do(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := k.Client.Do(req)

	if err != nil {
		return nil, fmt.Errorf("error completing http request: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return resp, errors.New(fmt.Sprintf("status code != OK, was %d", resp.StatusCode))
	}

	defer resp.Body.Close()

	if err = json.NewDecoder(resp.Body).Decode(v); err != nil {
		return resp, fmt.Errorf("error trying to decode json into `User` struct: %v", err)
	}
	return resp, nil
}

// newRequest creates an http.Request. A relative URL is resolved relative to the BaseURL of the Client.
// Relative URLs should always be specified with a preceding slash.
// If specified, the value pointed to by body is JSON encoded and included as the request body.
func (k *KikClient) newRequest(method, urlStr string, body interface{}) (*http.Request, error) {

	parsedUrl, err := k.BaseUrl.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(false)
		err := enc.Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, parsedUrl.String(), buf)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	return req, nil
}
