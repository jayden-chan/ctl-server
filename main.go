package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/rs/cors"

	"github.com/jayden-chan/ctl-server/routes"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/register", routes.Register).Methods("POST")
	r.HandleFunc("/login", routes.Login).Methods("POST")
	r.HandleFunc("/deregister", routes.Deregister).Methods("DELETE")
	r.HandleFunc("/folders", routes.Folders).Methods("GET", "POST")
	r.HandleFunc("/folders/{folderID}", routes.FoldersID).Methods("DELETE", "PATCH")
	r.HandleFunc("/items", routes.Items).Methods("GET", "POST")
	r.HandleFunc("/items/{itemID}", routes.ItemsID).Methods("DELETE", "PATCH")

	port := os.Getenv("PORT")
	if port == "" {
		port = "9000"
	}

	// Add filename into logging messages
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Printf("Running server on port %s...\n", port)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           180,
	})

	handler := c.Handler(r)
	http.ListenAndServe(":"+port, handler)
}
