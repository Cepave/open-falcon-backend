package fastweb

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/Jeffail/gabs"
	"github.com/emirpasic/gods/sets/hashset"

	"github.com/Cepave/open-falcon-backend/modules/fe/g"
	log "github.com/Sirupsen/logrus"
)

func QueryContact(platformName string) (contactList []Contactor, err error) {
	config := g.Config()
	fcname := config.Api.Name
	fctoken := getFctoken()
	url := config.Api.Contact
	params := map[string]string{
		"fcname":       fcname,
		"fctoken":      fctoken,
		"platform_key": platformName,
	}
	paramstr, _ := json.Marshal(params)
	log.Debugf("contact get url: %s", url)
	log.Debugf("contact get params: %v", params)
	reqPost, _ := http.NewRequest("POST", url, bytes.NewBuffer([]byte(paramstr)))
	reqPost.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(reqPost)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	parsedJSON, err := gabs.ParseJSONBuffer(resp.Body)
	// contactList = parsedJSON.Data().(map[string]interface{})
	if err != nil {
		return
	}
	contacts, err := parsedJSON.Search("result", platformName).Children()
	contactList = []Contactor{}
	for _, con := range contacts {
		contact := Contactor{
			con.Search("cell").Data().(string),
			con.Search("email").Data().(string),
			con.Search("realname").Data().(string),
		}
		contactList = append(contactList, contact)
	}
	return
}

func GetPlatfromContactInfo(platList *hashset.Set) (contactorMap map[string][]Contactor, err error) {
	contactorMap = map[string][]Contactor{}
	for _, name := range platList.Values() {
		sname := name.(string)
		res, err := QueryContact(sname)
		if err != nil {
			continue
		}
		contactorMap[sname] = res
	}
	return
}
