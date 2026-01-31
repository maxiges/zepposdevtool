package storage

import (
	"sync"
	"zepp-os-dev-tool/internal/models"
	"zepp-os-dev-tool/internal/utils"
)

var (
	storageData     = make(map[string][]*models.ZeppMemoryStruct, 0)
	storageDataLogs = make(map[string][]*models.ZeppLogs, 0)
)

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
	delete(storageData, appName)
	delete(storageDataLogs, appName)

}

func AddDataForApp(appName string) (ZeppMemoryData, bool) {
	data, exist := storageData[appName]
	return data, exist
}

func GetAppList() []string {
	appSet := utils.NewSet[string]()

	for appName := range storageData {
		appSet.Add(appName)
	}
	for appName := range storageDataLogs {
		appSet.Add(appName)
	}
	return appSet.ToList()
}

//------------

func GetLogsForApp(appName string) ([]*models.ZeppLogs, bool) {
	data, exist := storageDataLogs[appName]
	return data, exist
}
func AddLogsForApp(appName string, log *models.ZeppLogs) {
	data, exist := storageDataLogs[appName]
	if !exist {
		data = []*models.ZeppLogs{}
	}
	lock.Lock()
	defer lock.Unlock()
	data = append(data, log)
	storageDataLogs[appName] = data

}
