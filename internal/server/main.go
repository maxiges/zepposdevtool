package server

import (
	"net/http"

	"github.com/maxiges/ZeppOsDevTool/internal/api"

	"go.uber.org/zap"
)

const ServerPort = ":8081"

func RunServer() {

	//API HANDLERS
	http.HandleFunc("/app/echo", api.HandlerEcho)

	http.HandleFunc("/", api.HandlerDisplayMainPage)
	http.HandleFunc("/{appName}", api.HandlerDisplayData)

	http.HandleFunc("DELETE /add-data/{appName}", api.HandlerClearData)
	http.HandleFunc("POST /add-data/{appName}", api.HandlerAddData)

	logger, _ := zap.NewDevelopment()
	logger.Sugar().Infof("Server started at: localhost%s", ServerPort)

	//Run server
	http.ListenAndServe(ServerPort, nil)

}
