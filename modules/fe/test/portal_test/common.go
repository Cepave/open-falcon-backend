package PortalTest

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/Jeffail/gabs"
	"github.com/parnurzeal/gorequest"
)

var host = "127.0.0.1:1234"
var name = "root"
var password = "password"
var session = ""
var request = gorequest.New()

type Resps struct {
	Body  string
	Code  int
	Error error
}

func DoPost(url string, params string) (myresp Resps) {
	postdata := strings.NewReader(params)
	posturl := fmt.Sprintf("http://%s%s", host, url)
	resp, err := http.Post(posturl, "application/x-www-form-urlencoded", postdata)
	if err == nil {
		body, _ := ioutil.ReadAll(resp.Body)
		myresp.Body = string(body)
	}
	myresp.Code = resp.StatusCode
	myresp.Error = err
	return
}

func GetAuthSessoion() (string, string) {
	resp := DoPost("/api/v1/auth/login", fmt.Sprintf(`name=%s;password=%s`, name, password))
	jsParsed, _ := gabs.ParseJSON([]byte(resp.Body))
	session = jsParsed.Search("data", "sig").Data().(string)
	return name, session
}
