package handlers;

import (
	"net/http"
)

// MakeSecretHandler makes a handler for Create/List/Delete/Update of
//secrets in the Rancher API
func MakeSecretHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		w.WriteHeader(http.StatusNotImplemented)
		w.Write([]byte("Not implemented\n"))
	}
}