package integration

import (
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/r-kells/go-kik/kik"
)

var (
	kikClient *kik.KikClient
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
	user, err := kikClient.GetUser("rmdkelly")
	if err != nil {
		t.Fatalf("Could not find user: %v ", err)
	}
	t.Logf("Found User: %v ", user)
}
