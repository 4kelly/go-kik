package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/r-kells/go-kik/kik"
)

var (
	kikClient *kik.Client
	err       error
)

func init() {
	username := os.Getenv("IT_KIKBOT_USERNAME")
	key := os.Getenv("IT_KIKBOT_KEY")

	if username == "" {
		log.Fatal("!!! No IT_KIKBOT_USERNAME set. Tests can't run !!!\n\n")
	}
	if key == "" {
		log.Fatal("!!! No IT_KIKBOT_KEY set. Tests can't run !!!\n\n")
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
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())

	e.POST("/incoming", handleMessages)

	e.Logger.Fatal(e.Start(":8080"))
}

func handleMessages(c echo.Context) error {
	var messages kik.ReceivedMessages

	err := c.Bind(&messages)
	if err != nil {
		return err
	}

	outgoingMessages := respondTo(messages)

	kikClient.SendMessage(outgoingMessages)

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
