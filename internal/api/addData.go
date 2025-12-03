package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/maxiges/ZeppOsDevTool/internal/gui"
	"github.com/maxiges/ZeppOsDevTool/internal/models"
	"github.com/maxiges/ZeppOsDevTool/internal/storage"
	"go.uber.org/zap"
)

// HandlerAddData processes POST requests to add memory profiling data for a specific application.
//
// This handler:
// 1. Extracts the application name from the URL path parameter
// 2. Decodes incoming JSON data into a ZeppMemoryStruct
// 3. Automatically timestamps the incoming data
// 4. Retrieves existing data for the application from storage
// 5. Appends the new data point to the application's dataset
// 6. Persists the updated data to storage
// 7. Triggers a UI refresh if auto-refresh is enabled
//
// Request Path Parameters:
// - appName: the identifier of the application being profiled
//
// Request Body: JSON-encoded ZeppMemoryStruct containing memory metrics
//
// Response Status Codes:
// - 201 Created: Data successfully added
// - 500 Internal Server Error: JSON decoding failed
//
// Example POST: /api/data/my-app
// Body: {"memory": {...}, "description": "optional note"}
func HandlerAddData(w http.ResponseWriter, r *http.Request) {
	// Extract the application name from the URL path parameter
	appName := r.PathValue("appName")

	decoder := json.NewDecoder(r.Body)
	dataStruct := &models.ZeppMemoryStruct{}

	err := decoder.Decode(dataStruct)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	dataStruct.TimeStamp = time.Now()

	// Schedule a UI refresh after this function completes
	// This ensures the GUI updates with the new data if auto-refresh is enabled
	defer gui.TryRefreshUI(appName)

	val, exist := storage.GetDataForApp(appName)
	if !exist {
		val = []*models.ZeppMemoryStruct{}
	}
	// Append the newly received data point to the application's dataset
	val = append(val, dataStruct)
	// Persist the updated dataset to storage
	storage.SetData(appName, val)
	w.WriteHeader(http.StatusCreated)

}

// HandlerClearData processes POST/DELETE requests to clear all collected data for a specific application.
//
// This handler:
// 1. Extracts the application name from the URL path parameter
// 2. Removes all stored data points for the specified application
// 3. Sends a confirmation response
//
// This operation is useful for:
// - Starting fresh profiling sessions
// - Resetting data collection for an application
// - Cleaning up test data
//
// Request Path Parameters:
// - appName: the identifier of the application whose data should be cleared
//
// Response Status Codes:
// - 201 Created: Data successfully cleared (201 indicates the resource state has changed)
//
// Example: DELETE /api/data/my-app/clear
func HandlerClearData(w http.ResponseWriter, r *http.Request) {
	// Extract the application name from the URL path parameter
	appName := r.PathValue("appName")
	// Clear all data for this application by setting an empty slice
	storage.SetData(appName, []*models.ZeppMemoryStruct{})
	w.WriteHeader(http.StatusCreated)

}

func HandlerEcho(w http.ResponseWriter, r *http.Request) {
	logger, _ := zap.NewDevelopment()
	logger.Sugar().Info("req")

	if r.Body != nil {
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			fmt.Printf("Body reading error: %v", err)
			return
		}
		defer r.Body.Close()
		logger.Sugar().Info(string(bodyBytes))
	}

	logger.Sugar().Infof("Headers: %+v\n", r.Header)

	w.WriteHeader(http.StatusCreated)

}
