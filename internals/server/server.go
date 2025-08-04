package server

import (
	"context"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func StartServer(router *mux.Router) {
	server := &http.Server{
		Handler: router,
		Addr:    ":8888",
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Printf("Error starting server: %v", err)
		}
	}()

	log.Printf("Server started on port: %s", server.Addr)
	log.Printf("Press CTRL+C to stop the server")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server shutdown failed:", err)
	}
	log.Println("Server exiting")
	os.Exit(0)
}
