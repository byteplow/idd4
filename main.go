package main

import (
	"log"
	"net/http"

	"github.com/byteplow/idd4/internal/config"
	"github.com/byteplow/idd4/internal/container"
	"github.com/byteplow/idd4/routers"
	"github.com/gin-gonic/gin"
)

func init() {
	container.Setup()
}

func main() {
	if err := config.WatchConfig("/etc/byteplow/idd4/idd4.yml"); err != nil {
		log.Fatal(err)
	}

	gin.SetMode(config.Config.Server.RunMode)

	router := routers.InitRouter()

	server := &http.Server{
		Addr:    config.Config.Server.Endpoint,
		Handler: router,
	}

	log.Printf("listening on %s", config.Config.Server.Endpoint)

	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
