package go_anypoint

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestBusinessGroupGet(t *testing.T) {
	setup()
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