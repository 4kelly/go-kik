// go-kik is a client library for the [kik bot api](https://dev.kik.com/#/home).
package kik

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

const (
	GetUserUrl = "/v1/user/"
)

// KikClient is used to interface with the Kik bot API.
type KikClient struct {
	BotUsername string
	ApiKey      string
	Client      *http.Client
	BaseUrl     *url.URL
}

// NewKikClient is a simple convenience constructor for a KikClient, you do not have to use it.
func NewKikClient(baseUrl string, botUsername string, apiKey string, httpClient *http.Client) (*KikClient, error) {
	if httpClient == nil {
		httpClient = &http.Client{}
	}

	if !strings.HasSuffix(baseUrl, "/") {
		return nil, fmt.Errorf("BaseURL must have a trailing slash, but %s does not", baseUrl)
	}
	baseUrlParsed, err := url.Parse(baseUrl)
	if err != nil {
		return nil, err
	}

	return &KikClient{
		BotUsername: botUsername,
		ApiKey:      apiKey,
		Client:      httpClient,
		BaseUrl:     baseUrlParsed}, nil
}

//func (k *KikClient) setConfiguration() (http.Client, error) {
//	return apiResponse{}
//}
//
//func (k *KikClient) getConfiguration() (http.Client, error) {
//	return apiResponse{}
//}
//func (k *KikClient) sendMessages(messages []TextMessage) (http.Client, error) {
//}

//func (k *KikClient) sendBroadcast() (http.Client, error) {
//	return apiResponse{}
//}

// GetUser returns a users profile data as a User struct.
func (k *KikClient) GetUser(username string) (*User, error) {
	var user User

	req, err := k.newRequest("GET", GetUserUrl+username, nil)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(k.BotUsername, k.ApiKey)

	_, err = k.do(req, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

//
//func (k *KikClient) createCode() (http.Client, error) {
//	return apiResponse{}
//}
//
