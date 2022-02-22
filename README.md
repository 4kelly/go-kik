# kikbot

A Go client library for the [Kik bot API](https://dev.kik.com/#/home).

Example usage in [here](test/system/kik_test.go).

```go
package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/4kelly/go-kik/kik"
	"github.com/google/go-cmp/cmp"
)

var (
	kikClient *kik.Client
	err       error
)

func init() {
	username := os.Getenv("KIKBOT_USERNAME")
	key := os.Getenv("KIKBOT_API_KEY")
	webhook := os.Getenv("KIKBOT_WEBHOOK")

	client := &http.Client{
		Transport:     nil,
		CheckRedirect: nil,
		Jar:           nil,
		Timeout:       time.Duration(3) * time.Second,
	}

	kikClient, err = kik.NewKikClient(
		"https://api.kik.com/",
		username,
		key,
		client,
	)
	if err != nil {
		log.Fatalf("could not initiate client: %v ", err)
	}

	err = kikClient.SetConfiguration(&kik.Configuration{
		Webhook:        webhook,
		Features:       kik.Features{},
		StaticKeyboard: nil,
	})
	if err != nil {
		log.Fatalf("could not configure kik client: %v ", err)
	}
}
```

## Testing

Run system tests to validate integrity of the Kik API.
```go
go test ./...
```