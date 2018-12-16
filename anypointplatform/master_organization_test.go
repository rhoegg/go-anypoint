package anypointplatform_test

import (
	"fmt"
	"github.com/rhoegg/go-anypoint/anypointplatform"
	. "github.com/petergtz/pegomock"
	"net/http"
	"reflect"
	"testing"
)

func TestMasterOrganizationGet(t *testing.T) {
	setup()
	stubLogin()
	RegisterMockTestingT(t)
	defer teardown()

	const expectedOrgId = "5-6-7-8-9"

	client.Profile = NewMockProfileService()
	When(client.Profile.Get(ctx)).ThenReturn(&anypointplatform.Profile{ID: "MOCK-PROFILE-ID", OrganizationID: expectedOrgId}, nil, nil)

	mux.HandleFunc(fmt.Sprintf("/accounts/api/organizations/%s", expectedOrgId), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		response := fmt.Sprintf(` 
			{
				"name": "Test Organization",
				"id": "%s",
				"createdAt": "2018-07-26T21:16:33.464Z",
				"updatedAt": "2018-10-17T23:09:43.861Z"
			}`, expectedOrgId)

		fmt.Fprint(w, response)
	})

	o, _, err := client.MasterOrganization.Get(ctx)
	if err != nil {
		t.Errorf("MasterOrganization.Get returned error: %v", err)
	}

	expected := &anypointplatform.MasterOrganization{Name: "Test Organization", ID: expectedOrgId}
	if !reflect.DeepEqual(o, expected) {
		t.Errorf("MasterOrganization.Get returned %+v, expected %+v", o, expected)
	}
}