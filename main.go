package main

import (
	"context"
	"embed"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var binStore = map[uuid.UUID]Bin{}
var binRecordStore = map[uuid.UUID]map[uuid.UUID]Record{}
var binWatchStore = map[uuid.UUID]map[*websocket.Conn]bool{}

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
}

func main() {
	{
		binKey := uuid.MustParse("d45a2464-4bce-4628-95be-8b8dfebe90be")
		now := time.Now()
		bin := NewBin(binKey, now)
		binStore[bin.binKey] = bin
		binRecordStore[bin.binKey] = map[uuid.UUID]Record{}
		binWatchStore[bin.binKey] = map[*websocket.Conn]bool{}
	}
	// staticFileServer := http.StripPrefix("/static", http.FileServer(http.Dir("./static/")))
	staticFileServer := http.FileServer(http.FS(staticFS))
	r := mux.NewRouter().StrictSlash(true)
	r.Path("/bin").Methods(http.MethodPost).Name("bin-create").HandlerFunc(binCreateHandler)
	r.Path("/bin/{binKey}").Methods(http.MethodGet).Name("bin-read").HandlerFunc(binReadHandler)
	r.Path("/bin/{binKey}/watch").Methods(http.MethodGet).Name("bin-watch").HandlerFunc(binWatchHandler)
	r.Path("/bin/{binKey}/records/{recordKey}").Methods(http.MethodGet).Name("record-read").HandlerFunc(binRecordReadHandler)
	r.PathPrefix("/record/{binKey}").Name("record-create").HandlerFunc(recordCreateHandler)
	r.PathPrefix("/static").Name("static").Handler(staticFileServer)
	r.PathPrefix("/").Methods(http.MethodGet).Name("root").HandlerFunc(rootHandler)
	r.Use(recoveryMiddleware)
	r.Use(corsMiddleware)
	r.Use(loggingMiddleware)
	r.Use(handlers.ProxyHeaders)
	if err := r.Walk(walk); err != nil {
		log.Print(err)
		os.Exit(1)
	}
	runServer(r, fmt.Sprintf("%s:%s", host, port))
}

func recoveryMiddleware(next http.Handler) http.Handler {
	recoveryHandler := handlers.RecoveryHandler(
		handlers.PrintRecoveryStack(true),
	)
	return recoveryHandler(next)
}

func corsMiddleware(next http.Handler) http.Handler {
	corsHandler := handlers.CORS(
		handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "PATCH", "DELETE"}),
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedHeaders([]string{"Content-Type"}),
	)
	return corsHandler(next)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return handlers.CombinedLoggingHandler(os.Stdout, next)
}

func walk(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
	pathTemplate, err := route.GetPathTemplate()
	if err != nil {
		fmt.Println("ROUTE:", pathTemplate)
	}
	pathRegexp, err := route.GetPathRegexp()
	if err == nil {
		fmt.Println("Path regexp:", pathRegexp)
	}
	queriesTemplates, err := route.GetQueriesTemplates()
	if err == nil {
		fmt.Println("Queries templates:", strings.Join(queriesTemplates, ","))
	}
	queriesRegexps, err := route.GetQueriesRegexp()
	if err == nil {
		fmt.Println("Queries regexps:", strings.Join(queriesRegexps, ","))
	}
	methods, err := route.GetMethods()
	if err == nil {
		fmt.Println("Methods:", strings.Join(methods, ","))
	}
	fmt.Println()
	return nil
}

func runServer(handler http.Handler, addr string) {
	srv := &http.Server{
		Handler: handler,
		Addr:    addr,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	go func() {
		log.Printf("serving at: %s", addr)
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	// Block until we receive our signal.
	<-c
	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Println("shutting down")
	os.Exit(0)
}
