package api

import (
	"encoding/json"
	"net/http"
	"time"
	"zepp-os-dev-tool/internal/models"
	"zepp-os-dev-tool/internal/storage"
)

var RefreshFun func(appName string)

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

	if RefreshFun != nil {
		defer RefreshFun(appName)
	}

	val, exist := storage.GetDataForApp(appName)
	if !exist {
		val = []*models.ZeppMemoryStruct{}
	}
	val = append(val, dataStruct)
	storage.SetData(appName, val)
	w.WriteHeader(http.StatusCreated)

}

func HandlerClearData(w http.ResponseWriter, r *http.Request) {
	appName := r.PathValue("appName")
	storage.SetData(appName, []*models.ZeppMemoryStruct{})
	w.WriteHeader(http.StatusCreated)

}
