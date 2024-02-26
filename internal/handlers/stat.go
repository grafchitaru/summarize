package handlers

import (
	"encoding/json"
	"github.com/grafchitaru/summarize/internal/config"
	"github.com/grafchitaru/summarize/internal/middlewares/auth"
	"net/http"
)

func Stat(ctx config.HandlerContext, res http.ResponseWriter, req *http.Request) {
	userID, err := auth.GetUserID(req, ctx.Config.SecretKey)
	if err != nil {
		http.Error(res, err.Error(), http.StatusUnauthorized)
		return
	}

	result, err := ctx.Repos.GetStat(userID)
	if err != nil {
		http.Error(res, err.Error(), http.StatusNotFound)
		return
	}

	data, err := json.Marshal(result)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusOK)
	res.Write(data)
}
