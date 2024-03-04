package handlers

import (
	"net/http"
)

func (ctx *HandlerContext) Ping(res http.ResponseWriter, req *http.Request) {
	err := ctx.Repos.Ping()
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusOK)
}
