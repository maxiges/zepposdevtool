package api

import (
	"encoding/json"
	"net/http"
	"zepp-os-dev-tool/internal/storage"
)

type AppListResponse struct {
	Apps []string `json:"apps"`
}

func HandlerGetAppList(w http.ResponseWriter, r *http.Request) {

	appNames := storage.GetAppList()
	data := AppListResponse{
		Apps: appNames,
	}
	j, err := json.Marshal(data)

	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(j))
}
