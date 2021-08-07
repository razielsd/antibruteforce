package api

import (
	"net/http"
)

func (a *AbfAPI) ActionHealthProbe(w http.ResponseWriter, r *http.Request) {
	a.sendResult(w, NewSuccessOK())
}
