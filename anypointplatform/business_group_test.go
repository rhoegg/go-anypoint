package anypointplatform_test

import (
	"encoding/json"
	"fmt"
	. "github.com/petergtz/pegomock"
	"github.com/rhoegg/go-anypoint/anypointplatform"
	"net/http"
	"reflect"
	"testing"
)

func TestBusinessGroupGet(t *testing.T) {
	setupBusinessGroupTest(t)
	defer teardown()

	mux.HandleFunc("/accounts/api/organizations/0-1-2-3-4", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		response := ` 
			{
				"name": "Alpha Group",
				"id": "0-1-2-3-4",
				"createdAt": "2018-07-26T21:16:33.464Z",
				"updatedAt": "2018-10-17T23:09:43.861Z"
			}`

		fmt.Fprint(w, response)
	})

	bg, _, err := client.BusinessGroup.Get(ctx, "0-1-2-3-4")
	if err != nil {
		t.Errorf("BusinessGroup.Get returned error: %v", err)
	}

	expected := &anypointplatform.BusinessGroup{Name: "Alpha Group", ID: "0-1-2-3-4"}
	if !reflect.DeepEqual(bg, expected) {
		t.Errorf("BusinessGroup.Get returned %+v, expected %+v", bg, expected)
	}
}

func TestBusinessGroupCreate(t *testing.T) {
	setupBusinessGroupTest(t)
	defer teardown()

	createRequest := &anypointplatform.BusinessGroupCreateRequest{
		Name:     "Alpha Group",
		OwnerID:  "7-6-5-4-3",
		ParentID: "0-1-2-3-4",
	}

	handleHttp(t, "/accounts/api/organizations", http.MethodPost, func(w http.ResponseWriter, r *http.Request) {
		v := new(anypointplatform.BusinessGroupCreateRequest)
		err := json.NewDecoder(r.Body).Decode(v)
		if err != nil {
			t.Fatal(err)
		}

		response := businessGroupTestResponse("Alpha Group")

		if !reflect.DeepEqual(v, createRequest) {
			t.Errorf("Request body = %+v, expected %+v", v, createRequest)
		}

		fmt.Fprint(w, response)
	})

	bg, _, err := client.BusinessGroup.Create(ctx, createRequest)
	if err != nil {
		t.Errorf("BusinessGroup.Create returned error: %v", err)
	}

	expected := &anypointplatform.BusinessGroup{Name: "Alpha Group", ID: "0-1-2-3-4", ClientID: "00112233445566778899aabbccddeeff"}
	if !reflect.DeepEqual(bg, expected) {
		t.Errorf("BusinessGroup.Create returned %+v, expected %+v", bg, expected)
	}
}

func TestBusinessGroupCreateWithName(t *testing.T) {
	setupBusinessGroupTest(t)
	defer teardown()

	handleHttp(t, "/accounts/api/organizations", http.MethodPost, func(w http.ResponseWriter, r *http.Request) {
		v := new(anypointplatform.BusinessGroupCreateRequest)
		err := json.NewDecoder(r.Body).Decode(v)
		if err != nil {
			t.Fatal(err)
		}

		fmt.Fprint(w, businessGroupTestResponse("Alpha Group"))
	})

	created, err := client.BusinessGroup.CreateWithName(ctx, "Alpha Group")
	if err != nil {
		t.Errorf("BusinessGroup.CreateWithName returned error: %v", err)
	}

	if created.Name != "Alpha Group" {
		t.Errorf("Created BusinessGroup has name %v, expected %v", created.Name, "Alpha Group")
	}
}

func TestBusinessGroupCreate_UsesDefaultOwnerIDFromProfile(t *testing.T) {
	setupBusinessGroupTest(t)
	defer teardown()

	handleHttp(t, "/accounts/api/organizations", http.MethodPost, func(w http.ResponseWriter, r *http.Request) {
		v := new(anypointplatform.BusinessGroupCreateRequest)
		err := json.NewDecoder(r.Body).Decode(v)
		if err != nil {
			t.Fatal(err)
		}

		if v.OwnerID != "5-5-5-5-5" {
			t.Errorf("Created BusinessGroup has OwnerID %s, expected %s", v.OwnerID, "5-5-5-5-5")
		}
		fmt.Fprint(w, businessGroupTestResponse("Alpha Group"))
	})

	When(client.Profile.Get(ctx)).ThenReturn(&anypointplatform.Profile{ID: "5-5-5-5-5"}, nil, nil)

	createRequest := &anypointplatform.BusinessGroupCreateRequest{
		Name:     "Alpha Group",
		ParentID: "0-1-2-3-4",
	}
	_, _, err := client.BusinessGroup.Create(ctx, createRequest)
	if err != nil {
		t.Errorf("BusinessGroup.Create returned error: %v", err)
	}
}

func TestBusinessGroupCreate_UsesDefaultParentIDFromOrganizationInProfile(t *testing.T) {
	setupBusinessGroupTest(t)
	defer teardown()

	handleHttp(t, "/accounts/api/organizations", http.MethodPost, func(w http.ResponseWriter, r *http.Request) {
		v := new(anypointplatform.BusinessGroupCreateRequest)
		err := json.NewDecoder(r.Body).Decode(v)
		if err != nil {
			t.Fatal(err)
		}

		if v.ParentID != "6-6-6-6-6" {
			t.Errorf("Created BusinessGroup has ParentID %s, expected %s", v.ParentID, "6-6-6-6-6")
		}
		fmt.Fprint(w, businessGroupTestResponse("Alpha Group"))
	})

	When(client.Profile.Get(ctx)).ThenReturn(&anypointplatform.Profile{ID: "MOCK_PROFILE_ID", OrganizationID: "6-6-6-6-6"}, nil, nil)

	createRequest := &anypointplatform.BusinessGroupCreateRequest{
		Name:    "Alpha Group",
		OwnerID: "0-1-2-3-4",
	}
	_, _, err := client.BusinessGroup.Create(ctx, createRequest)
	if err != nil {
		t.Errorf("BusinessGroup.Create returned error: %v", err)
	}
}

func TestBusinessGroupDelete(t *testing.T) {
	setupBusinessGroupTest(t)
	defer teardown()

	deleted := false

	handleHttp(t, "/accounts/api/organizations/0-1-2-3-4", http.MethodDelete, func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "{}")
		deleted = true
	})

	_, err := client.BusinessGroup.Delete(ctx, "0-1-2-3-4")

	if err != nil {
		t.Errorf("BusinessGroup.Delete returned error: %+v", err)
	}

	if !deleted {
		t.Errorf("BusinessGroup.Delete did not invoke the API correctly.")
	}
}

func setupBusinessGroupTest(t *testing.T) {
	setup()
	stubLogin()
	stubGetProfile()
	RegisterMockTestingT(t)
}

func stubGetProfile() {
	client.Profile = NewMockProfileService()
	When(client.Profile.Get(ctx)).ThenReturn(&anypointplatform.Profile{ID: "MOCK-PROFILE-ID"}, nil, nil)
}

func businessGroupTestResponse(name string) string {
	return fmt.Sprintf(`{
			"name": "%v",
			"id": "0-1-2-3-4",
			"createdAt": "2018-07-26T21:16:33.464Z",
			"updatedAt": "2018-10-17T23:09:43.861Z",
			"ownerId": "7-6-5-4-3",
			"clientId": "00112233445566778899aabbccddeeff",
			"domain": null,
			"idprovider_id": "mulesoft",
			"isFederated": false,
			"parentOrganizationIds": ["0-1-2-3-4"],
			"subOrganizationIds": [],
			"tenantOrganizationIds": [],
			"isMaster": false
		}`, name)
}
