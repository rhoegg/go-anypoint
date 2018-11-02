package go_anypoint

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestBusinessGroupGet(t *testing.T) {
	setup()
	stubLogin()
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

	expected := &BusinessGroup{Name: "Alpha Group"}
	if !reflect.DeepEqual(bg, expected) {
		t.Errorf("BusinessGroup.Get returned %+v, expected %+v", bg, expected)
	}
}

func TestBusinessGroupCreate(t *testing.T) {
	setup()
	stubLogin()
	defer teardown()

	createRequest := &BusinessGroupCreateRequest{
		Name:                 "Alpha Group",
		OwnerID:              "7-6-5-4-3",
		ParentOrganizationID: "0-1-2-3-4",
	}

	mux.HandleFunc("/accounts/api/organizations", func(w http.ResponseWriter, r *http.Request) {
		v := new(BusinessGroupCreateRequest)
		err := json.NewDecoder(r.Body).Decode(v)
		if err != nil {
			t.Fatal(err)
		}

		testMethod(t, r, http.MethodPost)
		response := ` 
		{
			"name": "Alpha Group",
			"id": "0-1-2-3-4",
			"createdAt": "2018-07-26T21:16:33.464Z",
			"updatedAt": "2018-10-17T23:09:43.861Z",
			"ownerId": "7-6-5-4-3",
			"clientId": "00112233445566778899aabbccddeeff",
			"domain": null,
			"idprovider_id": "mulesoft",
			"isFederated": false,
			"parentOrganizationIds": ["0-1-2-3-4-5"],
			"subOrganizationIds": [],
			"tenantOrganizationIds": [],
			"isMaster": false
		}`

		if !reflect.DeepEqual(v, createRequest) {
			t.Errorf("Request body = %+v, expected %+v", v, createRequest)
		}

		fmt.Fprint(w, response)
	})

	bg, _, err := client.BusinessGroup.Create(ctx, createRequest)
	if err != nil {
		t.Errorf("BusinessGroup.Create returned error: %v", err)
	}

	expected := &BusinessGroup{Name: "Alpha Group"}
	if !reflect.DeepEqual(bg, expected) {
		t.Errorf("BusinessGroup.Create returned %+v, expected %+v", bg, expected)
	}
}
