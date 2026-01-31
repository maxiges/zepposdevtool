package models

import "time"

type ZeppMemoryStruct struct {
	Memory      Memory    `json:"memory"`
	Description string    `json:"description,omitempty"`
	TimeStamp   time.Time `json:"timestamp"`
}
type Memory struct {
	App       []memoryDataForApp `json:"app,omitempty"`
	Framework memoryData         `json:"framework"`
	System    memoryData         `json:"system"`
}
type memoryData struct {
	Peak  uint64 `json:"peak,omitempty"`
	Used  uint64 `json:"used,omitempty"`
	Total uint64 `json:"total,omitempty"`
}

type memoryDataForApp struct {
	AppID   uint64                    `json:"appid,omitempty"`
	Peak    uint64                    `json:"peak,omitempty"`
	Used    uint64                    `json:"used,omitempty"`
	Modules []memoryDataForAppModules `json:"modules,omitempty"`
}
type memoryDataForAppModules struct {
	File string `json:"file,omitempty"`
	Peak uint64 `json:"peak,omitempty"`
	Used uint64 `json:"used,omitempty"`
}
