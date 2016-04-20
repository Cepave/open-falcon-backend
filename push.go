package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func Push(params []ParamToAgent) {
	paramsBody, err := json.Marshal(params)
	if err != nil {
		fmt.Println(GetGeneralConfig().AgentPushURL+", format body error,", err)
	}

	postReq, err := http.NewRequest("POST", GetGeneralConfig().AgentPushURL, bytes.NewBuffer(paramsBody))
	postReq.Header.Set("Content-Type", "application/json; charset=UTF-8")
	postReq.Header.Set("Connection", "close")

	httpClient := &http.Client{}
	postResp, err := httpClient.Do(postReq)
	if err != nil {
		log.Print("Error in pushing to ", GetGeneralConfig().AgentPushURL, ": ", err, "\n")
	}
	defer postResp.Body.Close()
	fmt.Println("Body:", string(paramsBody))
}
