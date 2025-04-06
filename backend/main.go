package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/rs/cors"

	"go-parquet-read/routes"
)

func main() {

	router := routes.RouterConfig()

	port := 9000

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "DELETE", "PUT", "OPTIONS"},
	})

	server := http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		ReadTimeout:  90 * time.Second,
		WriteTimeout: 90 * time.Second,
		Handler:      c.Handler(router),
	}

	done := make(chan bool)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	//Graceful shut down
	go func() {
		<-quit
		log.Println("Server is shutting down...")

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		//Shutdown server
		server.SetKeepAlivesEnabled(false)
		if err := server.Shutdown(ctx); err != nil {
			log.Fatalf("Unable to gracefully shutdown the server: %v\n", err)
		}

		//Close channels
		close(quit)
		close(done)
	}()

	log.Printf("Listening on: %d", port)

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("Error in listening server: %s", err.Error())
	}

	<-done
	log.Fatal("Server stopped")

}
