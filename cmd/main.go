package main

import (
	"log"
	"net/http"

	"github.com/kpotier/carburant/pkg/api"
)

func main() {
	api, err := api.Start(log.Default(), "gob.db")
	if err != nil {
		log.Fatal(err)
	}
	defer api.Stop()

	fs := http.FileServer(http.Dir("./public/dist"))
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s: %s %s", r.RemoteAddr, r.Method, r.URL)
		fs.ServeHTTP(w, r)
	})
	mux.Handle("/api/", http.StripPrefix("/api", api.Handler()))
	http.ListenAndServe(":8080", mux)
}
