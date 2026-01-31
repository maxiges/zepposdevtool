package server

import (
	"net/http"
	"zepp-os-dev-tool/internal/api"

	"github.com/rs/cors"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"

	"go.uber.org/zap"
)

const ServerPort = ":8081"

func RunServer() {

	//API HANDLERS
	router := gin.Default()
	router.Use(static.Serve("/", static.LocalFile("./react/react-frontend/build", true)))

	go func() {
		if err := router.Run(":8082"); err != nil {
			panic(err)
		}
	}()

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "http://localhost:8081", "http://localhost:8082"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})
	mux := http.NewServeMux()

	mux.HandleFunc("/api/get-app-list", api.HandlerGetAppList)

	mux.HandleFunc("DELETE /api/app/{appName}", api.HandlerClearData)

	mux.HandleFunc("/api/data/{appName}", api.HandlerGetData)
	mux.HandleFunc("POST /api/data/{appName}", api.HandlerAddData)

	mux.HandleFunc("/api/logs/{appName}", api.HandlerGetLogs)
	mux.HandleFunc("POST /api/logs/{appName}", api.HandlerAddLog)

	logger, _ := zap.NewDevelopment()
	logger.Sugar().Infof("Server started at: localhost%s", ServerPort)

	handler := c.Handler(mux)
	//Run server
	err := http.ListenAndServe(ServerPort, handler)
	if err != nil {
		logger.Sugar().Error(err.Error())
	}

}
