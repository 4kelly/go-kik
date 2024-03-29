package kik_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/4kelly/go-kik/kik"
	"github.com/4kelly/go-kik/kiktest"
	"github.com/google/go-cmp/cmp"
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
		t.Errorf("GetUser(%s) returned an error = %+v; expected no error", username, err)
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
	if !strings.Contains(fmt.Sprint(err), "404 page not found") {
		t.Errorf("Expected 404, got %v", err)
	}
}

// This really testing the helper methods.
// Should drop this after explicitly adding tests for helpers.
// TODO maybe this test should validate the errors passed to the user of this library.
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

// TODO: can't verify this is working yet, seems like something on Kik side.
// func TestVerifySignature_Valid(t *testing.T) {
// 	client, _, teardown := kiktest.TestClient(t)
// 	defer teardown()

// 	got := client.VerifySignature("AC18D0105C2C257652859322B0499313342C6EB9", []byte("body"))
// 	want := true

// 	if got != want {
// 		t.Errorf("Expected signature validation to be correct.")
// 	}
// }

func TestVerifySignature_Invalid(t *testing.T) {
	client, _, teardown := kiktest.TestClient(t)
	defer teardown()
	got := client.VerifySignature("invalid sig", []byte("body"))
	want := false

	if got != want {
		t.Errorf("Expected signature validation to fail.")
	}
}
