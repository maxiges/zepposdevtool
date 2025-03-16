package api

import (
	"encoding/json"
	"my-chart-app/internal/models"
	"my-chart-app/internal/storage"
	"net/http"
	"time"
)

func HandlerAddData(w http.ResponseWriter, r *http.Request) {
	appName := r.PathValue("appName")
	decoder := json.NewDecoder(r.Body)
	dataStruct := &models.ZeppMemoryStruct{}
	err := decoder.Decode(dataStruct)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	dataStruct.TimeStamp = time.Now()

	val, exist := storage.GetDataForApp(appName)
	if !exist {
		val = []*models.ZeppMemoryStruct{}
	}
	val = append(val, dataStruct)
	storage.SetData(appName, val)
	w.WriteHeader(http.StatusCreated)

}
