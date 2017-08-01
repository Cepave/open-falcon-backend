package http

import (
	"fmt"
	"net/http"

	"github.com/Cepave/open-falcon-backend/modules/query/g"
	"github.com/juju/errors"
	log "github.com/sirupsen/logrus"
	"github.com/toolkits/file"
)

func configCommonRoutes() {
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok\n"))
	})

	http.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(fmt.Sprintf("%s\n", g.VERSION)))
	})

	http.HandleFunc("/workdir", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(fmt.Sprintf("%s\n", file.SelfDir())))
	})

	http.HandleFunc("/config", func(w http.ResponseWriter, r *http.Request) {
		RenderDataJson(w, g.Config())
	})

}

func setErrorMessage(message string, result map[string]interface{}) {
	jujuErr := errors.NewErr(message)
	jujuErr.SetLocation(1)

	putError(result, &jujuErr)
}

func setError(err error, result map[string]interface{}) {
	jujuErr := errors.NewErr("%v", err)
	jujuErr.SetLocation(1)

	putError(result, &jujuErr)
}

func putError(container map[string]interface{}, err error) {
	log.Errorf("Error has occurred: %s", errors.ErrorStack(err))

	container["error"] = append(container["error"].([]string), err.Error())
}
