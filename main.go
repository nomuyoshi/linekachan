package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/home", homeHandler).Methods("GET")
	router.HandleFunc("/callback", callbackHandler).Methods("POST")
	http.Handle("/", router)

	log.Print("Web Server starting port: 3000")
	if err := http.ListenAndServe(":3000", nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"url": "/home"})
}

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]string{"url": "/callback"})
}
