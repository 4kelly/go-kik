package integration

import (
	"github.com/r-kells/go-kik/kik"
	"log"
	"net/http"
	"os"
	"testing"
	"time"
)

var (
	kikClient    *kik.Client
	err          error
	testUserName string
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

	// Who we send test messages to
	testUserName = "rmdkelly"
}

func TestGetUser_HappyPath(t *testing.T) {

	_, err := kikClient.GetUser(testUserName)
	if err != nil {
		t.Errorf("Could not find user: %v ", err)
	}
}

func TestSendMessage_HappyPath(t *testing.T) {
	keyboard := []kik.SuggestedResponseKeyboard{
		{Type: "suggested",
			Responses: []interface{}{
				kik.KeyboardTextResponse{
					Type: "text",
					Body: "KeyboardTextResponse",
				},
			},
		},
	}
	msgs := []interface{}{
		kik.TextMessage{
			SendMessage: kik.SendMessage{
				To:   testUserName,
				Type: "text",
			},
			Body: "Test_SendMessage",
		},
		kik.PictureMessage{
			SendMessage: kik.SendMessage{
				To:   testUserName,
				Type: "picture",
			},
			PicUrl: "https://i.imgur.com/TsoLODG.png",
		},
		kik.LinkMessage{
			SendMessage: kik.SendMessage{
				To:   testUserName,
				Type: "link",
			},
			Url:       "https://duckduckgo.com/",
			PicUrl:    "https://i.imgur.com/hp5ix8B.jpg",
			Title:     "Link Test",
			Text:      "Such testing.",
			NoForward: false,
		},
		kik.VideoMessage{
			SendMessage: kik.SendMessage{
				To:        testUserName,
				Type:      "video",
				Delay:     1,
				Keyboards: keyboard,
			},
			VideoUrl: "https://media.tenor.com/videos/a912c20be335cfd78610916c97198438/mp4",
			Loop:     true,
			Muted:    false,
			Autoplay: true,
			NoSave:   false,
			Attribution: &kik.Attribution{
				Name: "Attribution Test",
			},
		},
	}

	err := kikClient.SendMessage(msgs)
	if err != nil {
		t.Errorf("Error while trying to send a message. %v.", err)
	}
}

func TestBroadcastMessage_HappyPath(t *testing.T) {

	keyboard := []kik.SuggestedResponseKeyboard{
		{Type: "suggested",
			Responses: []interface{}{
				kik.KeyboardTextResponse{
					Type: "text",
					Body: "KeyboardTextResponse",
				},
			},
		},
	}
	msgs := []interface{}{
		kik.TextMessage{
			SendMessage: kik.SendMessage{
				To:   testUserName,
				Type: "text",
			},
			Body: "Test_BroadcastMessage",
		},
		kik.PictureMessage{
			SendMessage: kik.SendMessage{
				To:   testUserName,
				Type: "picture",
			},
			PicUrl: "https://i.imgur.com/TsoLODG.png",
		},
		kik.LinkMessage{
			SendMessage: kik.SendMessage{
				To:   testUserName,
				Type: "link",
			},
			Url:       "https://duckduckgo.com/",
			PicUrl:    "https://i.imgur.com/hp5ix8B.jpg",
			Title:     "Link Test",
			Text:      "Such testing.",
			NoForward: false,
		},
		kik.VideoMessage{
			SendMessage: kik.SendMessage{
				To:        testUserName,
				Type:      "video",
				Delay:     1,
				Keyboards: keyboard,
			},
			VideoUrl: "https://media.tenor.com/videos/a912c20be335cfd78610916c97198438/mp4",
			Loop:     true,
			Muted:    false,
			Autoplay: true,
			NoSave:   false,
			Attribution: &kik.Attribution{
				Name: "Attribution Test",
			},
		},
	}
	err := kikClient.BroadcastMessage(msgs)
	if err != nil {
		t.Errorf("Error while trying to broadcast a message. %v.", err)
	}
}
