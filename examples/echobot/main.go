package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/4kelly/go-kik/kik"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var (
	kikClient *kik.Client
	err       error
)

func init() {
	username := os.Getenv("KIKBOT_USERNAME")
	key := os.Getenv("KIKBOT_API_KEY")
	webhook := os.Getenv("KIKBOT_WEBHOOK")

	if username == "" {
		log.Fatal("!!! No KIKBOT_USERNAME set. Tests can't run !!!\n\n")
	}
	if key == "" {
		log.Fatal("!!! No KIKBOT_API_KEY set. Tests can't run !!!\n\n")
	}

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
	config := kik.
	err = kikClient.SetConfiguration(&kik.Configuration{
		Webhook:        webhook,
		Features:       &kik.Features{},
		StaticKeyboard: nil,
	})
	if err != nil {
		log.Fatalf("could not configure kik client: %v ", err)
	}
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())

	e.POST("/incoming", handleMessages)

	e.Logger.Fatal(e.Start(":5000"))
}

func handleMessages(c echo.Context) error {
	var messages kik.ReceivedMessages

	err := c.Bind(&messages)
	if err != nil {
		return err
	}

	outgoingMessages := respondTo(messages)

	err = kikClient.SendMessage(outgoingMessages)
	if err != nil {
		log.Printf("could not send message: %v ", err)
	}

	return c.String(http.StatusOK, "")
}

func respondTo(m []kik.Receive) []kik.Message {
	var outgoingMessages []kik.Message

	for _, m := range m {
		var reply kik.Message
		switch v := m.(type) {

		case *kik.TextMessageReceive:
			reply = kik.TextMessage{
				SendMessage: kik.SendMessage{
					To:   v.From,
					Type: "text",
				},
				Body: v.Body,
			}

		case *kik.PictureMessageReceive:
			reply = kik.PictureMessage{
				SendMessage: kik.SendMessage{
					To:   v.From,
					Type: "picture",
				},
				PicUrl: v.PicUrl,
			}

		default:
			fmt.Errorf("Was not able to decode message type %T", v)
		}

		outgoingMessages = append(outgoingMessages, reply)
	}
	return outgoingMessages
}
