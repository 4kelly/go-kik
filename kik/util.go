package kik

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
)

func (k *Client) do(req *http.Request, v interface{}) error {
	resp, err := k.Client.Do(req)

	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		b, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("%v: %s %s %s returned: <%v> %s",
			HttpError, req.Method, req.URL, req.Body, resp.StatusCode, b)
	}

	if v != nil {
		defer resp.Body.Close()

		if err = json.NewDecoder(resp.Body).Decode(v); err != nil {
			return fmt.Errorf("error trying to decode json into `v` struct: %v", err)
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

// validateMessageTypes ensures that the message interfaces that are passed
// conform to valid message types. This is the best alternative without generics.
func validateMessageTypes(typesToTest []interface{}) error {
	for _, t := range typesToTest {

		switch messageType := t.(type) {
		case TextMessage, PictureMessage, LinkMessage, VideoMessage:
			err := validateKeyboardResponseTypes(messageType)
			if err != nil {
				return err
			}
			continue
		default:
			return fmt.Errorf("type: %T value:%v %w", messageType, t, NotMessageTypeError)
		}
	}
	return nil
}

/*
validateKeyboardResponseTypes scans through the known types of Messages and type checks the keyboard responses at runtime.

WARNING: This relies on the structure of the Message types.

Notes:
- You could recursively go through a type switch, but that would be non-linear to follow for the reader.
- This is the best alternative without generics.
*/
func validateKeyboardResponseTypes(message interface{}) error {
	sliceOfKeyboards := reflect.Indirect(reflect.ValueOf(message).FieldByName("SendMessage")).FieldByName("Keyboards")
	// For keyboard in keyboards
	for i := 0; i < sliceOfKeyboards.Len(); i++ {
		keyboardResponses := sliceOfKeyboards.Index(i).FieldByName("Responses")
		if keyboardResponses.Len() == 0 {
			continue
		}
		// For response in reponses
		for i := 0; i < keyboardResponses.Len(); i++ {
			r := keyboardResponses.Index(i).Interface()

			switch keyboardResponseType := r.(type) {
			case KeyboardTextResponse, KeyboardPictureResponse, KeyboardFriendPickerResponse:
				continue
			default:
				return fmt.Errorf("unrecongized type: %T value:%v %w", keyboardResponseType, keyboardResponseType, NotMessageTypeError)
			}
		}
	}
	return nil
}
