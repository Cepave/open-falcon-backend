package http

import (
	"bytes"
	"github.com/Cepave/query/g"
	"io/ioutil"
	"log"
	"net/http"
)

/**
 * @function name:	func graphInfo(rw http.ResponseWriter, req *http.Request)
 * @description:	This function sends a POST request in JSON format.
 * @related issues:	OWL-171
 * @param:			rw http.ResponseWriter
 * @param:			req *http.Request
 * @return:			void
 * @author:			Don Hsieh
 * @since:			11/12/2015
 * @last modified: 	11/13/2015
 * @called by:		func graphInfo(rw http.ResponseWriter, req *http.Request)
 *					func graphHistory(rw http.ResponseWriter, req *http.Request)
 */
func postByJson(rw http.ResponseWriter, req *http.Request, url string) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(req.Body)

	log.Println(buf.Len())
	s := buf.String()
	log.Println("s =", s)

	reqPost, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(s)))
	if err != nil {
		log.Println("Error =", err.Error())
	}
	reqPost.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(reqPost)
	if err != nil {
		log.Println("Error =", err.Error())
	}
	defer resp.Body.Close()

	log.Println("response Status:", resp.Status)	// 200 OK   TypeOf(resp.Status): string
	log.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	log.Println("response Body:", string(body))
	rw.Header().Set("Content-Type", "application/json; charset=UTF-8")
	rw.Write(body)
}

/**
 * @function name:	func graphInfo(rw http.ResponseWriter, req *http.Request)
 * @description:	This function handles /graph/info API request.
 * @related issues:	OWL-171
 * @param:			rw http.ResponseWriter
 * @param:			req *http.Request
 * @return:			void
 * @author:			Don Hsieh
 * @since:			11/12/2015
 * @last modified: 	11/13/2015
 * @called by:		func configApiRoutes()
 */
func graphInfo(rw http.ResponseWriter, req *http.Request) {
	log.Println("func graphInfo(rw http.ResponseWriter, req *http.Request)")
	url := g.Config().Api.Graph + "/graph/info"
	log.Println("url =", url)
	postByJson(rw, req, url)
}

/**
 * @function name:	func graphInfo(rw http.ResponseWriter, req *http.Request)
 * @description:	This function handles /graph/history API request.
 * @related issues:	OWL-171
 * @param:			rw http.ResponseWriter
 * @param:			req *http.Request
 * @return:			void
 * @author:			Don Hsieh
 * @since:			11/12/2015
 * @last modified: 	11/13/2015
 * @called by:		func configApiRoutes()
 */
func graphHistory(rw http.ResponseWriter, req *http.Request) {
	log.Println("func graphHistory(rw http.ResponseWriter, req *http.Request)")
	url := g.Config().Api.Graph + "/graph/history"
	log.Println("url =", url)
	postByJson(rw, req, url)
}

/**
 * @function name:	func dashboardEndpoints(rw http.ResponseWriter, req *http.Request)
 * @description:	This function handles /api/endpoints API request.
 * @related issues:	OWL-171
 * @param:			rw http.ResponseWriter
 * @param:			req *http.Request
 * @return:			void
 * @author:			Don Hsieh
 * @since:			11/12/2015
 * @last modified: 	11/13/2015
 * @called by:		func configApiRoutes()
 */
func dashboardEndpoints(rw http.ResponseWriter, req *http.Request) {
	log.Println("func dashboardEndpoints(rw http.ResponseWriter, req *http.Request)")
	url := g.Config().Api.Dashboard + req.URL.RequestURI()
	log.Println("url =", url)

	reqGet, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println("Error =", err.Error())
	}
	log.Println("reqGet =", reqGet)

	client := &http.Client{}
	resp, err := client.Do(reqGet)
	if err != nil {
		log.Println("Error =", err.Error())
	}
	defer resp.Body.Close()

	log.Println("response Status:", resp.Status)	// 200 OK   TypeOf(resp.Status): string
	log.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	log.Println("response Body:", string(body))
	rw.Header().Set("Content-Type", "application/json; charset=UTF-8")
	rw.Write(body)

}

/**
 * @function name:	func postByForm(rw http.ResponseWriter, req *http.Request)
 * @description:	This function sends a POST request in Form format.
 * @related issues:	OWL-171
 * @param:			rw http.ResponseWriter
 * @param:			req *http.Request
 * @return:			void
 * @author:			Don Hsieh
 * @since:			11/12/2015
 * @last modified: 	11/13/2015
 * @called by:		func dashboardCounters(rw http.ResponseWriter, req *http.Request)
 *					func dashboardChart(rw http.ResponseWriter, req *http.Request)
 */
func postByForm(rw http.ResponseWriter, req *http.Request, url string) {
	req.ParseForm()
	client := &http.Client{}
	resp, err := client.PostForm(url, req.PostForm)
	if err != nil {
		log.Println("Error =", err.Error())
	}
	defer resp.Body.Close()

	log.Println("response Status:", resp.Status)	// 200 OK   TypeOf(resp.Status): string
	log.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	log.Println("response Body:", string(body))
	rw.Header().Set("Content-Type", "application/json; charset=UTF-8")
	rw.Write(body)
}

/**
 * @function name:	func dashboardCounters(rw http.ResponseWriter, req *http.Request)
 * @description:	This function handles /api/counters API request.
 * @related issues:	OWL-171
 * @param:			rw http.ResponseWriter
 * @param:			req *http.Request
 * @return:			void
 * @author:			Don Hsieh
 * @since:			11/13/2015
 * @last modified: 	11/13/2015
 * @called by:		func configApiRoutes()
 */
func dashboardCounters(rw http.ResponseWriter, req *http.Request) {
	log.Println("func dashboardCounters(rw http.ResponseWriter, req *http.Request)")
	url := g.Config().Api.Dashboard + "/api/counters"
	log.Println("url =", url)
	postByForm(rw, req, url)
}

/**
 * @function name:	func dashboardChart(rw http.ResponseWriter, req *http.Request)
 * @description:	This function handles /api/chart API request.
 * @related issues:	OWL-171
 * @param:			rw http.ResponseWriter
 * @param:			req *http.Request
 * @return:			void
 * @author:			Don Hsieh
 * @since:			11/13/2015
 * @last modified: 	11/13/2015
 * @called by:		func configApiRoutes()
 */
func dashboardChart(rw http.ResponseWriter, req *http.Request) {
	log.Println("func dashboardChart(rw http.ResponseWriter, req *http.Request)")
	url := g.Config().Api.Dashboard + "/chart"
	log.Println("url =", url)
	postByForm(rw, req, url)
}

/**
 * @function name:	func configApiRoutes()
 * @description:	This function handles API requests.
 * @related issues:	OWL-171
 * @param:			void
 * @return:			void
 * @author:			Don Hsieh
 * @since:			11/12/2015
 * @last modified: 	11/13/2015
 * @called by:		func Start()
 *					 in http/http.go
 */
func configApiRoutes() {
	http.HandleFunc("/api/info", graphInfo)
	http.HandleFunc("/api/history", graphHistory)
	http.HandleFunc("/api/endpoints", dashboardEndpoints)
	http.HandleFunc("/api/counters", dashboardCounters)
	http.HandleFunc("/api/chart", dashboardChart)
}
