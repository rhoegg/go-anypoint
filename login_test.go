package go_anypoint

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestLogin(t *testing.T) {
	setup()
	defer teardown()

	loginRequest := &LoginRequest{
		Username: "max",
		Password: "mule",
	}

	mux.HandleFunc("/accounts/login", func(w http.ResponseWriter, r *http.Request) {
		v := new(LoginRequest)
		err := json.NewDecoder(r.Body).Decode(v)
		if err != nil {
			t.Fatal(err)
		}

		testMethod(t, r, http.MethodPost)

		if !reflect.DeepEqual(v, loginRequest) {
			t.Errorf("Request body = %+v, expected %+v", v, loginRequest)
		}

		response := ` 
		{
			"access_token": "a-b-c-d-e",
			"token_type": "bearer",
			"redirectUrl": "/home/"
		}`

		fmt.Fprint(w, response)
	})

	result, _, err := client.Login(ctx, loginRequest)
	if err != nil {
		t.Errorf("Login returned error: %v", err)
	}

	expected := &LoginResult{AccessToken: "a-b-c-d-e", TokenType: "bearer", RedirectURL: "/home/"}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Login returned %+v, expected %+v", result, expected)
	}
}

func stubLogin() {
	mux.HandleFunc("/accounts/login", func(w http.ResponseWriter, r *http.Request) {
		response := ` 
		{
			"access_token": "a-b-c-d-e",
			"token_type": "bearer",
			"redirectUrl": "/home/"
		}`

		fmt.Fprint(w, response)
	})
}
