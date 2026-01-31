package api

import (
	"encoding/json"
	"net/http"
	"zepp-os-dev-tool/internal/models"
	"zepp-os-dev-tool/internal/storage"
)

type DataResponse struct {
	Data []*DataResponseData `json:"data"`
}
type DataResponseData struct {
	Memory        models.Memory `json:"memory"`
	Description   string        `json:"description,omitempty"`
	TimeStampUnix int64         `json:"timestamp,omitempty"`
}

func HandlerGetData(w http.ResponseWriter, r *http.Request) {
	appName := r.PathValue("appName")
	data, exist := storage.GetDataForApp(appName)
	if !exist || len(data) == 0 {
		http.Error(w, "No data found for app: "+appName, http.StatusNotFound)
		return
	}

	respDataList := make([]*DataResponseData, 0, len(data))
	for _, singleData := range data {
		respDataList = append(respDataList, &DataResponseData{
			Memory:        singleData.Memory,
			Description:   singleData.Description,
			TimeStampUnix: singleData.TimeStamp.Unix(),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	resp := DataResponse{
		Data: respDataList}

	j, _ := json.Marshal(resp)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(j))

}

type LogsResponse struct {
	Data []*LogResponseData `json:"data"`
}
type LogResponseData struct {
	LogLevel    string `json:"log_level,omitempty"`
	Description string `json:"description,omitempty"`
	TimeStamp   int64  `json:"timestamp"`
}

func HandlerGetLogs(w http.ResponseWriter, r *http.Request) {
	appName := r.PathValue("appName")
	data, exist := storage.GetLogsForApp(appName)
	if !exist || len(data) == 0 {
		http.Error(w, "No data found for app: "+appName, http.StatusNotFound)
		return
	}

	respDataList := make([]*LogResponseData, 0, len(data))
	for _, singleData := range data {
		respDataList = append(respDataList, &LogResponseData{
			LogLevel:    singleData.LogLevel.String(),
			Description: singleData.Description,
			TimeStamp:   singleData.TimeStamp.Unix(),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	resp := LogsResponse{
		Data: respDataList,
	}
	j, _ := json.Marshal(resp)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(j))

}
