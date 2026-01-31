package api

import (
	"encoding/json"
	"net/http"
	"time"
	"zepp-os-dev-tool/internal/models"
	"zepp-os-dev-tool/internal/storage"
)

func HandlerAddLog(w http.ResponseWriter, r *http.Request) {
	appName := r.PathValue("appName")
	decoder := json.NewDecoder(r.Body)
	dataStruct := &models.ZeppLogs{}
	err := decoder.Decode(dataStruct)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	dataStruct.TimeStamp = time.Now()
	storage.AddLogsForApp(appName, &models.ZeppLogs{
		LogLevel:    dataStruct.LogLevel,
		Description: dataStruct.Description,
		TimeStamp:   dataStruct.TimeStamp,
	})
	w.WriteHeader(http.StatusCreated)
}
