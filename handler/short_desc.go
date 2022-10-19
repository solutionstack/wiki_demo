package handler

import (
	"fmt"
	logger "github.com/rs/zerolog"
	"io"
	"net/http"
	"os"
	"test/service"
)

var svc service.Service
var Log = logger.New(os.Stdout)

const searchQueryParam = "query"

func init() {
	svc = service.New()
}

func GetDescription(w http.ResponseWriter, r *http.Request) {
	param := r.URL.Query().Get(searchQueryParam)
	if len(param) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		Log.Debug().Caller().Msg("query parameter of invalid length supplied")
		io.WriteString(w, `{"message": "invalid query length in request"}`)
		return
	}

	val, err := svc.GetWikiShortDesc(param)

	if err != nil {
		if err.Error() == "article not found" {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		Log.Debug().Caller().Err(err).Msg("an error occurred during request processing ")

		io.WriteString(w, fmt.Sprintf(`{"message": "%s"}`, err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	io.WriteString(w, fmt.Sprintf(`{"message": "%s"}`, val))
}

func ErrorHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	io.WriteString(w, "{\"message\": \"bad route\"}")

}
