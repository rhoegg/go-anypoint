package anypointplatform_test

import (
	"fmt"
	"github.com/rhoegg/go-anypoint/anypointplatform"
	"net/http"
	"reflect"
	"testing"
)

func TestProfileGet(t *testing.T) {
	setupProfile(t)
	defer teardown()

	p, _, err := client.Profile.Get(ctx)

	if err != nil {
		t.Errorf("Profile.Get returned error: %v", err)
	}

	expected := &anypointplatform.Profile{ID: "4-5-6-7-8", OrganizationID: "0-9-8-7-6"}

	if !reflect.DeepEqual(p, expected) {
		t.Errorf("Profile.Get returned %v, expected %v", p, expected)
	}
}

func setupProfile(t *testing.T) {
	setup()
	stubLogin()
	mux.HandleFunc("/accounts/api/profile", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		response := `
	{
	    "id": "4-5-6-7-8",
	    "createdAt": "2016-12-02T07:29:00.925Z",
	    "updatedAt": "2018-12-09T09:47:26.759Z",
	    "organizationId": "0-9-8-7-6",
	    "firstName": "Cloud03",
	    "lastName": "Instructor",
	    "email": "Cloud03Instructor@gmail.com"
	}`
		fmt.Fprint(w, response)
	})
}
