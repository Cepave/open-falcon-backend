package http

import (
	"encoding/json"
	"net/http"
	"os"
	"os/exec"
	"strconv"
)

type tailReq struct {
	MaxLineNumber int `json:"MaxLineNumber"`
}

func configTailRoutes() {
	http.HandleFunc("/v1/tail", func(w http.ResponseWriter, req *http.Request) {
		var reqData tailReq
		if req.ContentLength == 0 {
			reqData.MaxLineNumber = 10
		} else {
			decoder := json.NewDecoder(req.Body)
			err := decoder.Decode(&reqData)
			if err != nil {
				http.Error(w, "connot decode body", http.StatusBadRequest)
				return
			}
		}

		pwd, err := os.Getwd()
		if err != nil {
			http.Error(w, "connot get current directory", http.StatusBadRequest)
			return
		}
		filepath := pwd + "/var/app.log"
		//Debug use only
		//w.Write([]byte("log file is: " + filepath + "\n"))

		cmd := exec.Command("tail", "-n", strconv.Itoa(reqData.MaxLineNumber), filepath)
		out, err := cmd.CombinedOutput()
		if err != nil {
			w.Write([]byte("tail command failed\n"))
			w.Write([]byte(out))
			w.Write([]byte("line number is: " + strconv.Itoa(reqData.MaxLineNumber)))
			return
		}
		w.Write([]byte(out))
	})
}
