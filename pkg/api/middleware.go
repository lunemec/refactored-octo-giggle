package api

import (
	"encoding/json"
	"net/http"

	log "github.com/mgutz/logxi/v1"
	"github.com/pkg/errors"
)

// Error is interface which allows us to set specific error code to be sent to user
// along with the error message. More metadata about the error may be easily added.
type Error interface {
	error
	StatusCode() int
}

// errJSON represents JSON error to be sent to user.
type errJSON struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"error"`
}

// handler is regular http.Handler but returns error which is processed using
// errHandler.
type handler func(http.ResponseWriter, *http.Request) error

// panicHandler recovers from runtime panics, logs and returns message to user.
func panicHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error
		defer func() {
			r := recover()
			if r != nil {
				switch t := r.(type) {
				case string:
					err = errors.New(t)
				case error:
					err = t
				default:
					err = errors.New("Unknown error")
				}
				log.Error("Panic recovered", "err", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}()
		handler.ServeHTTP(w, r)
	})
}

// ErrHandler allows us to have http.Handler that can return error which is handled
// here and encoded as JSON struct containing the error message and with correct http
// status code.
func ErrHandler(handler handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := handler(w, r)
		if err != nil {
			var e errJSON

			switch v := err.(type) {
			case Error:
				w.WriteHeader(v.StatusCode())
				e = errJSON{
					StatusCode: v.StatusCode(),
					Message:    v.Error(),
				}
			default:
				w.WriteHeader(http.StatusInternalServerError)
				e = errJSON{
					StatusCode: http.StatusInternalServerError,
					Message:    v.Error(),
				}
			}
			w.Header().Set("Content-Type", "application/json")

			out, err := json.Marshal(e)
			if err != nil {
				log.Error("Error returning JSON error", "err", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			_, err = w.Write(out)
			if err != nil {
				log.Error("Error writing response with JSON error", "err", err)
				return
			}
			return
		}
	})
}
