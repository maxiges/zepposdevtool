package storage

import (
	"my-chart-app/internal/models"
	"sync"
)

var storageData = make(map[string][]*models.ZeppMemoryStruct, 0)
var lock = sync.Mutex{}

type ZeppMemoryData []*models.ZeppMemoryStruct

func init() {

}

func GetDataForApp(appName string) (ZeppMemoryData, bool) {
	data, exist := storageData[appName]
	return data, exist
}
func SetData(appName string, data []*models.ZeppMemoryStruct) {
	lock.Lock()
	defer lock.Unlock()
	storageData[appName] = data
}

func AddDataForApp(appName string) (ZeppMemoryData, bool) {
	data, exist := storageData[appName]
	return data, exist
}
