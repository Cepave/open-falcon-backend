package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

func Push(params []ParamToAgent) {
	paramsBody, err := json.Marshal(params)
	if err != nil {
		log.Fatalln(GetGeneralConfig().AgentPushURL+", Error on formatting body:,", err)
	}

	postReq, err := http.NewRequest("POST", GetGeneralConfig().AgentPushURL, bytes.NewBuffer(paramsBody))
	postReq.Header.Set("Content-Type", "application/json; charset=UTF-8")
	postReq.Header.Set("Connection", "close")

	httpClient := &http.Client{}
	postResp, err := httpClient.Do(postReq)
	if err != nil {
		log.Fatalln("Error on pushing to [", GetGeneralConfig().AgentPushURL, "]:", err)
	}
	defer postResp.Body.Close()
	log.Println("Pushing the HTTP Body...succeeded")
}
