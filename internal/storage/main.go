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
func ClearAllDataForApp(appName string) {
	lock.Lock()
	defer lock.Unlock()
	_, exist := storageData[appName]
	if exist {
		storageData[appName] = []*models.ZeppMemoryStruct{}
	}
}

func AddDataForApp(appName string) (ZeppMemoryData, bool) {
	data, exist := storageData[appName]
	return data, exist
}

func GetAppList() []string {
	appNameList := make([]string, 0, len(storageData))
	for appName := range storageData {
		appNameList = append(appNameList, appName)
	}
	return appNameList
}
