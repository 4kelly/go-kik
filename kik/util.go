package kik

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

func (k *Client) do(req *http.Request, v interface{}) error {
	resp, err := k.Client.Do(req)

	if err != nil {
		return err
	}

	// defer resp.Body.Close()
	// b, _ := ioutil.ReadAll(resp.Body)
	// fmt.Printf("WTF %s END", string(b))

	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		b, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("%v: %s %s %s returned: <%v> %s",
			HttpError, req.Method, req.URL, req.Body, resp.StatusCode, b)
	}

	if v != nil {
		defer resp.Body.Close()

		if err = json.NewDecoder(resp.Body).Decode(v); err != nil {
			return fmt.Errorf("error trying to decode json into struct: %v", err)
		}
	}
	return nil
}

// newRequest creates an http.Request. A relative URL is resolved relative to the BaseURL of the Client.
// Relative URLs should always be specified with a preceding slash.
// If specified, the value pointed to by body is JSON encoded and included as the request body.
func (k *Client) newRequest(method, urlStr string, body interface{}) (*http.Request, error) {

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

	log.Printf("%s %s %s", method, parsedUrl.String(), buf)

	req, err := http.NewRequest(method, parsedUrl.String(), buf)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	return req, nil
}
