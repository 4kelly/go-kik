// go-kik is a client library for the [kik bot api](https://dev.kik.com/#/home).
// Documentation can be found [here](https://dev.kik.com/#/docs/messaging).
package kik

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

const (
	GetUserUrl     = "/v1/user/"
	SendMessageUrl = "/v1/message"
	BroadcastUrl   = "/v1/broadcast"
	ConfigtUrl     = "/v1/config"
)

// Client is used to interface with the Kik bot API.
type Client struct {
	BotUsername string
	ApiKey      string
	Client      *http.Client
	BaseUrl     *url.URL
}

// NewKikClient is a simple convenience constructor for a Client, you do not have to use it.
func NewKikClient(baseUrl string, botUsername string, apiKey string, httpClient *http.Client) (*Client, error) {
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

	return &Client{
		BotUsername: botUsername,
		ApiKey:      apiKey,
		Client:      httpClient,
		BaseUrl:     baseUrlParsed}, nil
}

func (k *Client) SetConfiguration(c *Configuration) error {
	req, err := k.newRequest("POST", ConfigtUrl, c)
	if err != nil {
		return err
	}

	req.SetBasicAuth(k.BotUsername, k.ApiKey)

	err = k.do(req, &c)
	if err != nil {
		return err
	}
	return nil
}

func (k *Client) GetConfiguration() (*Configuration, error) {
	var config Configuration

	req, err := k.newRequest("GET", ConfigtUrl, nil)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(k.BotUsername, k.ApiKey)

	err = k.do(req, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func (k *Client) SendMessage(messages []interface{}) error {
	err := validateMessageTypes(messages)
	if err != nil {
		return err
	}

	type m struct {
		Messages []interface{} `json:"messages"`
	}
	payload := m{Messages: messages}

	req, err := k.newRequest("POST", SendMessageUrl, payload)
	if err != nil {
		return err
	}

	req.SetBasicAuth(k.BotUsername, k.ApiKey)

	return k.do(req, nil)
}

func (k *Client) BroadcastMessage(messages []interface{}) error {
	err := validateMessageTypes(messages)
	if err != nil {
		return err
	}

	type m struct {
		Messages []interface{} `json:"messages"`
	}
	payload := m{Messages: messages}

	req, err := k.newRequest("POST", BroadcastUrl, payload)
	if err != nil {
		return err
	}

	req.SetBasicAuth(k.BotUsername, k.ApiKey)

	return k.do(req, nil)
}

// GetUser returns a users profile data as a User struct.
func (k *Client) GetUser(username string) (*User, error) {
	var user User

	req, err := k.newRequest("GET", GetUserUrl+username, nil)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(k.BotUsername, k.ApiKey)

	err = k.do(req, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (k *Client) createCode() (Code, error) {
	
}
