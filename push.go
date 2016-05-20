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
		log.Fatalln(GetGeneralConfig().Agent.PushURL+", Error on formatting body:,", err)
	}

	postReq, err := http.NewRequest("POST", GetGeneralConfig().Agent.PushURL, bytes.NewBuffer(paramsBody))
	postReq.Header.Set("Content-Type", "application/json; charset=UTF-8")
	postReq.Header.Set("Connection", "close")

	httpClient := &http.Client{}
	postResp, err := httpClient.Do(postReq)
	if err != nil {
		log.Fatalln("Error on pushing to [", GetGeneralConfig().Agent.PushURL, "]:", err)
	}
	defer postResp.Body.Close()
	log.Println("Pushing the HTTP Body...succeeded")
}
