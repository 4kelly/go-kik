// Package integration runs tests against the real Kik API.
// Generally they will only fail if there is an error in the payloads being sent.
// This is useful to verify the Kik API hasn't introduced breaking request / response types.
package integration

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
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

func TestGetUser_HappyPath(t *testing.T) {

	_, err := kikClient.GetUser(testUserName)
	if err != nil {
		t.Errorf("Could not find user: %v ", err)
	}
}

// func TestSendMessage_HappyPath(t *testing.T) {

// 	err := kikClient.SendMessage(allMessageTypesTestData)
// 	if err != nil {
// 		t.Errorf("Error while trying to send a message. %v.", err)
// 	}
// }

// func TestBroadcastMessage_HappyPath(t *testing.T) {

// 	err := kikClient.BroadcastMessage(allMessageTypesTestData)
// 	if err != nil {
// 		t.Errorf("Error while trying to broadcast a message. %v.", err)
// 	}
// }

// TestConfig_HappyPath Sets then gets Kik bot configuration.
func TestConfig_HappyPath(t *testing.T) {
	keyboard := &kik.SuggestedResponseKeyboard{
		Type: "suggested",
		Responses: []interface{}{
			kik.KeyboardTextResponse{
				Type: "text",
				Body: "StaticKeyboardTest",
			},
		},
	}
	wantConfig := &kik.Configuration{
		Webhook: "http://example.com",
		Features: &kik.Features{
			ManuallySendReadReceipts: false,
			ReceiveReadReceipts:      false,
			ReceiveDeliveryReceipts:  false,
			ReceiveIsTyping:          false,
		},
		StaticKeyboard: keyboard,
	}
	err := kikClient.SetConfiguration(wantConfig)
	if err != nil {
		t.Errorf("Error while trying to set configuration. %v.", err)
	}
	gotConfig, err := kikClient.GetConfiguration()
	if err != nil {
		t.Errorf("Error while trying to get configuration. %v.", err)
	}

	if !cmp.Equal(gotConfig, wantConfig) {
		t.Errorf("SetConfiguration() = %v; want %v", gotConfig, wantConfig)
	}

}

// TestCreateCode_HappyPath Creates then gets image for a Kik Scan Code.
func TestCreateCode_HappyPath(t *testing.T) {
	scanCodeData := &kik.ScanData{
		Data: "Kik Scan Code Example Data!",
	}
	code, err := kikClient.CreateCode(scanCodeData)
	if err != nil {
		t.Errorf("Error while trying create a Kik scan code. %v.", err)
	}

	url := fmt.Sprintf("https://api.kik.com/v1/code/%s?c=1", code.Id)
	codeID, err := http.Get(url)
	if err != nil {
		t.Errorf("Error while trying to get scan code image. %v.", err)
	}

	defer codeID.Body.Close()
	body, _ := ioutil.ReadAll(codeID.Body)
	contentType := http.DetectContentType(body)

	if contentType != "image/png" {
		t.Errorf("Returned data must be an image.")
	}
}

/*
Test data
*/

// Who we send test messages to
var testUserName = "rmdkelly"

// Contains an example of all the keyboard response types.
var allKeyboardTypesTestData = []kik.SuggestedResponseKeyboard{
	{Type: "suggested",
		Responses: []interface{}{
			kik.KeyboardPictureResponse{
				Type:     "picture",
				PicUrl:   "https://i.imgur.com/8rqLdgy.png",
				Metadata: "picture1",
			},
			kik.KeyboardFriendPickerResponse{
				Type:        "friend-picker",
				Body:        "Test",
				Min:         0,
				Max:         2,
				Preselected: []string{"cacolvil"},
			},
			kik.KeyboardTextResponse{
				Type: "text",
				Body: "KeyboardTextResponse",
			},
		},
	},
}

// Contains an example of all the message types
var allMessageTypesTestData = []kik.Message{
	kik.TextMessage{
		SendMessage: kik.SendMessage{
			To:   testUserName,
			Type: "text",
		},
		Body: "Test_TestMessage",
	},
	kik.PictureMessage{
		SendMessage: kik.SendMessage{
			To:   testUserName,
			Type: "picture",
		},
		PicUrl: "https://i.imgur.com/TsoLODG.png",
		Attribution: &kik.Attribution{
			Name: "Attribution Test",
		},
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
		Attribution: &kik.Attribution{
			Name: "Attribution Test",
		},
	},
	kik.VideoMessage{
		SendMessage: kik.SendMessage{
			To:        testUserName,
			Type:      "video",
			Delay:     1,
			Keyboards: allKeyboardTypesTestData,
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
