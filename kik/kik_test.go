package kik_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/r-kells/go-kik/kik"
	"github.com/r-kells/go-kik/kiktest"
)

// User is an example kik username used for testing, it can be anything.
const username = "kikteam"

func TestGetUser_HappyPath(t *testing.T) {
	client, mux, teardown := kiktest.TestClient(t)
	defer teardown()

	expectedUser := &kik.User{
		FirstName:              "Ryan",
		LastName:               ".",
		ProfilePicLastModified: 1560526317131,
		ProfilePicUrl:          "https://cdn.kik.com/User/pic/rmdkelly/big",
	}

	mux.HandleFunc(kik.GetUserUrl, func(w http.ResponseWriter, r *http.Request) {
		err := json.NewEncoder(w).Encode(expectedUser)
		if err != nil {
			panic(err)
		}
		fmt.Fprint(w, expectedUser)
	})

	gotUser, err := client.GetUser("Foo")

	if err != nil {
		t.Errorf("GetUser(%s) returned an error = %s; expected no error", username, err)
	}
	if !cmp.Equal(gotUser, expectedUser) {
		t.Errorf("GetUser(%s) = %v; want %v", username, gotUser, expectedUser)
	}
}

func TestGetUser_404(t *testing.T) {
	client, _, teardown := kiktest.TestClient(t)
	defer teardown()

	_, err := client.GetUser(username)

	// TODO custom error types.
	if !strings.Contains(fmt.Sprint(err), "status code != OK") {
		t.Errorf("Expected 404, got %v", err)
	}
}

// This really testing the helper methods.
// Should drop this after explicitly adding tests for helpers.
func TestGetUser_ShouldFailToDecodeUser(t *testing.T) {
	client, mux, teardown := kiktest.TestClient(t)
	defer teardown()

	mux.HandleFunc(kik.GetUserUrl, func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"firstName": true}`)
	})

	_, err := client.GetUser(username)

	if !strings.Contains(
		fmt.Sprint(err),
		"cannot unmarshal bool into Go struct field User.FirstName of type string") {
		t.Errorf("Expected a json decode error, got %v", err)
	}
}
