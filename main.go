package main

import (
	"embed"
	"fmt"
	"log"
	"os"

	"github.com/hellupline/winspector/pkg/datastore"
	"github.com/hellupline/winspector/pkg/server"
	"github.com/hellupline/winspector/pkg/service"
)

//go:embed static
var staticFS embed.FS

var host string
var port string

func init() {
	var ok bool
	host, ok = os.LookupEnv("HOST")
	if !ok {
		host = "127.0.0.1"
	}
	port, ok = os.LookupEnv("PORT")
	if !ok {
		port = "8000"
	}
	log.SetFlags(log.Ldate | log.Ltime | log.Llongfile | log.LUTC)
}

func main() {
	dataStore := datastore.NewDataStore()
	webapp := service.NewService(dataStore, staticFS)
	r := webapp.Router()
	r.Use(server.RecoveryMiddleware)
	r.Use(server.CorsMiddleware)
	r.Use(server.ProxyHeaders)
	// r.Use(server.LoggingMiddleware)
	if err := r.Walk(server.Walk); err != nil {
		log.Print(err)
		os.Exit(1)
	}
	server.RunServer(r, fmt.Sprintf("%s:%s", host, port))
}
