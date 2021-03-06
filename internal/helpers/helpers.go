package helpers

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/NastyInitiative/bookings/internal/config"
)

var app *config.AppConfig

func NewHelpers(a *config.AppConfig) {
	app = a
}

func ClientError(w http.ResponseWriter, status int) {
	app.InfoLog.Println("Client error with status of", status)
	http.Error(w, http.StatusText(status), status)
}

func ServerError(w http.ResponseWriter, pointOfFailure string, err error) {
	trace := fmt.Sprintf("%s -- %s\n%s", pointOfFailure, err.Error(), debug.Stack())
	app.ErrorLog.Println(trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}
