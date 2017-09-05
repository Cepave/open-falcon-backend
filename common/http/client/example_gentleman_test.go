package client_test

import (
	"errors"
	"fmt"
	"net/http"

	gt "gopkg.in/h2non/gentleman.v2"
	"gopkg.in/h2non/gock.v1"

	"github.com/Cepave/open-falcon-backend/common/http/client"
)

// Send request and confirm status code.
func ExampleGentlemanRequest_sendAndStatusMatch() {
	defer gock.Off()

	/**
	 * Mock service
	 */
	gock.New("http://example-1.gock/success").
		Reply(http.StatusOK)
	// :~)

	req := newClientByGock().
		URL("http://example-1.gock/").
		Get().
		AddPath("success")

	_, err := client.ToGentlemanReq(req).SendAndStatusMatch(http.StatusOK)
	if err != nil {
		fmt.Printf("Err: %v\n", err)
		return
	}

	fmt.Println("Success")

	// Output:
	// Success
}

// Send request and confirm by customized matcher.
func ExampleGentlemanRequest_sendAndMatch() {
	defer gock.Off()

	/**
	 * Mock service
	 */
	gock.New("http://example-2.gock/success").
		Reply(http.StatusOK)
	// :~)

	req := newClientByGock().
		URL("http://example-2.gock/").
		Get().
		AddPath("success")

	_, err := client.ToGentlemanReq(req).SendAndMatch(
		func(resp *gt.Response) error {
			if resp.StatusCode != http.StatusOK {
				return errors.New("Not Success")
			}

			return nil
		},
	)
	if err != nil {
		fmt.Printf("Err: %v\n", err)
		return
	}

	fmt.Println("Success")

	// Output:
	// Success
}

// Gets JSON object from response object(panic if some error has occurred)
func ExampleGentlemanResponse_mustGetJson() {
	defer gock.Off()

	/**
	 * Mock service
	 */
	gock.New("http://example-3.gock/json-1").
		Reply(http.StatusOK).
		JSON(map[string]interface{}{
			"name": "King",
			"age":  33,
		})
	// :~)

	req := newClientByGock().
		URL("http://example-3.gock/").
		Get().
		AddPath("json-1")

	resp := client.ToGentlemanReq(req).SendAndStatusMustMatch(http.StatusOK)
	json := client.ToGentlemanResp(resp).MustGetJson()

	fmt.Printf("%s %d", json.Get("name").MustString(), json.Get("age").MustInt())

	// Output:
	// King 33
}

// Binds JSON object from response object(panic if some error has occurred)
func ExampleGentlemanResponse_mustBindJson() {
	defer gock.Off()

	/**
	 * Mock service
	 */
	gock.New("http://example-4.gock/json-2").
		Reply(http.StatusOK).
		JSON(map[string]interface{}{
			"name": "Jon Snow",
			"age":  18,
		})
	// :~)

	req := newClientByGock().
		URL("http://example-4.gock/").
		Get().
		AddPath("json-2")

	jsonBody := &struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}{}

	resp := client.ToGentlemanReq(req).SendAndStatusMustMatch(http.StatusOK)
	client.ToGentlemanResp(resp).MustBindJson(&jsonBody)

	fmt.Printf("%s %d", jsonBody.Name, jsonBody.Age)

	// Output:
	// Jon Snow 18
}
