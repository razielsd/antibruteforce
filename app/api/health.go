package api

import (
	"net/http"
)

func (a *AbfAPI) handlerHealthProbe(w http.ResponseWriter, r *http.Request) {
	a.sendResult(w, NewSuccessOK())
}
