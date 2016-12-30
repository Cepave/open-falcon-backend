package http

import (
	"encoding/json"
	"github.com/Cepave/open-falcon-backend/modules/agent/g"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strconv"
)

type tailReq struct {
	MaxLineNumber int `json:"MaxLineNumber"`
}

func configTailRoutes() {
	http.HandleFunc("/v1/tail", func(w http.ResponseWriter, r *http.Request) {
		if !g.IsTrustable(r.RemoteAddr) {
			w.Write([]byte("no privilege"))
			return
		}

		var reqData tailReq
		reqData.MaxLineNumber = -1
		err := json.NewDecoder(r.Body).Decode(&reqData)
		switch {
		case err == io.EOF:
			// empty body
			reqData.MaxLineNumber = 100
		case err != nil:
			// parsing error
			http.Error(w, "connot decode body", http.StatusBadRequest)
			return
		}
		if reqData.MaxLineNumber < 0 {
			http.Error(w, "Need positive integer of MaxLineNumber", http.StatusBadRequest)
			return
		}

		pwd, err := os.Getwd()
		if err != nil {
			http.Error(w, "connot get current directory", http.StatusBadRequest)
			return
		}
		filepath := pwd + "/var/app.log"

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
