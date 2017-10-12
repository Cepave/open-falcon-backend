package client

import (
	"fmt"
	sjson "github.com/bitly/go-simplejson"
	"github.com/dghubble/sling"
)

type SlintExt sling.Sling

func ToSlintExt(slingObject *sling.Sling) *SlintExt {
	return (*SlintExt)(slingObject)
}
func (c *SlintExt) DoReceive(expectedStatus int, successV interface{}) error {
	slingClient := (*sling.Sling)(c)

	jsonError := sjson.New()

	resp, err := slingClient.Receive(successV, jsonError)
	if err != nil {
		return fmt.Errorf("HTTP request error: %v.", err)
	}

	if resp.StatusCode != expectedStatus {
		jsonResp, marshalJsonError := jsonError.MarshalJSON()
		if marshalJsonError != nil {
			return fmt.Errorf("Marshal JSON body of response error: %v", marshalJsonError)
		}

		return fmt.Errorf(
			"Status: %d(needing [%d]). Json Response: %s",
			resp.StatusCode, expectedStatus, string(jsonResp),
		)
	}

	return nil
}
